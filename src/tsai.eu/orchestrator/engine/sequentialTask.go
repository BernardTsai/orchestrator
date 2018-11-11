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
	domain   string           `yaml:"domain"`   // domain of task
	uuid     string           `yaml:"uuid"`     // uuid of task
	parent   string           `yaml:"parent"`   // uuid of parent task
	status   model.TaskStatus `yaml:"status"`   // status of task: (execution/completion/failure)
	phase    int              `yaml:"phase"`    // phase of task
	subtasks []string         `yaml:"subtasks"` // list of subtasks
}

// NewSequentialTask creates a new task
func NewSequentialTask(domain string, parent string, subtasks []string) (*SequentialTask, error) {
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
		return nil, errors.New("unknown domain")
	}

	// add task to domain
	err = d.AddTask(&task)
	if err != nil {
		return nil, err
	}

	// success
	return &task, nil
}

// UUID delivers the universal unique identifier of the task.
func (task *SequentialTask) UUID() string {
	return task.uuid
}

// Parent delivers the universal unique identifier of the parent task.
func (task *SequentialTask) Parent() string {
	return task.parent
}

// Type delivers the type of the task.
func (task *SequentialTask) Type() model.TaskType {
	return model.TaskTypeSequential
}

// Status delivers the status of the task.
func (task *SequentialTask) Status() model.TaskStatus {
	return task.status
}

// Phase delivers the internal status of the task.
func (task *SequentialTask) Phase() int {
	return task.phase
}

// GetSubtask provides the subtask with a given uuid.
func (task *SequentialTask) GetSubtask(uuid string) (model.Task, error) {
	// check if uuid is in slice of substasks
	found := false
	for _, suuid := range task.subtasks {
		if suuid == uuid {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("unknown subtask")
	}

	// get domain
	domain, _ := model.GetModel().GetDomain(task.domain)

	// get subtask
	subtask, err := domain.GetTask(uuid)
	if err != nil {
		return nil, errors.New("unknown subtask")
	}

	// success
	return subtask, nil
}

// GetSubtasks provides a slice of subtask uuids.
func (task *SequentialTask) GetSubtasks() []string {
	return task.subtasks
}

// AddSubtask adds a subtask to the list of subtasks.
func (task *SequentialTask) AddSubtask(subtask model.Task) {
	task.subtasks = append(task.subtasks, subtask.UUID())
}

// Execute is the main task execution routine.
func (task *SequentialTask) Execute(channel chan model.Event) error {
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
			channel <- model.Event{
				Domain: task.domain,
				UUID:   uuid.New().String(),
				Task:   task.parent,
				Type:   model.EventTypeTaskExecution,
				Source: task.uuid,
			}
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
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   subtask.UUID(),
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}
	// do nothing if task is still executing
	case model.TaskStatusExecuting:
	// do nothing if subtask has been terminated
	case model.TaskStatusTerminated:
	// proceed to next task if subtask has been completed
	case model.TaskStatusCompleted:
		task.phase++

		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   task.uuid,
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}
	// check if subtask has failed
	case model.TaskStatusFailed:
		task.status = model.TaskStatusFailed

		// inform parent
		if task.parent != "" {
			channel <- model.Event{
				Domain: task.domain,
				UUID:   uuid.New().String(),
				Task:   task.parent,
				Type:   model.EventTypeTaskExecution,
				Source: task.uuid,
			}
		}
	// check if subtask has run into a timeout
	case model.TaskStatusTimeout:
		task.status = model.TaskStatusTimeout

		// inform parent
		if task.parent != "" {
			channel <- model.Event{
				Domain: task.domain,
				UUID:   uuid.New().String(),
				Task:   task.parent,
				Type:   model.EventTypeTaskTimeout,
				Source: task.uuid,
			}
		}
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Terminate handles the termination of the task
func (task *SequentialTask) Terminate(channel chan model.Event) error {
	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusTerminated

		// terminate all subtasks
		for _, subtask := range task.subtasks {
			channel <- model.Event{
				Domain: task.domain,
				UUID:   uuid.New().String(),
				Task:   subtask,
				Type:   model.EventTypeTaskTermination,
				Source: task.uuid,
			}
		}
	}

	// success
	return nil
}

// Failed handles the failure of the task
func (task *SequentialTask) Failed(channel chan model.Event) error {
	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusFailed

		// retrigger execution of parent
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   task.parent,
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}
	}

	// success
	return nil
}

// Timeout handles the timeput of the task
func (task *SequentialTask) Timeout(channel chan model.Event) error {
	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusTimeout

		// retrigger execution of parent
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   task.parent,
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}
	}

	// success
	return nil
}

// Completed handles the completion of the task
func (task *SequentialTask) Completed(channel chan model.Event) error {
	// check if task is regarded to be executing
	if task.status == model.TaskStatusExecuting {
		// update status
		task.status = model.TaskStatusCompleted

		// retrigger execution of parent
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   task.parent,
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}
	}

	// success
	return nil
}

// Save writes the task as json data to a file
func (task *SequentialTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

// Show displays the task information as yaml
func (task *SequentialTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
