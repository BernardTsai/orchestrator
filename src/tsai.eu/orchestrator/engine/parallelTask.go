package engine

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ParallelTask executes a set of subtasks in parallel.
type ParallelTask struct {
	AbstractTask
}

//------------------------------------------------------------------------------

// NewParallelTask creates a new task
func NewParallelTask(domain string, parent string, subtasks []string) (ParallelTask, error) {
	var task ParallelTask

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

// Execute triggers the execution of the task
func (task *ParallelTask) Execute() {
	fmt.Println(util.GID())

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

// Save writes the task as json data to a file
func (task ParallelTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task ParallelTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
