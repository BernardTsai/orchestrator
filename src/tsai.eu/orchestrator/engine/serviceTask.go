package engine

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ServiceSetup captures all required configurations for a service.
type ServiceSetup struct {
	name     string
	versions map[string]VersionSetup
}

// VersionSetup captures all required configurations for a version of a service.
type VersionSetup struct {
	version string
	states  map[string]StateSetup
}

// StateSetup captures the sizing of a version of a service with a specific state.
type StateSetup struct {
	state     string
	instances map[string]string
}

//------------------------------------------------------------------------------

func determineCurrentSetup(domain string, service string) ServiceSetup {
	// create ServiceSetup
	serviceSetup := ServiceSetup{
		name:     service,
		versions: map[string]VersionSetup{},
	}

	// loop over all instances of a component/service
	d, _ := model.GetModel().GetDomain(domain) // domain
	c, _ := d.GetComponent(service)            // component
	l, _ := c.ListInstances()                  // list of instances
	for n := range l {
		u := l[n]                // uuid
		i, _ := c.GetInstance(u) // instance

		// check if version exists
		versionSetup, found := serviceSetup.versions[i.Version]
		if !found {
			versionSetup = VersionSetup{
				version: i.Version,
				states:  map[string]StateSetup{},
			}
		}

		// check if state exists
		stateSetup, found := versionSetup.states[i.State]
		if !found {
			stateSetup = StateSetup{
				state:     i.State,
				instances: map[string]string{},
			}
		}

		// add instance
		stateSetup.instances[i.UUID] = i.UUID
	}

	// success
	return serviceSetup
}

func determineTargetSetup(domain string, architecture string, service string) ServiceSetup {
	// create ServiceSetup
	serviceSetup := ServiceSetup{
		name:     service,
		versions: map[string]VersionSetup{},
	}

	// loop over all instances of a component/service
	d, _ := model.GetModel().GetDomain(domain) // domain
	a, _ := d.GetArchitecture(architecture)    // architecture
	s, _ := a.GetService(service)              // service
	l, _ := s.ListSetups()                     // list of setups
	for i := range l {
		n := l[i]             // setup name
		t, _ := s.GetSetup(n) // setup

		// check if version exists
		versionSetup, found := serviceSetup.versions[t.Version]
		if !found {
			versionSetup = VersionSetup{
				version: t.Version,
				states:  map[string]StateSetup{},
			}
		}

		// check if state exists
		stateSetup, found := versionSetup.states[t.State]
		if !found {
			stateSetup = StateSetup{
				state:     t.State,
				instances: map[string]string{},
			}
		}

		// add instances
		for j := 0; j < t.Size; j++ {
			u := uuid.New().String()
			stateSetup.instances[u] = u
		}
	}

	// success
	return serviceSetup
}

func determineTasks(domain string, architecture string, service string) ([]InstanceTask, []InstanceTask, []InstanceTask) {
	targetSetup := determineTargetSetup(domain, architecture, service)
	currentSetup := determineCurrentSetup(domain, service)
	updateTasks := []InstanceTask{}
	createTasks := []InstanceTask{}
	removeTasks := []InstanceTask{}

	// determine all unchanged instances
	for _, targetVersionSetup := range targetSetup.versions {
		for _, targetStateSetup := range targetVersionSetup.states {
			for targetInstance := range targetStateSetup.instances {

				// try to find matching current instance
				currentVersionSetup, found := currentSetup.versions[targetVersionSetup.version]
				if !found {
					continue
				}

				currentStateSetup, found := currentVersionSetup.states[targetStateSetup.state]
				if !found {
					continue
				}

				for currentInstance := range currentStateSetup.instances {
					// instance has been found - now remove instances from the setup
					delete(targetStateSetup.instances, targetInstance)
					delete(currentStateSetup.instances, currentInstance)
					break
				}
			}
		}
	}

	// determine all instances which need to be updated
	for targetVersion, targetVersionSetup := range targetSetup.versions {
		for targetState, targetStateSetup := range targetVersionSetup.states {
			for targetInstance := range targetStateSetup.instances {

				// try to find matching current instance with matching version
				currentVersionSetup, found := currentSetup.versions[targetVersionSetup.version]
				if !found {
					continue
				}

				for _, currentStateSetup := range currentVersionSetup.states {
					for currentInstance := range currentStateSetup.instances {
						// create new update task
						updateTask, _ := NewInstanceTask(domain, "TODO: unknown", service, targetVersion, currentInstance, targetState)

						// append new task to set of update tasks
						updateTasks = append(updateTasks, updateTask)

						// instance has been found - now remove instances from the setup
						delete(targetStateSetup.instances, targetInstance)
						delete(currentStateSetup.instances, currentInstance)
						break
					}

				}
			}
		}
	}

	// all leftover current instances need to be removed
	for currentVersion, currentVersionSetup := range currentSetup.versions {
		for _, currentStateSetup := range currentVersionSetup.states {
			for currentInstance := range currentStateSetup.instances {
				// create new remove task
				removeTask, _ := NewInstanceTask(domain, "TODO: unknown", service, currentVersion, currentInstance, "initial")

				// append new task to set of remove tasks
				removeTasks = append(removeTasks, removeTask)
			}
		}
	}

	// all leftover target instances need to be created
	for targetVersion, targetVersionSetup := range targetSetup.versions {
		for targetState, targetStateSetup := range targetVersionSetup.states {
			for targetInstance := range targetStateSetup.instances {
				// create new create task
				createTask, _ := NewInstanceTask(domain, "TODO: unknown", service, targetVersion, targetInstance, targetState)

				// append new task to set of create tasks
				createTasks = append(createTasks, createTask)
			}
		}
	}

	// success
	return updateTasks, createTasks, removeTasks
}

