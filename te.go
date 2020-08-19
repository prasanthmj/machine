package machine

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"log"
	"reflect"
	"sync"
)

type TaskExecutor interface {
	Execute(interface{}) error
}

type InternalTaskExecutor struct {
	access    sync.RWMutex
	executors map[string]TaskExecutor
}

func newTaskExecutor() *InternalTaskExecutor {
	te := &InternalTaskExecutor{}
	te.executors = make(map[string]TaskExecutor)
	return te
}

func (te *InternalTaskExecutor) DoTask(jb []uint8) error {
	buf := bytes.NewBuffer(jb)
	dec := gob.NewDecoder(buf)
	var job Job
	err := dec.Decode(&job)
	if err != nil {
		return err
	}
	log.Printf("Task ID %s received ", job.Task.GetTaskID())

	name := reflect.TypeOf(job.Task).String()
	te.access.RLock()
	tex, exists := te.executors[name]
	te.access.RUnlock()
	if !exists {
		err = fmt.Errorf("Executor for type %s not registered", name)
		return err
	}
	return tex.Execute(job.Task)
}
func (te *InternalTaskExecutor) Register(task interface{}, tex TaskExecutor) {
	gob.Register(task)
	te.access.Lock()
	defer te.access.Unlock()
	taskName := reflect.TypeOf(task).String()
	te.executors[taskName] = tex
}

func (te *InternalTaskExecutor) MakeTask(job *Job) (*tasks.Signature, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(job)
	if err != nil {
		return nil, err
	}
	bt := buf.Bytes()
	signature := &tasks.Signature{
		Name: "MachineTask",
		Args: []tasks.Arg{
			{
				Type:  "[]uint8",
				Value: bt,
			},
		},
	}
	if job.IsScheduled() {
		signature.ETA = &job.Due
	}
	return signature, nil
}
