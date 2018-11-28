package engine

import (
	"errors"
	"fmt"

	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// AbstractTask is the base class for Task implementations
type AbstractTask struct {
	domain   string           `yaml:"domain"`   // domain of task
	uuid     string           `yaml:"uuid"`     // uuid of task
	parent   string           `yaml:"parent"`   // uuid of parent task
	status   model.TaskStatus `yaml:"status"`   // status of task: (execution/completion/failure)
	phase    int              `yaml:"phase"`    // phase of task
	subtasks []string         `yaml:"subtasks"` // list of subtasks
}

//------------------------------------------------------------------------------

// Domain delivers the domain of the task.
func (task AbstractTask) Domain() string {
	return task.domain
}

// UUID delivers the universal unique identifier of the task.
func (task AbstractTask) UUID() string {
	return task.uuid
}

// Parent delivers the universal unique identifier of the parent task.
func (task AbstractTask) Parent() string {
	return task.parent
}

// Type delivers the type of the task.
func (task AbstractTask) Type() model.TaskType {
	return model.TaskTypeParallel
}

// Status delivers the status of the task.
func (task AbstractTask) Status() model.TaskStatus {
	return task.status
}

// Phase delivers the internal status of the task.
func (task AbstractTask) Phase() int {
	return task.phase
}

// GetSubtask provides the subtask with a given uuid.
func (task AbstractTask) GetSubtask(uuid string) (model.Task, error) {
	// check if uuid is in slice of substasks
	found := false
	for _, suuid := range task.subtasks {
		if suuid == uuid {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("unknown subtask")
	}

	// get domain
	domain, _ := model.GetModel().GetDomain(task.domain)

	// get subtask
	subtask, err := domain.GetTask(uuid)
	if err != nil {
		return nil, errors.New("unknown subtask")
	}

	// success
	return subtask, nil
}

// GetSubtasks provides a slice of subtask uuids.
func (task AbstractTask) GetSubtasks() []string {
	return task.subtasks
}

// AddSubtask adds a subtask to the list of subtasks.
func (task AbstractTask) AddSubtask(subtask model.Task) {
	task.subtasks = append(task.subtasks, subtask.UUID())
}

//------------------------------------------------------------------------------

// Terminate handles the termination of the task
func (task AbstractTask) Terminate() {
	fmt.Println(util.GID())

	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusTerminated

		// terminate all subtasks
		for _, subtask := range task.subtasks {
			channel <- model.NewEvent(task.domain, subtask, model.EventTypeTaskTermination, task.uuid)
		}
	}
}

//------------------------------------------------------------------------------

// Failed handles the failure of the task
func (task AbstractTask) Failed() {
	fmt.Println(util.GID())

	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusFailed

		// retrigger execution of parent
		channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskFailure, task.uuid)
	}
}

//------------------------------------------------------------------------------

// Timeout handles the timeput of the task
func (task AbstractTask) Timeout() {
	fmt.Println(util.GID())

	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusTimeout

		// signal timeout to parent
		channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskTimeout, task.uuid)
	}
}

//------------------------------------------------------------------------------

// Completed handles the completion of the task
func (task AbstractTask) Completed() {
	fmt.Println(util.GID())

	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusCompleted

		// retrigger execution of parent
		channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskExecution, task.uuid)
	}
}

//------------------------------------------------------------------------------
