package machine

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"sync"
	"time"
)

type JobQueue struct {
	server           *machinery.Server
	worker           *machinery.Worker
	recurringTickers []*time.Ticker
	NumWorkers       int
	wg               sync.WaitGroup
	te               *InternalTaskExecutor
	executors        map[string]TaskExecutor
	close            chan struct{}
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
	log.SetDebug(NewEmptyLog())
	log.SetInfo(NewEmptyLog())

	jq := &JobQueue{server: server, NumWorkers: 10}
	jq.te = newTaskExecutor()
	jq.close = make(chan struct{})

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

	for _, ticker := range jq.recurringTickers {
		ticker.Stop()
	}
	if jq.worker == nil {
		return
	}
	close(jq.close)
	jq.worker.Quit()
	jq.wg.Wait()
}

func (jq *JobQueue) Register(task interface{}, tex TaskExecutor) {
	jq.te.Register(task, tex)
}

func (jq *JobQueue) ScheduleRecurringJob(job *Job, repeat time.Duration) {
	ticker := jq.runRecurring(job, repeat)
	jq.recurringTickers = append(jq.recurringTickers, ticker)
}

func (jq *JobQueue) runRecurring(job *Job, repeat time.Duration) *time.Ticker {
	ticker := time.NewTicker(repeat)
	jq.wg.Add(1)
	go func() {
		defer jq.wg.Done()
		for {
			select {
			case <-ticker.C:
				jq.QueueUp(job)
			case <-jq.close:
				return
			}
		}
	}()
	return ticker
}
