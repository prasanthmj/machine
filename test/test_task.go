package test

import (
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

type TestTask struct {
	TaskTime time.Duration
	TaskID   string
}

func (task *TestTask) WithTaskTime(d time.Duration) *TestTask {
	task.TaskTime = d
	return task
}

func (task *TestTask) GetTaskID() string {
	return task.TaskID
}

func (task *TestTask) RunTask(t *testing.T) error {

	t.Logf("Test Task:%s Running run time %v ", task.TaskID, task.TaskTime)
	//t.logger.Logf("TestTask Running ... ")

	if task.TaskTime > 0 {
		time.Sleep(task.TaskTime)
	}
	t.Logf("TestTask Completed.")
	//t.logger.Logf("TestTask Completed.")
	return nil
}
func CreateTestTask(te *TestTaskExecutor) *TestTask {
	taskID := faker.RandomString(12)
	te.InitTask(taskID)
	return &TestTask{TaskID: taskID}
}
