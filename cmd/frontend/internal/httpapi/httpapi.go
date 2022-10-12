package httpapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/graph-gophers/graphql-go"
	"github.com/inconshreveable/log15"

	sglog "github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/enterprise"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/envvar"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/app/updatecheck"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/handlerutil"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/httpapi/releasecache"
	apirouter "github.com/sourcegraph/sourcegraph/cmd/frontend/internal/httpapi/router"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/httpapi/webhookhandlers"
	frontendsearch "github.com/sourcegraph/sourcegraph/cmd/frontend/internal/search"
	registry "github.com/sourcegraph/sourcegraph/cmd/frontend/registry/api"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/webhooks"
	"github.com/sourcegraph/sourcegraph/internal/actor"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/authz"
	internalcodeintel "github.com/sourcegraph/sourcegraph/internal/codeintel"
	"github.com/sourcegraph/sourcegraph/internal/conf"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/encryption/keyring"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/repos"
	"github.com/sourcegraph/sourcegraph/internal/search"
	"github.com/sourcegraph/sourcegraph/internal/search/searchcontexts"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type Handlers struct {
	GitHubWebhook                   webhooks.Registerer
	GitLabWebhook                   http.Handler
	BitbucketServerWebhook          http.Handler
	BitbucketCloudWebhook           http.Handler
	BatchesChangesFileGetHandler    http.Handler
	BatchesChangesFileExistsHandler http.Handler
	BatchesChangesFileUploadHandler http.Handler
	NewCodeIntelUploadHandler       enterprise.NewCodeIntelUploadHandler
	NewComputeStreamHandler         enterprise.NewComputeStreamHandler
}

// NewHandler returns a new API handler that uses the provided API
// router, which must have been created by httpapi/router.New, or
// creates a new one if nil.
//
// 🚨 SECURITY: The caller MUST wrap the returned handler in middleware that checks authentication
// and sets the actor in the request context.
func NewHandler(
	db database.DB,
	m *mux.Router,
	schema *graphql.Schema,
	rateLimiter graphqlbackend.LimitWatcher,
	handlers *Handlers,
) http.Handler {
	logger := sglog.Scoped("Handler", "frontend HTTP API handler")

	if m == nil {
		m = apirouter.New(nil)
	}
	m.StrictSlash(true)

	handler := jsonMiddleware(&errorHandler{
		// Only display error message to admins when in debug mode, since it
		// may contain sensitive info (like API keys in net/http error
		// messages).
		WriteErrBody: env.InsecureDev,
	})

	// Set handlers for the installed routes.
	m.Get(apirouter.RepoShield).Handler(trace.Route(handler(serveRepoShield(db))))

	m.Get(apirouter.RepoRefresh).Handler(trace.Route(handler(serveRepoRefresh(db))))

	gh := webhooks.GitHubWebhook{
		ExternalServices: db.ExternalServices(),
	}

	webhookhandlers.Init(db, &gh)
	webhookMiddleware := webhooks.NewLogMiddleware(
		db.WebhookLogs(keyring.Default().WebhookLogKey),
	)

	handlers.GitHubWebhook.Register(&gh)

	m.Get(apirouter.GitHubWebhooks).Handler(trace.Route(webhookMiddleware.Logger(&gh)))
	m.Get(apirouter.GitLabWebhooks).Handler(trace.Route(webhookMiddleware.Logger(handlers.GitLabWebhook)))
	m.Get(apirouter.BitbucketServerWebhooks).Handler(trace.Route(webhookMiddleware.Logger(handlers.BitbucketServerWebhook)))
	m.Get(apirouter.BitbucketCloudWebhooks).Handler(trace.Route(webhookMiddleware.Logger(handlers.BitbucketCloudWebhook)))
	m.Get(apirouter.BatchesFileGet).Handler(trace.Route(handlers.BatchesChangesFileGetHandler))
	m.Get(apirouter.BatchesFileExists).Handler(trace.Route(handlers.BatchesChangesFileExistsHandler))
	m.Get(apirouter.BatchesFileUpload).Handler(trace.Route(handlers.BatchesChangesFileUploadHandler))
	m.Get(apirouter.LSIFUpload).Handler(trace.Route(handlers.NewCodeIntelUploadHandler(true)))
	m.Get(apirouter.ComputeStream).Handler(trace.Route(handlers.NewComputeStreamHandler()))

	ghSync := repos.GitHubWebhookHandler{}
	ghSync.Register(&gh)

	if envvar.SourcegraphDotComMode() {
		m.Path("/updates").Methods("GET", "POST").Name("updatecheck").Handler(trace.Route(http.HandlerFunc(updatecheck.HandlerWithLog(logger))))
	}

	m.Get(apirouter.GraphQL).Handler(trace.Route(handler(serveGraphQL(logger, schema, rateLimiter, false))))

	m.Get(apirouter.SearchStream).Handler(trace.Route(restTokenScopesMiddleware([]string{"search-stream:read"}, frontendsearch.StreamHandler(db))))

	// Return the minimum src-cli version that's compatible with this instance
	m.Get(apirouter.SrcCli).Handler(trace.Route(restTokenScopesMiddleware([]string{"src-cli:read"}, newSrcCliVersionHandler(logger))))

	// Set up the src-cli version cache handler (this will effectively be a
	// no-op anywhere other than dot-com).
	m.Get(apirouter.SrcCliVersionCache).Handler(trace.Route(restTokenScopesMiddleware([]string{"src-cli:read"}, releasecache.NewHandler(logger))))

	m.Get(apirouter.Registry).Handler(trace.Route(handler(registry.HandleRegistry(db))))

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API no route: %s %s from %s", r.Method, r.URL, r.Referer())
		http.Error(w, "no route", http.StatusNotFound)
	})

	return m
}

