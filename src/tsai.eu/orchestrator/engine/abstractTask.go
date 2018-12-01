package engine

import (
	"errors"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// AbstractTask is the base class for Task implementations
type AbstractTask struct {
	Domain   string           `yaml:"domain"`   // domain of task
	UUID     string           `yaml:"uuid"`     // uuid of task
	Parent   string           `yaml:"parent"`   // uuid of parent task
	Status   model.TaskStatus `yaml:"status"`   // status of task: (execution/completion/failure)
	Phase    int              `yaml:"phase"`    // phase of task
	Subtasks []string         `yaml:"subtasks"` // list of subtasks
}

//------------------------------------------------------------------------------

// GetDomain delivers the domain of the task.
func (task AbstractTask) GetDomain() string {
	return task.Domain
}

// GetUUID delivers the universal unique identifier of the task.
func (task AbstractTask) GetUUID() string {
	return task.UUID
}

// GetParent delivers the universal unique identifier of the parent task.
func (task AbstractTask) GetParent() string {
	return task.Parent
}

// GetType delivers the type of the task.
func (task AbstractTask) GetType() model.TaskType {
	return model.TaskTypeParallel
}

// GetStatus delivers the status of the task.
func (task AbstractTask) GetStatus() model.TaskStatus {
	return task.Status
}

// GetPhase delivers the internal status of the task.
func (task AbstractTask) GetPhase() int {
	return task.Phase
}

// GetSubtask provides the subtask with a given uuid.
func (task AbstractTask) GetSubtask(uuid string) (model.Task, error) {
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
	domain, _ := model.GetModel().GetDomain(task.Domain)

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
	return task.Subtasks
}

// AddSubtask adds a subtask to the list of subtasks.
func (task AbstractTask) AddSubtask(subtask model.Task) {
	task.Subtasks = append(task.Subtasks, subtask.GetUUID())
}

//------------------------------------------------------------------------------

// Terminate handles the termination of the task
func (task AbstractTask) Terminate() {
	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.Status == model.TaskStatusExecuting {
		// update status
		task.Status = model.TaskStatusTerminated

		// terminate all subtasks
		for _, subtask := range task.Subtasks {
			channel <- model.NewEvent(task.Domain, subtask, model.EventTypeTaskTermination, task.UUID)
		}
	}
}

//------------------------------------------------------------------------------

// Failed handles the failure of the task
func (task AbstractTask) Failed() {
	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.Status == model.TaskStatusExecuting {
		// update status
		task.Status = model.TaskStatusFailed

		// retrigger execution of parent
		channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskFailure, task.UUID)
	}
}

//------------------------------------------------------------------------------

// Timeout handles the timeput of the task
func (task AbstractTask) Timeout() {
	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.Status == model.TaskStatusExecuting {
		// update status
		task.Status = model.TaskStatusTimeout

		// signal timeout to parent
		channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskTimeout, task.UUID)
	}
}

//------------------------------------------------------------------------------

// Completed handles the completion of the task
func (task AbstractTask) Completed() {
	// get event channel
	channel := GetEventChannel()

	// check if task is regarded to be executing
	if task.Status == model.TaskStatusExecuting {
		// update status
		task.Status = model.TaskStatusCompleted

		// retrigger execution of parent
		channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskExecution, task.UUID)
	}
}

//------------------------------------------------------------------------------
