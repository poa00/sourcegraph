package workerutil

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/executor"
)

// Record is a generic interface for record conforming to the requirements of the store.
type Record interface {
	// RecordID returns the integer primary key of the record.
	RecordID() int
}

// Store is the persistence layer for the workerutil package that handles worker-side operations.
type Store[T Record] interface {
	// QueuedCount returns the number of records in the queued state.
	QueuedCount(ctx context.Context) (int, error)

	// Dequeue selects a record for processing. Any extra arguments supplied will be used in accordance with the
	// concrete persistence layer (e.g. additional SQL conditions for a database layer). This method returns a boolean
	// flag indicating the existence of a processable record.
	Dequeue(ctx context.Context, workerHostname string, extraArguments any) (T, bool, error)

	// Heartbeat updates last_heartbeat_at of all the given jobs, when they're processing. All IDs of records that were
	// touched are returned. Additionally, jobs in the working set that are flagged as to be canceled are returned.
	Heartbeat(ctx context.Context, jobIDs []int) (knownIDs, cancelIDs []int, err error)

	// AddExecutionLogEntry adds an executor log entry to the record and
	// returns the ID of the new entry (which can be used with
	// UpdateExecutionLogEntry) and a possible error.
	AddExecutionLogEntry(ctx context.Context, id int, entry executor.ExecutionLogEntry) (int, error)

	// UpdateExecutionLogEntry updates the executor log entry with the given ID
	// on the given record.
	UpdateExecutionLogEntry(ctx context.Context, recordID, entryID int, entry executor.ExecutionLogEntry) error

	// MarkComplete attempts to update the state of the record to complete. This method returns a boolean flag indicating
	// if the record was updated.
	MarkComplete(ctx context.Context, id int) (bool, error)

	// MarkErrored attempts to update the state of the record to errored. This method returns a boolean flag indicating
	// if the record was updated.
	MarkErrored(ctx context.Context, id int, failureMessage string) (bool, error)

	// MarkFailed attempts to update the state of the record to failed. This method returns a boolean flag indicating
	// if the record was updated.
	MarkFailed(ctx context.Context, id int, failureMessage string) (bool, error)
}