// NewInternalHandler returns a new API handler for internal endpoints that uses
// the provided API router, which must have been created by httpapi/router.NewInternal.
//
// 🚨 SECURITY: This handler should not be served on a publicly exposed port. 🚨
// This handler is not guaranteed to provide the same authorization checks as
// public API handlers.
func NewInternalHandler(
	m *mux.Router,
	db database.DB,
	schema *graphql.Schema,
	newCodeIntelUploadHandler enterprise.NewCodeIntelUploadHandler,
	newComputeStreamHandler enterprise.NewComputeStreamHandler,
	rateLimitWatcher graphqlbackend.LimitWatcher,
	codeIntelServices internalcodeintel.Services,
) http.Handler {
	logger := sglog.Scoped("InternalHandler", "")
	if m == nil {
		m = apirouter.New(nil)
	}
	m.StrictSlash(true)

	handler := jsonMiddleware(&errorHandler{
		// Internal endpoints can expose sensitive errors
		WriteErrBody: true,
	})

	m.Get(apirouter.ExternalServiceConfigs).Handler(trace.Route(handler(serveExternalServiceConfigs(db))))

	// zoekt-indexserver endpoints
	gsClient := gitserver.NewClient(db)
	indexer := &searchIndexerServer{
		db:            db,
		ListIndexable: backend.NewRepos(logger, db, gsClient).ListIndexable,
		RepoStore:     db.Repos(),
		SearchContextsRepoRevs: func(ctx context.Context, repoIDs []api.RepoID) (map[api.RepoID][]string, error) {
			return searchcontexts.RepoRevs(ctx, db, repoIDs)
		},
		Indexers: search.Indexers(),

		MinLastChangedDisabled: os.Getenv("SRC_SEARCH_INDEXER_EFFICIENT_POLLING_DISABLED") != "",

		codeIntelServices: codeIntelServices,
	}
	m.Get(apirouter.SearchConfiguration).Handler(trace.Route(handler(indexer.serveConfiguration)))
	m.Get(apirouter.ReposIndex).Handler(trace.Route(handler(indexer.serveList)))
	m.Get(apirouter.RepoRank).Handler(trace.Route(handler(indexer.serveRepoRank)))
	m.Get(apirouter.DocumentRanks).Handler(trace.Route(handler(indexer.serveDocumentRanks)))

	m.Get(apirouter.ExternalURL).Handler(trace.Route(handler(serveExternalURL)))
	m.Get(apirouter.SendEmail).Handler(trace.Route(handler(serveSendEmail)))
	gitService := &gitServiceHandler{
		Gitserver: gsClient,
	}
	m.Get(apirouter.GitInfoRefs).Handler(trace.Route(handler(gitService.serveInfoRefs())))
	m.Get(apirouter.GitUploadPack).Handler(trace.Route(handler(gitService.serveGitUploadPack())))
	m.Get(apirouter.Telemetry).Handler(trace.Route(telemetryHandler(db)))
	m.Get(apirouter.GraphQL).Handler(trace.Route(handler(serveGraphQL(logger, schema, rateLimitWatcher, true))))
	m.Get(apirouter.Configuration).Handler(trace.Route(handler(serveConfiguration)))
	m.Path("/ping").Methods("GET").Name("ping").HandlerFunc(handlePing)
	m.Get(apirouter.StreamingSearch).Handler(trace.Route(frontendsearch.StreamHandler(db)))
	m.Get(apirouter.ComputeStream).Handler(trace.Route(newComputeStreamHandler()))

	m.Get(apirouter.LSIFUpload).Handler(trace.Route(newCodeIntelUploadHandler(false)))

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API no route: %s %s from %s", r.Method, r.URL, r.Referer())
		http.Error(w, "no route", http.StatusNotFound)
	})

	return m
}

