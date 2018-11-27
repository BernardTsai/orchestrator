package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// EventCommand executes the event related subcommands
func EventCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		EventUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		EventUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 2 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// list events
		events, _ := domain.ListEvents()
		result, err := util.ConvertToJSON(events)
		handleResult(context, err, "events could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 5 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new event
		eventType, err := model.String2EventType(context.Args[3])
		if err != nil {
			handleResult(context, err, "invalid event type", "")
			return
		}

		event := model.NewEvent(context.Args[1], context.Args[2], eventType, context.Args[4])
		if err != nil {
			handleResult(context, err, "unable to create a new event", "")
			return
		}

		// add event to domain
		err = domain.AddEvent(&event)
		handleResult(context, err, "unable to create event", "event has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 3 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new event
		event := model.NewEvent("", "", model.EventTypeTaskUnknown, "")

		// load event
		err = event.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "event could not be loaded", "")
		}

		// add event to domain
		err = domain.AddEvent(&event)
		handleResult(context, err, "unable to load event", "event has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 4 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine event
		event, err := domain.GetEvent(context.Args[2])

		if err != nil {
			handleResult(context, err, "event can not be identified", "")
			return
		}

		// save event
		err = event.Save(context.Args[3])
		handleResult(context, err, "unable to save event", "event has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 3 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get event
		event, err := domain.GetEvent(context.Args[2])

		if err != nil {
			handleResult(context, err, "event can not be identified", "")
			return
		}

		// execute the command
		result, err := event.Show()
		handleResult(context, err, "event can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 3 {
			EventUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// execute command
		err = domain.DeleteEvent(context.Args[2])
		handleResult(context, err, "event can not be deleted", "event has been deleted")
	default:
		EventUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// EventUsage describes how to make use of the subcommand
func EventUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  event list <domain>`)
	context.Println(`        create <domain> <event> <type> <source>`)
	context.Println(`        load <domain> <filename>`)
	context.Println(`        save <domain> <event> <filename>`)
	context.Println(`        show <domain> <event>`)
	context.Println(`        delete <domain> <event>`)
}

//------------------------------------------------------------------------------
