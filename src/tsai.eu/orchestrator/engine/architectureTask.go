package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// NewArchitectureTask creates a new task
func NewArchitectureTask(domain string, parent string, architecture *model.Architecture) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type = "ArchitectureTask"
	task.Domain = domain
	task.Architecture = architecture.Name
	task.Component = ""
	task.Version = ""
	task.Instance = ""
	task.State = ""
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = []string{}

	// add handlers
	task.SetExecute(ExecuteTask)
	task.SetTerminate(TerminateTask)
	task.SetFailed(FailedTask)
	task.SetTimeout(TimeoutTask)
	task.SetCompleted(CompletedTask)

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

	// construct all required subtasks (one for each service)
	architecture.Services.RLock()
	for service := range architecture.Services.Map {
		subtask, err := NewServiceTask(domain, task.UUID, architecture.Name, service)
		if err != nil {
			return task, errors.New("unable to create subtask for a required service")
		}

		task.AddSubtask(&subtask)
	}
	architecture.Services.RUnlock()

	// success
	return task, nil
}

//------------------------------------------------------------------------------
