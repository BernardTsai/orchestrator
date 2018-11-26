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
	AbstractTask

	component string `yaml:"component"` // component
	version   string `yaml:"version"`   // version of the component
	instance  string `yaml:"instance"`  // uuid of the instance
	state     string `yaml:"state"`     // desired state
}

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, component string, version string, instance string, state string) (InstanceTask, error) {
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

// Execute is the main task execution routine.
func (task InstanceTask) Execute() error {
	// get event channel
	channel := GetEventChannel()

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

// Save writes the task as json data to a file
func (task InstanceTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task InstanceTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
