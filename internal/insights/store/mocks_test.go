// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package store

import (
	"context"
	"sync"

	sqlf "github.com/keegancsmith/sqlf"
	api "github.com/sourcegraph/sourcegraph/internal/api"
)

// MockInsightPermissionStore is a mock implementation of the
// InsightPermissionStore interface (from the package
// github.com/sourcegraph/sourcegraph/internal/insights/store) used for unit
// testing.
type MockInsightPermissionStore struct {
	// GetUnauthorizedRepoIDsFunc is an instance of a mock function object
	// controlling the behavior of the method GetUnauthorizedRepoIDs.
	GetUnauthorizedRepoIDsFunc *InsightPermissionStoreGetUnauthorizedRepoIDsFunc
	// GetUnauthorizedRepoIDsQueryFunc is an instance of a mock function
	// object controlling the behavior of the method
	// GetUnauthorizedRepoIDsQuery.
	GetUnauthorizedRepoIDsQueryFunc *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc
	// GetUserPermissionsFunc is an instance of a mock function object
	// controlling the behavior of the method GetUserPermissions.
	GetUserPermissionsFunc *InsightPermissionStoreGetUserPermissionsFunc
}

// NewMockInsightPermissionStore creates a new mock of the
// InsightPermissionStore interface. All methods return zero values for all
// results, unless overwritten.
func NewMockInsightPermissionStore() *MockInsightPermissionStore {
	return &MockInsightPermissionStore{
		GetUnauthorizedRepoIDsFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsFunc{
			defaultHook: func(context.Context) (r0 []api.RepoID, r1 error) {
				return
			},
		},
		GetUnauthorizedRepoIDsQueryFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc{
			defaultHook: func(context.Context) (r0 *sqlf.Query, r1 error) {
				return
			},
		},
		GetUserPermissionsFunc: &InsightPermissionStoreGetUserPermissionsFunc{
			defaultHook: func(context.Context) (r0 []int, r1 []int, r2 error) {
				return
			},
		},
	}
}

// NewStrictMockInsightPermissionStore creates a new mock of the
// InsightPermissionStore interface. All methods panic on invocation, unless
// overwritten.
func NewStrictMockInsightPermissionStore() *MockInsightPermissionStore {
	return &MockInsightPermissionStore{
		GetUnauthorizedRepoIDsFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsFunc{
			defaultHook: func(context.Context) ([]api.RepoID, error) {
				panic("unexpected invocation of MockInsightPermissionStore.GetUnauthorizedRepoIDs")
			},
		},
		GetUnauthorizedRepoIDsQueryFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc{
			defaultHook: func(context.Context) (*sqlf.Query, error) {
				panic("unexpected invocation of MockInsightPermissionStore.GetUnauthorizedRepoIDsQuery")
			},
		},
		GetUserPermissionsFunc: &InsightPermissionStoreGetUserPermissionsFunc{
			defaultHook: func(context.Context) ([]int, []int, error) {
				panic("unexpected invocation of MockInsightPermissionStore.GetUserPermissions")
			},
		},
	}
}

// NewMockInsightPermissionStoreFrom creates a new mock of the
// MockInsightPermissionStore interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockInsightPermissionStoreFrom(i InsightPermissionStore) *MockInsightPermissionStore {
	return &MockInsightPermissionStore{
		GetUnauthorizedRepoIDsFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsFunc{
			defaultHook: i.GetUnauthorizedRepoIDs,
		},
		GetUnauthorizedRepoIDsQueryFunc: &InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc{
			defaultHook: i.GetUnauthorizedRepoIDsQuery,
		},
		GetUserPermissionsFunc: &InsightPermissionStoreGetUserPermissionsFunc{
			defaultHook: i.GetUserPermissions,
		},
	}
}

// InsightPermissionStoreGetUnauthorizedRepoIDsFunc describes the behavior
// when the GetUnauthorizedRepoIDs method of the parent
// MockInsightPermissionStore instance is invoked.
type InsightPermissionStoreGetUnauthorizedRepoIDsFunc struct {
	defaultHook func(context.Context) ([]api.RepoID, error)
	hooks       []func(context.Context) ([]api.RepoID, error)
	history     []InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall
	mutex       sync.Mutex
}

