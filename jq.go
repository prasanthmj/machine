package machine

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"sync"
)

type JobQueue struct {
	server     *machinery.Server
	worker     *machinery.Worker
	NumWorkers int
	wg         sync.WaitGroup
	te         *InternalTaskExecutor
	executors  map[string]TaskExecutor
}

func New(redisURL string) (*JobQueue, error) {
	var cnf = config.Config{
		Broker:        redisURL,
		ResultBackend: redisURL,
		NoUnixSignals: true,
	}
	server, err := machinery.NewServer(&cnf)
	if err != nil {
		return nil, err
	}

	jq := &JobQueue{server: server, NumWorkers: 10}
	jq.te = newTaskExecutor()

	err = jq.server.RegisterTask("MachineTask", jq.te.DoTask)
	if err != nil {
		return nil, err
	}
	return jq, nil
}

func (d *JobQueue) Workers(n int) *JobQueue {
	d.NumWorkers = n
	return d
}

func (jq *JobQueue) Start() {

	jq.wg.Add(1)
	go func() {
		defer jq.wg.Done()
		jq.worker = jq.server.NewWorker("main_worker", jq.NumWorkers)
		err := jq.worker.Launch()
		if err != nil {
			//Log error
		}
	}()

}

func (jq *JobQueue) QueueUp(job *Job) error {
	signature, err := jq.te.MakeTask(job)
	if err != nil {
		return err
	}
	jq.server.SendTask(signature)
	return nil
}

func (jq *JobQueue) Stop() {
	if jq.worker == nil {
		return
	}

	jq.worker.Quit()
	jq.wg.Wait()
}

func (jq *JobQueue) Register(task interface{}, tex TaskExecutor) {
	jq.te.Register(task, tex)
}
