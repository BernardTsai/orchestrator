package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// NewSequentialTask creates a new task
func NewSequentialTask(domain string, parent string, subtasks []string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type = "SequentialTask"
	task.Domain = domain
	task.Architecture = ""
	task.Component = ""
	task.Version = ""
	task.Instance = ""
	task.State = ""
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = subtasks

	// add handlers
	task.SetExecute(ExecuteSequentialTask)
	task.SetTerminate(TerminateTask)
	task.SetFailed(FailedTask)
	task.SetTimeout(TimeoutTask)
	task.SetCompleted(CompletedTask)

	// get domain
	d, err := model.GetModel().GetDomain(domain)
	if err != nil {
		return task, errors.New("unknown domain")
	}

	// determine parent node
	if parent != "" {
		parentTask, err := d.GetTask(parent)
		if err != nil {
			return task, errors.New("unknown parent")
		}

		// add parent context
		task.Architecture = parentTask.Architecture
		task.Component = parentTask.Component
		task.Version = parentTask.Version
		task.Instance = parentTask.Instance
		task.State = parentTask.State
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

// ExecuteSequentialTask is the main task execution routine.
func ExecuteSequentialTask(task *model.Task) {
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
