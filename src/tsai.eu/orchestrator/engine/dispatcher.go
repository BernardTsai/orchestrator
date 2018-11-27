package engine

import (
	"fmt"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// Dispatcher receives events from a channel and triggers a task coroutine.
type Dispatcher struct {
	Model   *model.Model           // repository
	Channel chan model.Event       // the channel for event notification
	Tasks   map[string]*model.Task // map of running tasks
}

//------------------------------------------------------------------------------

// StartDispatcher creates a dispatcher and returns a channel for new tasks.
func StartDispatcher(m *model.Model) chan model.Event {

	// create the communication channel
	channel := GetEventChannel()

	// create the dispatcher
	dispatcher := Dispatcher{
		Model:   m,
		Channel: channel,
		Tasks:   map[string]*model.Task{},
	}

	// start the dispatcher
	go dispatcher.Run()

	return channel
}

//------------------------------------------------------------------------------

// Run starts the dispatcher loop receiving events and triggering tasks.
func (d *Dispatcher) Run() {
	// loop until exit is requested
	for {
		// get next event
		event := <-d.Channel

		// terminate if domain is empty = exit request
		if event.Domain == "" {
			return
		}

		// get corresponding domain from the model
		domain, err := d.Model.GetDomain(event.Domain)
		if err != nil {
			// TODO: log unknown domain
			continue
		}

		// save event
		domain.AddEvent(&event)

		// get task
		task, err := domain.GetTask(event.Task)
		if err != nil {
			// TODO: log unknown task
			continue
		}

		// determine action by type of event
		// Event types: execute, completed, failed, timeout, terminate
		// Task types can be:
		// - set component state
		// - set instance state
		// - transition component
		// - transition instance
		// - parallel execute tasks
		// - sequentially execute tasks
		switch event.Type {
		// execute the task
		case model.EventTypeTaskExecution:
			err := task.Execute()
			if err != nil {
				fmt.Println(err)
				// TODO: handle error
			}

		// handle task completion
		case model.EventTypeTaskCompletion:
			err := task.Completed()
			if err != nil {
				// TODO: handle error
			}

		// handle task failure
		case model.EventTypeTaskFailure:
			err := task.Failed()
			if err != nil {
				// TODO: handle error
			}

		// handle timeout of a task
		case model.EventTypeTaskTimeout:
			err := task.Timeout()
			if err != nil {
				// TODO: handle error
			}

		// handle termination of a task
		case model.EventTypeTaskTermination:
			err := task.Terminate()
			if err != nil {
				// TODO: handle error
			}

		}
	}
}

//------------------------------------------------------------------------------
