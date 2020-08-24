package machine_test

import (
	"github.com/prasanthmj/machine"
	mtest "github.com/prasanthmj/machine/test"
	"log"
	"math/rand"
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

//TODO: machinery disable tracing
func TestStartAndStop(t *testing.T) {
	//defer goleak.VerifyNone(t)
	jq, err := machine.New("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Error creating queue %v", err)
	}
	t.Logf("Starting job queue ...")
	jq.Start()
	<-time.After(2 * time.Second)
	t.Logf("stopping the queue ...")
	jq.Stop()
	t.Logf("Done.")
}

type MyTask struct {
	ID        string
	CreatedAt time.Time
}

func (t *MyTask) GetTaskID() string {
	return t.ID
}

type MyTaskExecutor struct {
}

func (te *MyTaskExecutor) Execute(p interface{}) error {
	task := p.(*MyTask)
	log.Printf("Executing task with ID %s ", task.ID)

	return nil
}

func TestSendingJobs(t *testing.T) {
	myte := &MyTaskExecutor{}
	jq, err := machine.New("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Error creating queue %v", err)
	}
	jq.Start()
	jq.Register(&MyTask{}, myte)
	<-time.After(time.Second)
	var task MyTask
	task.ID = faker.RandomString(8)
	task.CreatedAt = time.Now()
	job := machine.NewJob(&task)
	err = jq.QueueUp(job)
	if err != nil {
		t.Errorf("Error queueing up job %v", err)
	}
	jq.Stop()
	t.Logf("Done.")
}

func TestAllAreJobsRun(t *testing.T) {
	myte := mtest.NewTaskExecutor(t)

	jq, err := machine.New("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Error creating queue %v", err)
	}
	jq.Start()
	jq.Register(&mtest.TestTask{}, myte)
	for i := 0; i < 100; i++ {
		task := mtest.CreateTestTask(myte).WithTaskTime(time.Duration(rand.Intn(100)) * time.Millisecond)
		jq.QueueUp(machine.NewJob(task))
	}
	<-time.After(1 * time.Second)
	jq.Stop()
	myte.PrintStatus()
	if !myte.AssertAllTasksExecutedExactlyOnce() {
		t.Errorf("Expected all tasks to execute exactly once")
	}

}

func TestDelayedJob(t *testing.T) {
	myte := mtest.NewTaskExecutor(t)

	jq, err := machine.New("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Error creating queue %v", err)
	}
	jq.Start()
	jq.Register(&mtest.TestTask{}, myte)
	task := mtest.CreateTestTask(myte).WithTaskTime(time.Duration(rand.Intn(100)) * time.Millisecond)
	job := machine.NewJob(task).After(500 * time.Millisecond)
	jq.QueueUp(job)
	<-time.After(1 * time.Second)
	jq.Stop()
	myte.PrintStatus()
	if myte.GetExecutionCount(task.GetTaskID()) != 1 {
		t.Errorf("Delayed task didn't execute ")
	}

	time_exec := myte.GetTaskExecutedAt(task.GetTaskID())

	time_created := myte.GetTaskCreatedAt(task.GetTaskID())

	d := time_exec.Sub(time_created)
	if d <= (500 * time.Millisecond) {
		t.Errorf("Delayed task executed before required delay ")
	}

}
func TestRecurringJob(t *testing.T) {
	myte := mtest.NewTaskExecutor(t)

	jq, err := machine.New("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatalf("Error creating queue %v", err)
	}
	jq.Start()
	jq.Register(&mtest.TestTask{}, myte)
	task := mtest.CreateTestTask(myte).WithTaskTime(time.Duration(rand.Intn(100)) * time.Millisecond)
	job := machine.NewJob(task)
	jq.ScheduleRecurringJob(job, 200*time.Millisecond)
	<-time.After(1 * time.Second)
	jq.Stop()
	myte.PrintStatus()

	if !myte.AssertAllTasksExecutedAtleastOnce() {
		t.Errorf("Expected task to execute at least once")
	}
	if myte.GetExecutionCount(task.GetTaskID()) < 3 {
		t.Errorf("Expected the recurring task to execute at leat 3 times ")
	}
}
