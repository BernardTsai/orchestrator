package model

import (
	"errors"

	"tsai.eu/orchestrator/util"
)

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

// TaskHandler is function capable of processing a task related event.
type TaskHandler func(task *Task)

//------------------------------------------------------------------------------

// Task specifies the basic behaviour of a task
type Task struct {
	Type         string     `yaml:"type"`         // type of task
	Domain       string     `yaml:"domain"`       // domain of task
	Architecture string     `yaml:"architecture"` // architecture of entity
	Component    string     `yaml:"component"`    // component of entity
	Version      string     `yaml:"version"`      // version of entity
	Instance     string     `yaml:"instance"`     // instance of entity
	State        string     `yaml:"state"`        // desired state of entity
	UUID         string     `yaml:"uuid"`         // uuid of task
	Parent       string     `yaml:"parent"`       // uuid of parent task
	Status       TaskStatus `yaml:"status"`       // status of task: (execution/completion/failure)
	Phase        int        `yaml:"phase"`        // phase of task
	Subtasks     []string   `yaml:"subtasks"`     // list of subtasks
	execute      TaskHandler
	terminate    TaskHandler
	failed       TaskHandler
	timeout      TaskHandler
	completed    TaskHandler
}

//------------------------------------------------------------------------------

// GetType delivers the type of the task.
func (task *Task) GetType() string {
	return task.Type
}

//------------------------------------------------------------------------------

// GetDomain delivers the domain of the task.
func (task *Task) GetDomain() string {
	return task.Domain
}

//------------------------------------------------------------------------------

// GetArchitecture delivers the architecture of the entity.
func (task *Task) GetArchitecture() string {
	return task.Architecture
}

//------------------------------------------------------------------------------

// GetComponent delivers the component of the entity.
func (task *Task) GetComponent() string {
	return task.Component
}

//------------------------------------------------------------------------------

// GetVersion delivers the version of the entity.
func (task *Task) GetVersion() string {
	return task.Version
}

//------------------------------------------------------------------------------

// GetInstance delivers the instance of entity task.
func (task *Task) GetInstance() string {
	return task.Instance
}

//------------------------------------------------------------------------------

// GetState delivers the state of the entity.
func (task *Task) GetState() string {
	return task.State
}

//------------------------------------------------------------------------------

// GetUUID delivers the universal unique identifier of the task.
func (task *Task) GetUUID() string {
	return task.UUID
}

//------------------------------------------------------------------------------

// GetParent delivers the universal unique identifier of the parent task.
func (task *Task) GetParent() string {
	return task.Parent
}

//------------------------------------------------------------------------------

// GetStatus delivers the status of the task.
func (task *Task) GetStatus() TaskStatus {
	return task.Status
}

//------------------------------------------------------------------------------

// GetPhase delivers the internal status of the task.
func (task *Task) GetPhase() int {
	return task.Phase
}

//------------------------------------------------------------------------------

// GetSubtask provides the subtask with a given uuid.
func (task *Task) GetSubtask(uuid string) (*Task, error) {
	// check if uuid is in slice of substasks
	found := false
	for _, suuid := range task.Subtasks {
		if suuid == uuid {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("unknown subtask")
	}

	// get domain
	domain, _ := GetModel().GetDomain(task.Domain)

	// get subtask
	subtask, err := domain.GetTask(uuid)
	if err != nil {
		return nil, errors.New("unknown subtask")
	}

	// success
	return subtask, nil
}

//------------------------------------------------------------------------------

// GetSubtasks provides a slice of subtask uuids.
func (task *Task) GetSubtasks() []string {
	return task.Subtasks
}

//------------------------------------------------------------------------------

// AddSubtask adds a subtask to the list of subtasks.
func (task *Task) AddSubtask(subtask *Task) {
	task.Subtasks = append(task.Subtasks, subtask.GetUUID())
}

//------------------------------------------------------------------------------

// Save writes the task as json data to a file
func (task *Task) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task *Task) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------

// SetExecute defines the execute event handler of a task
func (task *Task) SetExecute(handler TaskHandler) {
	task.execute = handler
}

//------------------------------------------------------------------------------

// Execute processes a task
func (task *Task) Execute() {
	// execute task if appropriate handler has been defined
	if task.execute != nil {
		task.execute(task)
	}
}

//------------------------------------------------------------------------------

// SetTerminate defines the terminate event handler of a task
func (task *Task) SetTerminate(handler TaskHandler) {
	task.terminate = handler
}

//------------------------------------------------------------------------------

// Terminate stops a task
func (task *Task) Terminate() {
	// execute task if appropriate handler has been defined
	if task.terminate != nil {
		task.terminate(task)
	}
}

//------------------------------------------------------------------------------

// SetFailed defines the failed event handler of a task
func (task *Task) SetFailed(handler TaskHandler) {
	task.failed = handler
}

//------------------------------------------------------------------------------

// Failed handles the failure of the task
func (task *Task) Failed() {
	// execute task if appropriate handler has been defined
	if task.failed != nil {
		task.failed(task)
	}
}

//------------------------------------------------------------------------------

// SetTimeout defines the timeout event handler of a task
func (task *Task) SetTimeout(handler TaskHandler) {
	task.timeout = handler
}

//------------------------------------------------------------------------------

// Timeout handles the timeput of the task
func (task *Task) Timeout() {
	// execute task if appropriate handler has been defined
	if task.timeout != nil {
		task.timeout(task)
	}
}

//------------------------------------------------------------------------------

// SetCompleted defines the completed event handler of a task
func (task *Task) SetCompleted(handler TaskHandler) {
	task.completed = handler
}

//------------------------------------------------------------------------------

// Completed handles the completion of the task
func (task *Task) Completed() {
	// execute task if appropriate handler has been defined
	if task.completed != nil {
		task.completed(task)
	}
}

//------------------------------------------------------------------------------
