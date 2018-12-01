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

	Component string `yaml:"component"` // component
	Version   string `yaml:"version"`   // version of the component
	Instance  string `yaml:"instance"`  // uuid of the instance
	State     string `yaml:"state"`     // desired state
}

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, component string, version string, instance string, state string) (InstanceTask, error) {
	var task InstanceTask

	// TODO: check parameters if context exists
	task.Domain = domain
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = []string{}
	task.Component = component
	task.Version = version
	task.Instance = instance
	task.State = state

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
func (task InstanceTask) Execute() {
	// get event channel
	channel := GetEventChannel()

	// check status
	status := task.GetStatus()

	if status != model.TaskStatusInitial && status != model.TaskStatusExecuting {
		return
	}

	// initialize if needed
	if status == model.TaskStatusInitial {
		// update status
		task.Status = model.TaskStatusExecuting
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
	domain, _ := model.GetModel().GetDomain(task.Domain)
	component, _ := domain.GetComponent(task.Component)
	instance, _ := component.GetInstance(task.Instance)
	controller, _ := ctrl.GetController(component.Type)
	configuration, _ := model.GetConfiguration(domain.Name, component.Name, instance.UUID)

	// determine current state and target state of instance and derive the required transition
	currentState, _ := controller.Status(configuration)
	targetState := task.State
	transition, err := model.GetTransition(currentState.InstanceState, targetState)

	// check for invalid states
	if err != nil {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
	}

	// check if reconfiguration is required
	newDependencies := model.DetermineDependencies(domain, component, instance)
	oldDependencies := instance.GetDependencies()

	// execute the required transition
	switch transition {
	case "create":
		instance.SetDependencies(newDependencies)
		_, err = controller.Create(configuration)
	case "start":
		instance.SetDependencies(newDependencies)
		_, err = controller.Start(configuration)
	case "stop":
		instance.SetDependencies(newDependencies)
		_, err = controller.Stop(configuration)
	case "destroy":
		instance.SetDependencies(newDependencies)
		_, err = controller.Destroy(configuration)
	case "reset":
		instance.SetDependencies(newDependencies)
		_, err = controller.Reset(configuration)
	case "configure":
		instance.SetDependencies(newDependencies)
		_, err = controller.Configure(configuration)
	case "none":
		if !util.AreEqual(oldDependencies, newDependencies) {
			_, err = controller.Configure(configuration)
		}
	}

	// check for errors
	if err != nil {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskFailure, task.UUID)
	} else {
		channel <- model.NewEvent(task.Domain, task.UUID, model.EventTypeTaskCompletion, task.UUID)
	}
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
