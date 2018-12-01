package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// ArchitectureTask instantiates an architecture definition.
type ArchitectureTask struct {
	ParallelTask
}

// NewArchitectureTask creates a new task
func NewArchitectureTask(domain string, parent string, architecture *model.Architecture) (ArchitectureTask, error) {
	var task ArchitectureTask

	// TODO: check parameters if context exists
	task.Domain = domain
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = []string{}

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
	for service := range architecture.Services {
		subtask, err := NewServiceTask(domain, task.UUID, architecture.Name, service)
		if err != nil {
			return task, errors.New("unable to create subtask for a required service")
		}

		task.AddSubtask(&subtask)
	}

	// success
	return task, nil
}

//------------------------------------------------------------------------------
