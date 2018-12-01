package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// SequentialTask sequentially executes a set of subtasks.
type SequentialTask struct {
	AbstractTask
}

// NewSequentialTask creates a new task
func NewSequentialTask(domain string, parent string, subtasks []string) (SequentialTask, error) {
	var task SequentialTask

	// TODO: check parameters if context exists
	task.Domain = domain
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = subtasks

	// get domain
	d, err := model.GetModel().GetDomain(domain)
	if err != nil {
		return task, errors.New("unknown domain")
	}

	// add task to domain
	err = d.AddTask(&task)
	if err != nil {
		return task, err
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// Execute is the main task execution routine.
func (task SequentialTask) Execute() {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	// check if the task has finished
	if task.Phase >= len(task.Subtasks) {
		// update status
		task.Status = model.TaskStatusCompleted

		// inform parent
		if task.Parent != "" {
			channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskExecution, task.UUID)
		}

		// success
		return
	}

	// check status of current subtask
	domain, _ := model.GetModel().GetDomain(task.Domain)
	subtask, _ := domain.GetTask(task.Subtasks[task.Phase])

	switch subtask.GetStatus() {
	// trigger subtask which may not have started yet
	case model.TaskStatusInitial:
		channel <- model.NewEvent(task.Domain, subtask.GetUUID(), model.EventTypeTaskExecution, task.UUID)
	// do nothing if task is still executing
	case model.TaskStatusExecuting:
	// do nothing if subtask has been terminated
	case model.TaskStatusTerminated:
	// proceed to next task if subtask has been completed
	case model.TaskStatusCompleted:
		task.Phase++

		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskExecution, task.UUID)
	// check if subtask has failed
	case model.TaskStatusFailed:
		task.Status = model.TaskStatusFailed

		// inform parent
		if task.Parent != "" {
			channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskFailure, task.UUID)
		}
	// check if subtask has run into a timeout
	case model.TaskStatusTimeout:
		task.Status = model.TaskStatusTimeout

		// inform parent
		if task.Parent != "" {
			channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskTimeout, task.UUID)
		}
	}

	// success
	return
}

//------------------------------------------------------------------------------

// Save writes the task as json data to a file
func (task SequentialTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task SequentialTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