// GetUnauthorizedRepoIDs delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockInsightPermissionStore) GetUnauthorizedRepoIDs(v0 context.Context) ([]api.RepoID, error) {
	r0, r1 := m.GetUnauthorizedRepoIDsFunc.nextHook()(v0)
	m.GetUnauthorizedRepoIDsFunc.appendCall(InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// GetUnauthorizedRepoIDs method of the parent MockInsightPermissionStore
// instance is invoked and the hook queue is empty.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) SetDefaultHook(hook func(context.Context) ([]api.RepoID, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetUnauthorizedRepoIDs method of the parent MockInsightPermissionStore
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) PushHook(hook func(context.Context) ([]api.RepoID, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) SetDefaultReturn(r0 []api.RepoID, r1 error) {
	f.SetDefaultHook(func(context.Context) ([]api.RepoID, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) PushReturn(r0 []api.RepoID, r1 error) {
	f.PushHook(func(context.Context) ([]api.RepoID, error) {
		return r0, r1
	})
}

func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) nextHook() func(context.Context) ([]api.RepoID, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) appendCall(r0 InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall objects describing
// the invocations of this function.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsFunc) History() []InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall {
	f.mutex.Lock()
	history := make([]InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall is an object that
// describes an invocation of method GetUnauthorizedRepoIDs on an instance
// of MockInsightPermissionStore.
type InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []api.RepoID
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightPermissionStoreGetUnauthorizedRepoIDsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc describes the
// behavior when the GetUnauthorizedRepoIDsQuery method of the parent
// MockInsightPermissionStore instance is invoked.
type InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc struct {
	defaultHook func(context.Context) (*sqlf.Query, error)
	hooks       []func(context.Context) (*sqlf.Query, error)
	history     []InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall
	mutex       sync.Mutex
}

// GetUnauthorizedRepoIDsQuery delegates to the next hook function in the
// queue and stores the parameter and result values of this invocation.
func (m *MockInsightPermissionStore) GetUnauthorizedRepoIDsQuery(v0 context.Context) (*sqlf.Query, error) {
	r0, r1 := m.GetUnauthorizedRepoIDsQueryFunc.nextHook()(v0)
	m.GetUnauthorizedRepoIDsQueryFunc.appendCall(InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// GetUnauthorizedRepoIDsQuery method of the parent
// MockInsightPermissionStore instance is invoked and the hook queue is
// empty.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) SetDefaultHook(hook func(context.Context) (*sqlf.Query, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetUnauthorizedRepoIDsQuery method of the parent
// MockInsightPermissionStore instance invokes the hook at the front of the
// queue and discards it. After the queue is empty, the default hook
// function is invoked for any future action.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) PushHook(hook func(context.Context) (*sqlf.Query, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) SetDefaultReturn(r0 *sqlf.Query, r1 error) {
	f.SetDefaultHook(func(context.Context) (*sqlf.Query, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) PushReturn(r0 *sqlf.Query, r1 error) {
	f.PushHook(func(context.Context) (*sqlf.Query, error) {
		return r0, r1
	})
}

func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) nextHook() func(context.Context) (*sqlf.Query, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) appendCall(r0 InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall objects
// describing the invocations of this function.
func (f *InsightPermissionStoreGetUnauthorizedRepoIDsQueryFunc) History() []InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall {
	f.mutex.Lock()
	history := make([]InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall is an object
// that describes an invocation of method GetUnauthorizedRepoIDsQuery on an
// instance of MockInsightPermissionStore.
type InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *sqlf.Query
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightPermissionStoreGetUnauthorizedRepoIDsQueryFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// InsightPermissionStoreGetUserPermissionsFunc describes the behavior when
// the GetUserPermissions method of the parent MockInsightPermissionStore
// instance is invoked.
type InsightPermissionStoreGetUserPermissionsFunc struct {
	defaultHook func(context.Context) ([]int, []int, error)
	hooks       []func(context.Context) ([]int, []int, error)
	history     []InsightPermissionStoreGetUserPermissionsFuncCall
	mutex       sync.Mutex
}

// GetUserPermissions delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockInsightPermissionStore) GetUserPermissions(v0 context.Context) ([]int, []int, error) {
	r0, r1, r2 := m.GetUserPermissionsFunc.nextHook()(v0)
	m.GetUserPermissionsFunc.appendCall(InsightPermissionStoreGetUserPermissionsFuncCall{v0, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the GetUserPermissions
// method of the parent MockInsightPermissionStore instance is invoked and
// the hook queue is empty.
func (f *InsightPermissionStoreGetUserPermissionsFunc) SetDefaultHook(hook func(context.Context) ([]int, []int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetUserPermissions method of the parent MockInsightPermissionStore
// instance invokes the hook at the front of the queue and discards it.
// After the queue is empty, the default hook function is invoked for any
// future action.
func (f *InsightPermissionStoreGetUserPermissionsFunc) PushHook(hook func(context.Context) ([]int, []int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *InsightPermissionStoreGetUserPermissionsFunc) SetDefaultReturn(r0 []int, r1 []int, r2 error) {
	f.SetDefaultHook(func(context.Context) ([]int, []int, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *InsightPermissionStoreGetUserPermissionsFunc) PushReturn(r0 []int, r1 []int, r2 error) {
	f.PushHook(func(context.Context) ([]int, []int, error) {
		return r0, r1, r2
	})
}

func (f *InsightPermissionStoreGetUserPermissionsFunc) nextHook() func(context.Context) ([]int, []int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *InsightPermissionStoreGetUserPermissionsFunc) appendCall(r0 InsightPermissionStoreGetUserPermissionsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// InsightPermissionStoreGetUserPermissionsFuncCall objects describing the
// invocations of this function.
func (f *InsightPermissionStoreGetUserPermissionsFunc) History() []InsightPermissionStoreGetUserPermissionsFuncCall {
	f.mutex.Lock()
	history := make([]InsightPermissionStoreGetUserPermissionsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// InsightPermissionStoreGetUserPermissionsFuncCall is an object that
// describes an invocation of method GetUserPermissions on an instance of
// MockInsightPermissionStore.
type InsightPermissionStoreGetUserPermissionsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []int
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 []int
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c InsightPermissionStoreGetUserPermissionsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c InsightPermissionStoreGetUserPermissionsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}
