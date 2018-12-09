package engine

import (
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// TerminateTask handles the termination of the task
func TerminateTask(task *model.Task) {
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

// FailedTask handles the failure of the task
func FailedTask(task *model.Task) {
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

// ExecuteTask handles the execution of the task
func ExecuteTask(task *model.Task) {
}

//------------------------------------------------------------------------------

// TimeoutTask handles the timeout of the task
func TimeoutTask(task *model.Task) {
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

// CompletedTask handles the completion of the task
func CompletedTask(task *model.Task) {
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
