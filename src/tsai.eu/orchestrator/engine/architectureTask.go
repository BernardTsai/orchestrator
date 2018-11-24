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
func NewArchitectureTask(domain string, parent string, architecture *model.Architecture) (*ArchitectureTask, error) {
	var task ArchitectureTask

	// TODO: check parameters if context exists
	task.domain = domain
	task.uuid = uuid.New().String()
	task.parent = parent
	task.status = model.TaskStatusInitial
	task.phase = 0
	task.subtasks = []string{}

	// get domain
	d, err := model.GetModel().GetDomain(domain)
	if err != nil {
		return nil, errors.New("unknown domain")
	}

	// add task to domain
	err = d.AddTask(&task)
	if err != nil {
		return nil, err
	}

	// construct all required subtasks (one for each service)
	for service := range architecture.Services {
		subtask, err := NewServiceTask(domain, task.uuid, architecture.Name, service)
		if err != nil {
			return nil, errors.New("unable to create subtask for a required service")
		}

		task.AddSubtask(subtask)
	}

	// success
	return &task, nil
}

//------------------------------------------------------------------------------
