package engine

import (
	"errors"

	"github.com/google/uuid"
	ctrl "tsai.eu/orchestrator/controller"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// InstanceTask evolves an instance towards a desired target state.
type InstanceTask struct {
	domain   string           `yaml:"domain"`   // domain
	uuid     string           `yaml:"uuid"`     // uuid of task
	parent   string           `yaml:"parent"`   // uuid of parent task
	status   model.TaskStatus `yaml:"status"`   // status of task: (execution/completion/failure/timeout/terminated)
	phase    int              `yaml:"phase"`    // internal phase of task
	subtasks []string         `yaml:"subtasks"` // list of subtasks

	component string `yaml:"component"` // component
	version   string `yaml:"version"`   // version of the component
	instance  string `yaml:"instance"`  // uuid of the instance
	state     string `yaml:"state"`     // desired state
}

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, component string, version string, instance string, state string) (*InstanceTask, error) {
	var task InstanceTask

	// TODO: check parameters if context exists
	task.domain = domain
	task.uuid = uuid.New().String()
	task.parent = parent
	task.status = model.TaskStatusInitial
	task.phase = 0
	task.subtasks = []string{}
	task.component = component
	task.version = version
	task.instance = instance
	task.state = state

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
func (task *InstanceTask) UUID() string {
	return task.uuid
}

// Parent delivers the universal unique identifier of the parent task.
func (task *InstanceTask) Parent() string {
	return task.parent
}

// Type delivers the type of the task.
func (task *InstanceTask) Type() model.TaskType {
	return model.TaskTypeComponent
}

// Status delivers the status of the task.
func (task *InstanceTask) Status() model.TaskStatus {
	return task.status
}

// Phase delivers the internal status of the task.
func (task *InstanceTask) Phase() int {
	return task.phase
}

// GetSubtask provides the subtask with a given uuid.
func (task *InstanceTask) GetSubtask(uuid string) (model.Task, error) {
	// find the corresponding subtask
	for _, suuid := range task.subtasks {
		if suuid == uuid {
			// get domain
			domain, err := model.GetModel().GetDomain(task.domain)
			if err != nil {
				return nil, errors.New("unknown domain")
			}

			// get actual subtask
			subtask, err := domain.GetTask(uuid)
			if err != nil {
				return nil, errors.New("unknown subtask")
			}

			return subtask, nil
		}
	}

	// subtask was not found
	return nil, errors.New("unknown subtask")
}

// GetSubtasks provides a slice of subtask uuids.
func (task *InstanceTask) GetSubtasks() []string {
	return task.subtasks
}

// AddSubtask adds a subtask to the list of subtasks.
func (task *InstanceTask) AddSubtask(subtask model.Task) {
	task.subtasks = append(task.subtasks, subtask.UUID())
}

// Execute is the main task execution routine.
func (task *InstanceTask) Execute(channel chan model.Event) error {
	// check status
	status := task.Status()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return errors.New("invalid task state")
	}

	// initialize if needed
	if status == model.TaskStatusInitial {
		// update status
		task.status = model.TaskStatusExecuting
	}

	// TODO: implement and proper error handling
	// - collect relevant information
	// - determine current state and target state of instance and derive the required transition
	// - trigger controller to execute transition
	// - obtain endpoint
	// - calculate endpoint information
	// - trigger new execution
	// - in case of error trigger failure

	// collect relevant information
	domain, _ := model.GetModel().GetDomain(task.domain)
	component, _ := domain.GetComponent(task.component)
	instance, _ := component.GetInstance(task.instance)
	controller, _ := ctrl.GetController(component.Type)
	configuration, _ := model.GetConfiguration(domain.Name, component.Name, instance.UUID)

	// determine current state and target state of instance and derive the required transition
	currentState, _ := controller.Status(configuration)
	targetState := task.state
	transition, err := model.GetTransition(currentState.InstanceState, targetState)

	if err != nil {
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   task.uuid,
			Type:   model.EventTypeTaskFailure,
			Source: task.uuid,
		}
		return errors.New("invalid state")
	}

	// check if reconfiguration is required
	newDependencies := model.DetermineDependencies(domain, component, instance)
	oldDependencies := instance.GetDependencies()

	// execute the required transition
	switch transition {
	case "create":
		instance.SetDependencies(newDependencies)
		_, err := controller.Create(configuration)
		return err
	case "start":
		instance.SetDependencies(newDependencies)
		_, err := controller.Start(configuration)
		return err
	case "stop":
		instance.SetDependencies(newDependencies)
		_, err := controller.Stop(configuration)
		return err
	case "destroy":
		instance.SetDependencies(newDependencies)
		_, err := controller.Destroy(configuration)
		return err
	case "reset":
		instance.SetDependencies(newDependencies)
		_, err := controller.Reset(configuration)
		return err
	case "configure":
		instance.SetDependencies(newDependencies)
		_, err := controller.Configure(configuration)
		return err
	case "none":
		if !util.AreEqual(oldDependencies, newDependencies) {
			_, err := controller.Configure(configuration)
			return err
		}
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Terminate handles the termination of the task
func (task *InstanceTask) Terminate(channel chan model.Event) error {
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
func (task *InstanceTask) Failed(channel chan model.Event) error {
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
func (task *InstanceTask) Timeout(channel chan model.Event) error {
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
func (task *InstanceTask) Completed(channel chan model.Event) error {
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
func (task *InstanceTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

// Show displays the task information as yaml
func (task *InstanceTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
