package machine

import (
	"encoding/gob"
	"time"
)

type RunnableTask interface {
	GetTaskID() string
}

type Job struct {
	Task RunnableTask
	Due  time.Time
}

func NewJob(t RunnableTask) *Job {
	gob.Register(&Job{})
	return &Job{Task: t}
}

func (j *Job) After(d time.Duration) *Job {
	j.Due = time.Now().Add(d)
	return j
}

func (j *Job) IsScheduled() bool {
	if j.Due.IsZero() {
		return false
	}
	return true
}
