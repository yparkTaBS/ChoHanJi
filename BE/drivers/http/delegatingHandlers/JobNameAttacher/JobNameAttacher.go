package JobNameAttacher

import (
	"ChoHanJi/infrastructure/ContextKeys"
	"context"
	"net/http"
)

type JobNameAttacher struct {
	jobName string
	next    http.Handler
}

var _ http.Handler = (*JobNameAttacher)(nil)

func New(jobName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return new(jobName, next)
	}
}

func new(jobName string, next http.Handler) *JobNameAttacher {
	return &JobNameAttacher{jobName, next}
}

// ServeHTTP implements http.Handler.
func (j *JobNameAttacher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextKeys.JobName, j.jobName)
	r = r.WithContext(ctx)
	j.next.ServeHTTP(w, r)
}
