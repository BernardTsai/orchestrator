package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ParallelTask executes a set of subtasks in parallel.
type ParallelTask struct {
	domain   string           `yaml:"domain"`   // domain of task
	uuid     string           `yaml:"uuid"`     // uuid of task
	parent   string           `yaml:"parent"`   // uuid of parent task
	status   model.TaskStatus `yaml:"status"`   // status of task: (execution/completion/failure)
	phase    int              `yaml:"phase"`    // phase of task
	subtasks []string         `yaml:"subtasks"` // list of subtasks
}

// NewParallelTask creates a new task
func NewParallelTask(domain string, parent string, subtasks []string) (*ParallelTask, error) {
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
func (task *ParallelTask) UUID() string {
	return task.uuid
}

// Parent delivers the universal unique identifier of the parent task.
func (task *ParallelTask) Parent() string {
	return task.parent
}

// Type delivers the type of the task.
func (task *ParallelTask) Type() model.TaskType {
	return model.TaskTypeParallel
}

// Status delivers the status of the task.
func (task *ParallelTask) Status() model.TaskStatus {
	return task.status
}

// Phase delivers the internal status of the task.
func (task *ParallelTask) Phase() int {
	return task.phase
}

// GetSubtask provides the subtask with a given uuid.
func (task *ParallelTask) GetSubtask(uuid string) (model.Task, error) {
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
func (task *ParallelTask) GetSubtasks() []string {
	return task.subtasks
}

// AddSubtask adds a subtask to the list of subtasks.
func (task *ParallelTask) AddSubtask(subtask model.Task) {
	task.subtasks = append(task.subtasks, subtask.UUID())
}

// Execute is the main task execution routine.
func (task *ParallelTask) Execute(channel chan model.Event) error {
	// check status
	status := task.Status()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return errors.New("invalid task state")
	}

	// initially trigger all subtasks
	if status == model.TaskStatusInitial {
		// update status
		task.status = model.TaskStatusExecuting

		// execute all subtasks
		for _, subtask := range task.subtasks {
			channel <- model.Event{
				Domain: task.domain,
				UUID:   uuid.New().String(),
				Task:   subtask,
				Type:   model.EventTypeTaskExecution,
				Source: task.uuid,
			}
		}

		// success
		return nil
	}

	// check status of currently running subtasks
	completed := 0
	for _, suuid := range task.subtasks {
		domain, _ := model.GetModel().GetDomain(task.domain)
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

			return nil
		}
	}

	// check if task has completed
	if completed == len(task.subtasks) {
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
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Terminate handles the termination of the task
func (task *ParallelTask) Terminate(channel chan model.Event) error {
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
func (task *ParallelTask) Failed(channel chan model.Event) error {
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
func (task *ParallelTask) Timeout(channel chan model.Event) error {
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
func (task *ParallelTask) Completed(channel chan model.Event) error {
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
func (task *ParallelTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

// Show displays the task information as yaml
func (task *ParallelTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
