package model

//------------------------------------------------------------------------------

// TaskStatus resembles the state of a task.
type TaskStatus int

// Enumeration of possible states of a task.
const (
	// TaskStatusInitial resembles the initial state of a task
	TaskStatusInitial TaskStatus = iota
	// TaskStatusExecuting resembles the execution state of a task
	TaskStatusExecuting
	// TaskStatusCompleted resembles the completed state of a task
	TaskStatusCompleted
	// TaskStatusFailed resembles the failed state of a task
	TaskStatusFailed
	// TaskStatusTimeout resembles the timeout state of a task
	TaskStatusTimeout
	// TaskStatusTerminated resembles the terminated state of a task
	TaskStatusTerminated
)

//------------------------------------------------------------------------------

// TaskType resembles the type of a task.
type TaskType int

// Enumeration of possible types of a task.
const (
	// TaskTypeComponent resembles the type of ComponentTask
	TaskTypeComponent TaskType = iota
	// TaskTypeInstance resembles the type of InstanceTask
	TaskTypeInstance
	// TaskTypeTransition resembles the type of TransitionTask
	TaskTypeTransition
	// TaskTypeInstanceState resembles the type of ParallelTask
	TaskTypeParallel
	// TaskTypeInstanceState resembles the type of SequentialTask
	TaskTypeSequential
)

//------------------------------------------------------------------------------

// Task specifies the abstract behaviour of a task
type Task interface {
	Domain() string
	UUID() string
	Parent() string
	Type() TaskType
	Status() TaskStatus
	Phase() int
	GetSubtask(uuid string) (Task, error)
	GetSubtasks() []string
	AddSubtask(subtask Task)
	Save(filename string) error
	Show() (string, error)
	// TODO: marshal and unmarshal
	Execute()
	Terminate()
	Failed()
	Timeout()
	Completed()
}

//------------------------------------------------------------------------------
