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

// Execute triggers the execution of the task
func (task ParallelTask) Execute() {
	fmt.Println(util.GID())

	// get event channel
	channel := GetEventChannel()

	// get domain
	domain, err := model.GetModel().GetDomain(task.domain)
	if err != nil {
		fmt.Println("invalid domain")
		channel <- model.NewEvent(task.domain, task.uuid, model.EventTypeTaskFailure, task.uuid)
		return
	}

	// check status
	status := task.Status()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		fmt.Println("invalid task state")
		channel <- model.NewEvent(task.domain, task.uuid, model.EventTypeTaskFailure, task.uuid)
		return
	}

	// initially trigger all subtasks
	if status == model.TaskStatusInitial {
		// update status
		task.status = model.TaskStatusExecuting

		// execute all subtasks
		for _, subtask := range task.subtasks {
			// create event
			channel <- model.NewEvent(task.domain, subtask, model.EventTypeTaskExecution, task.uuid)
		}
	}

	// check status of currently running subtasks
	completed := 0
	for _, suuid := range task.subtasks {
		subtask, _ := domain.GetTask(suuid)

		switch subtask.Status() {
		// do nothing if subtask has not started yet or is still executing
		case model.TaskStatusInitial, model.TaskStatusExecuting:
		// increment counter for completed subtasks
		case model.TaskStatusCompleted:
			completed++
		// check if subtask has failed
		case model.TaskStatusTerminated, model.TaskStatusFailed, model.TaskStatusTimeout:
			task.status = model.TaskStatusFailed
			// inform parent of failure
			if task.parent != "" {
				channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskFailure, task.uuid)
			}

			// trigger closure
			fmt.Println("subtask failed")
			channel <- model.NewEvent(task.domain, task.uuid, model.EventTypeTaskFailure, task.uuid)
			return
		}
	}

	// check if task has completed
	if completed == len(task.subtasks) {
		task.status = model.TaskStatusCompleted
		// retrigger parent execution
		if task.parent != "" {
			channel <- model.NewEvent(task.domain, task.parent, model.EventTypeTaskExecution, task.uuid)
		}

		// trigger clouser
		fmt.Println("completed")
		channel <- model.NewEvent(task.domain, task.uuid, model.EventTypeTaskCompletion, task.uuid)
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