//------------------------------------------------------------------------------

// ServiceTask evolves a component with its instances towards a desired service setup.
type ServiceTask struct {
	AbstractTask

	architecture string `yaml:"architecture"` // name of architecture
	service      string `yaml:"service"`      // name of the service to be instantiated
}

// NewServiceTask creates a new task
func NewServiceTask(domain string, parent string, architecture string, service string) (ServiceTask, error) {
	var task ServiceTask

	// TODO: check parameters if context exists
	task.domain = domain
	task.uuid = uuid.New().String()
	task.parent = parent
	task.status = model.TaskStatusInitial
	task.phase = 0
	task.subtasks = []string{}
	task.architecture = architecture
	task.service = service

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
func (task ServiceTask) Execute() error {
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

		// determine required subtasks
		updateTasks, createTasks, removeTasks := determineTasks(task.domain, task.architecture, task.service)

		// add tasks to domain

		// create task groups
		mainTask, _ := NewParallelTask(task.domain, task.uuid, []string{})
		task.AddSubtask(mainTask)

		// add update subtasks
		updateTask, _ := NewParallelTask(task.domain, mainTask.UUID(), []string{})
		mainTask.AddSubtask(updateTask)
		for _, s := range updateTasks {
			subTask, _ := NewInstanceTask(s.domain, mainTask.UUID(), s.component, s.version, s.instance, s.state)

			updateTask.AddSubtask(subTask)
		}

		// add create subtasks
		createTask, _ := NewParallelTask(task.domain, mainTask.UUID(), []string{})
		mainTask.AddSubtask(createTask)
		for _, s := range createTasks {
			subTask, _ := NewInstanceTask(s.domain, mainTask.UUID(), s.component, s.version, s.instance, s.state)

			createTask.AddSubtask(subTask)
		}

		// add remove subtasks
		removeTask, _ := NewParallelTask(task.domain, mainTask.UUID(), []string{})
		mainTask.AddSubtask(removeTask)
		for _, s := range removeTasks {
			subTask, _ := NewInstanceTask(s.domain, mainTask.UUID(), s.component, s.version, s.instance, s.state)

			removeTask.AddSubtask(subTask)
		}

		// trigger execution of main subtask
		channel <- model.Event{
			Domain: task.domain,
			UUID:   uuid.New().String(),
			Task:   mainTask.UUID(),
			Type:   model.EventTypeTaskExecution,
			Source: task.uuid,
		}

		// success
		return nil
	}

	// success
	return nil
}

//------------------------------------------------------------------------------

// Save writes the task as json data to a file
func (task ServiceTask) Save(filename string) error {
	return util.SaveYAML(filename, task)
}

//------------------------------------------------------------------------------

// Show displays the task information as yaml
func (task ServiceTask) Show() (string, error) {
	return util.ConvertToYAML(task)
}

//------------------------------------------------------------------------------
