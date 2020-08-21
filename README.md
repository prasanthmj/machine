# Machine
Wrapper around the awesome [machinery](https://github.com/RichardKnop/machinery) go library for background task execution 

Status: Work In progress

## Why?
* Limited and simplified wrapper.
* Pass around struct as task parameter (gets converted to gob)

### Usage
Create the JobQueue.
Pass the redis URL as parameter to New()
```go
jq, err := machine.New("redis://127.0.0.1:6379")

jq.Start()

```
Implement a "TaskExecuter" that can execute a "task"
TaskExecuter interface

```go
type TaskExecutor interface {
	Execute(interface{}) error
}
```
Task

```go
type RunnableTask interface {
	GetTaskID() string
}
```

Register the task executer with JobQueue
```go
jq.Register(&OrderEmail{}, orderSvc)
```

Now push tasks to the queue
```go
task := &OrderEmail{email, orderID}
job := machine.NewJob(&task)
jq.QueueUp(job)
```
Your TaskExecutor gets called when the task is due
```go
func (*OrderService) Execute(t interface{})error{
    switch t.(type){
        case OrderEmail:
            return sendEmail(t)
        case PaymentUpdate:
            return updatePayment(t)
    }
    return nil
}
```

## Delayed jobs
```go
task := &ConfirmationEmail{email}

job := machine.NewJob(task).After(5*time.Minutes)
jq.QueueUp(job)

```

## Stop the queue (and the workers)
```go
jq.Stop()
```