var schemaDecoder = schema.NewDecoder()

func init() {
	schemaDecoder.IgnoreUnknownKeys(true)

	// Register a converter for unix timestamp strings -> time.Time values
	// (needed for Appdash PageLoadEvent type).
	schemaDecoder.RegisterConverter(time.Time{}, func(s string) reflect.Value {
		ms, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.ValueOf(err)
		}
		return reflect.ValueOf(time.Unix(0, ms*int64(time.Millisecond)))
	})
}

type errorHandler struct {
	WriteErrBody bool
}

func (h *errorHandler) Handle(w http.ResponseWriter, r *http.Request, status int, err error) {
	trace.SetRequestErrorCause(r.Context(), err)

	// Handle custom errors
	var e *handlerutil.URLMovedError
	if errors.As(err, &e) {
		err := handlerutil.RedirectToNewRepoName(w, r, e.NewRepo)
		if err != nil {
			log15.Error("error redirecting to new URI", "err", err, "new_url", e.NewRepo)
		}
		return
	}

	// Never cache error responses.
	w.Header().Set("cache-control", "no-cache, max-age=0")

	errBody := err.Error()

	var displayErrBody string
	if h.WriteErrBody {
		displayErrBody = errBody
	}
	http.Error(w, displayErrBody, status)
	traceID := trace.ID(r.Context())
	traceURL := trace.URL(traceID, conf.DefaultClient())

	if status < 200 || status >= 500 {
		log15.Error("API HTTP handler error response", "method", r.Method, "request_uri", r.URL.RequestURI(), "status_code", status, "error", err, "trace", traceURL, "traceID", traceID)
	}
}

func jsonMiddleware(errorHandler *errorHandler) func(func(http.ResponseWriter, *http.Request) error) http.Handler {
	return func(h func(http.ResponseWriter, *http.Request) error) http.Handler {
		return handlerutil.HandlerWithErrorReturn{
			Handler: func(w http.ResponseWriter, r *http.Request) error {
				w.Header().Set("Content-Type", "application/json")
				return h(w, r)
			},
			Error: errorHandler.Handle,
		}
	}
}

func restTokenScopesMiddleware(scopes []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(scopes) > 0 {
			a := actor.FromContext(r.Context())
			// only care about token based auth and non-internal tokens for now
			isUserAll := a.Scopes[authz.ScopeUserAll]
			isSiteAdminSudo := a.Scopes[authz.ScopeSiteAdminSudo]
			if a.FromToken && !a.Internal && !(isUserAll || isSiteAdminSudo) {
				// for now all scopes are required, but this can be changed in the future
				// to be more flexible
				for _, scope := range scopes {
					if ok := a.Scopes[scope]; !ok {
						http.Error(w, fmt.Sprintf("forbidden, missing token scope: %s", scope), http.StatusForbidden)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
		return
	})
}
