package engine

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// NewParallelTask creates a new task
func NewParallelTask(domain string, parent string, subtasks []string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type = "ParallelTask"
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
	task.SetExecute(ExecuteParallelTask)
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

// ExecuteParallelTask triggers the execution of the task
func ExecuteParallelTask(task *model.Task) {
	// get event channel
	channel := GetEventChannel()

	// get domain
	domain, err := model.GetModel().GetDomain(task.Domain)
	if err != nil {
		fmt.Println("invalid domain")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
		return
	}

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		fmt.Println("invalid task state")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
		return
	}

	// initially trigger all subtasks
	if status == model.TaskStatusInitial {
		// update status
		task.Status = model.TaskStatusExecuting

		// execute all subtasks
		for _, subtask := range task.Subtasks {
			// create event
			channel <- model.NewEvent(task.Domain, subtask, model.EventTypeTaskExecution, task.UUID)
		}
	}

	// check status of currently running subtasks
	completed := 0
	for _, suuid := range task.Subtasks {
		subtask, _ := domain.GetTask(suuid)

		switch subtask.GetStatus() {
		// do nothing if subtask has not started yet or is still executing
		case model.TaskStatusInitial, model.TaskStatusExecuting:
		// increment counter for completed subtasks
		case model.TaskStatusCompleted:
			completed++
		// check if subtask has failed
		case model.TaskStatusTerminated, model.TaskStatusFailed, model.TaskStatusTimeout:
			task.Status = model.TaskStatusFailed
			// inform parent of failure
			if task.Parent != "" {
				channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskFailure, task.UUID)
			}

			// trigger closure
			fmt.Println("subtask failed")
			channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
			return
		}
	}

	// check if task has completed
	if completed == len(task.Subtasks) {
		task.Status = model.TaskStatusCompleted
		// retrigger parent execution
		if task.Parent != "" {
			channel <- model.NewEvent(task.Domain, task.Parent, model.EventTypeTaskExecution, task.UUID)
		}

		// trigger clouser
		fmt.Println("completed")
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID)
	}
}

//------------------------------------------------------------------------------
