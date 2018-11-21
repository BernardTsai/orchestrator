package engine

import (
	"sync"

	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

var eventChannel chan model.Event // the channel for event notification
var eventChannelOnce sync.Once

//------------------------------------------------------------------------------

// GetEventChannel initialises and returns a channel for model.Event objects.
func GetEventChannel() chan model.Event {

	// initialise singleton once
	eventChannelOnce.Do(func() {
		eventChannel = make(chan model.Event)
	})

	return eventChannel
}

//------------------------------------------------------------------------------
