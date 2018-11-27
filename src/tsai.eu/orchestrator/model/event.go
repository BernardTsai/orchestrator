package model

import (
	"errors"

	"github.com/google/uuid"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// EventType resembles the type of an event.
type EventType string

// Enumeration of possible types of an event.
const (
	// EventTypeTaskExecution resembles an event which should trigger the execution of a task.
	EventTypeTaskExecution EventType = "execution"
	// EventTypeTaskCompletion resembles an event which should trigger the closure of a task.
	EventTypeTaskCompletion EventType = "completion"
	// EventTypeTaskFailure resembles an event which should trigger failure handling of a task.
	EventTypeTaskFailure EventType = "failure"
	// EventTypeTaskTimeout resemblesan an event which should trigger timeout handling of a task.
	EventTypeTaskTimeout EventType = "timeout"
	// EventTypeTaskTermination resembles an event which should trigger termination handling of a task.
	EventTypeTaskTermination EventType = "termination"
	// EventTypeTaskUnknown resembles an unknown event.
	EventTypeTaskUnknown EventType = "unknown"
)

// EventType2String converts EventType to a string
func EventType2String(eventType EventType) (string, error) {
	switch eventType {
	case EventTypeTaskExecution:
		return "execution", nil
	case EventTypeTaskCompletion:
		return "completion", nil
	case EventTypeTaskFailure:
		return "failure", nil
	case EventTypeTaskTimeout:
		return "timeout", nil
	case EventTypeTaskTermination:
		return "termination", nil
	}
	return "", errors.New("unknown type")
}

// String2EventType converts a string to an EventType
func String2EventType(eventType string) (EventType, error) {
	switch eventType {
	case "execution":
		return EventTypeTaskExecution, nil
	case "completion":
		return EventTypeTaskCompletion, nil
	case "failure":
		return EventTypeTaskFailure, nil
	case "timeout":
		return EventTypeTaskTimeout, nil
	case "termination":
		return EventTypeTaskTermination, nil
	}
	return EventTypeTaskUnknown, errors.New("unknown type")
}

//------------------------------------------------------------------------------
// Event
// =====
//
// Attributes:
//   - Domain
//   - UUID
//   - Task
//   - type
//   - Source
//
// Functions:
//   - NewEvent
//
//   - event.Show
//   - event.Load
//   - event.Save
//------------------------------------------------------------------------------

// Event describes a situation which may trigger further tasks.
type Event struct {
	Domain string    `yaml:"domain"` // domain of event
	UUID   string    `yaml:"uuid"`   // uuid of event
	Task   string    `yaml:"task"`   // uuid of task
	Type   EventType `yaml:"type"`   // type of event: "execution", "completion", "failure"
	Source string    `yaml:"source"` // source of the event (uuid of the task or "")
}

//------------------------------------------------------------------------------

// NewEvent creates a new event
func NewEvent(domain string, task string, etype EventType, source string) Event {
	var event Event

	event.Domain = domain
	event.UUID = uuid.New().String()
	event.Task = task
	event.Type = etype
	event.Source = source

	// success
	return event
}

//------------------------------------------------------------------------------

// Show displays the event information as json
func (event *Event) Show() (string, error) {
	return util.ConvertToYAML(event)
}

//------------------------------------------------------------------------------

// Save writes the event as json data to a file
func (event *Event) Save(filename string) error {
	return util.SaveYAML(filename, event)
}

//------------------------------------------------------------------------------

// Load reads the event from a file
func (event *Event) Load(filename string) error {
	return util.LoadYAML(filename, event)
}

//------------------------------------------------------------------------------
