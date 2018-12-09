package engine

import (
	"errors"

	"github.com/google/uuid"
	ctrl "tsai.eu/orchestrator/controller"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// NewInstanceTask creates a new instance task
func NewInstanceTask(domain string, parent string, architecture string, component string, version string, instance string, state string) (model.Task, error) {
	var task model.Task

	// TODO: check parameters if context exists
	task.Type = "InstanceTask"
	task.Domain = domain
	task.Architecture = architecture
	task.Component = component
	task.Version = version
	task.Instance = instance
	task.State = state
	task.UUID = uuid.New().String()
	task.Parent = parent
	task.Status = model.TaskStatusInitial
	task.Phase = 0
	task.Subtasks = []string{}

	// add handlers
	task.SetExecute(ExecuteInstanceTask)
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

	// success
	return task, nil
}

//------------------------------------------------------------------------------

// ExecuteInstanceTask is the main task execution routine.
func ExecuteInstanceTask(task *model.Task) {
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
