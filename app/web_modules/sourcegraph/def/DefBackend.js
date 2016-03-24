import * as DefActions from "sourcegraph/def/DefActions";
import DefStore from "sourcegraph/def/DefStore";
import Dispatcher from "sourcegraph/Dispatcher";
import defaultXhr from "sourcegraph/util/xhr";

const DefBackend = {
	xhr: defaultXhr,

	__onDispatch(action) {
		switch (action.constructor) {
		case DefActions.WantDef:
			{
				let def = DefStore.defs.get(action.url);
				if (def === null) {
					DefBackend.xhr({
						uri: `/.ui${action.url}`,
						headers: {
							"X-Definition-Data-Only": "yes",
						},
						json: {},
					}, function(err, resp, body) {
						if (err) {
							console.error(err);
							return;
						}
						Dispatcher.dispatch(new DefActions.DefFetched(action.url, body));
					});
				}
				break;
			}

		case DefActions.WantDefs:
			{
				let defs = DefStore.defs.list(action.repo, action.rev, action.query);
				if (defs === null) {
					DefBackend.xhr({
						uri: `/.api/.defs?RepoRevs=${encodeURIComponent(action.repo)}@${encodeURIComponent(action.rev)}&Nonlocal=true&Query=${encodeURIComponent(action.query)}`,
						json: {},
					}, function(err, resp, body) {
						if (err) {
							console.error(err);
							return;
						}
						Dispatcher.dispatch(new DefActions.DefsFetched(action.repo, action.rev, action.query, body));
					});
				}
				break;
			}

		case DefActions.WantExample:
			{
				let example = DefStore.examples.get(action.defURL, action.index);
				if (example === null && action.index < DefStore.examples.getCount(action.defURL)) {
					DefBackend.xhr({
						uri: `/.ui${action.defURL}/.examples?PerPage=1&Page=${action.index + 1}`,
						json: {},
					}, function(err, resp, body) {
						if (err) {
							console.error(err);
							return;
						}
						if (body === null || body.Error) {
							Dispatcher.dispatch(new DefActions.NoExampleAvailable(action.defURL, action.index));
							return;
						}
						Dispatcher.dispatch(new DefActions.ExampleFetched(action.defURL, action.index, body[0]));
					});
				}
				break;
			}

		}
	},
};

Dispatcher.register(DefBackend.__onDispatch);

export default DefBackend;
