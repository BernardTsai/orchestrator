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
	task.domain = domain
	task.uuid = uuid.New().String()
	task.parent = parent
	task.status = model.TaskStatusInitial
	task.phase = 0
	task.subtasks = subtasks

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
func (task SequentialTask) Execute() error {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.Status()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return errors.New("invalid task state")
	}

	// check if the task has finished
	if task.phase >= len(task.subtasks) {
		// update status
		task.status = model.TaskStatusCompleted

		// inform parent
		if task.parent != "" {
			channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskExecution, task.uuid)
		}

		// success
		return nil
	}

	// check status of current subtask
	domain, _ := model.GetModel().GetDomain(task.domain)
	subtask, _ := domain.GetTask(task.subtasks[task.phase])

	switch subtask.Status() {
	// trigger subtask which may not have started yet
	case model.TaskStatusInitial:
		channel <- model.NewEvent(task.domain, subtask.UUID(), model.EventTypeTaskExecution, task.uuid)
	// do nothing if task is still executing
	case model.TaskStatusExecuting:
	// do nothing if subtask has been terminated
	case model.TaskStatusTerminated:
	// proceed to next task if subtask has been completed
	case model.TaskStatusCompleted:
		task.phase++

		channel <- model.NewEvent(task.domain, task.uuid, model.EventTypeTaskExecution, task.uuid)
	// check if subtask has failed
	case model.TaskStatusFailed:
		task.status = model.TaskStatusFailed

		// inform parent
		if task.parent != "" {
			channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskFailure, task.uuid)
		}
	// check if subtask has run into a timeout
	case model.TaskStatusTimeout:
		task.status = model.TaskStatusTimeout

		// inform parent
		if task.parent != "" {
			channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskTimeout, task.uuid)
		}
	}

	// success
	return nil
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
