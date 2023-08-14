package types

import (
	"strconv"
	"time"
)

// ExhaustiveSearchRepoJob is a job that runs the exhaustive search on a repository.
// Maps to the `exhaustive_search_repo_jobs` database table.
type ExhaustiveSearchRepoJob struct {
	WorkerJob

	ID int64

	RepoID      int64
	SearchJobID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (j *ExhaustiveSearchRepoJob) RecordID() int {
	return int(j.ID)
}

func (j *ExhaustiveSearchRepoJob) RecordUID() string {
	return strconv.FormatInt(j.ID, 10)
}
