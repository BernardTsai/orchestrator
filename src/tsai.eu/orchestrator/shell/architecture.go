package shell

import (
	"fmt"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/engine"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ArchitectureCommand executes the architecture related subcommands
func ArchitectureCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ArchitectureUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ArchitectureUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 2 {
			ArchitectureUsage(true, context)
			return
		}

		// execute command
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		architectures, _ := domain.ListArchitectures()
		result, err := util.ConvertToYAML(architectures)
		handleResult(context, err, "architectures could not be listed", result)

	case "create":
		// check availability of arguments
		if len(context.Args) < 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new architecture
		architecture, err := model.NewArchitecture(context.Args[2])
		if err != nil {
			handleResult(context, err, "unable to create a new architecture", "")
			return
		}

		// add architecture to domain
		err = domain.AddArchitecture(architecture)
		handleResult(context, err, "unable to create architecture", "architecture has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new architecture
		architecture, _ := model.NewArchitecture("")

		// load architecture
		err = architecture.Load(context.Args[2])

		if err != nil {
			fmt.Println(err)
			handleResult(context, err, "architecture could not be loaded", "")
		}

		// add architecture to domain
		err = domain.AddArchitecture(architecture)
		handleResult(context, err, "unable to load architecture", "architecture has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 4 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// save architecture
		err = architecture.Save(context.Args[3])
		handleResult(context, err, "unable to save architecture", "architecture has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		t, err := d.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// execute the command
		result, err := t.Show()
		handleResult(context, err, "architecture can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		_, err = domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// execute command
		err = domain.DeleteArchitecture(context.Args[2])
		handleResult(context, err, "architecture can not be deleted", "architecture has been deleted")
	case "execute":
		// check availability of arguments
		if len(context.Args) != 3 {
			ArchitectureUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// create task and start it by signalling an event
		task, _ := engine.NewArchitectureTask(domain.Name, "", architecture)
		if err != nil {
			handleResult(context, err, "task can not be created", "")
			return
		}

		// get event channel
		channel := engine.GetEventChannel()

		// create event
		channel <- model.NewEvent(domain.Name, task.GetUUID(), model.EventTypeTaskExecution, "")

		handleResult(context, nil, "architecture can not be executed", "architecture execution has been initiated")
	default:
		ArchitectureUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ArchitectureUsage describes how to make use of the subcommand
func ArchitectureUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  architecture list <domain>`)
	context.Println(`               create <domain> <architecture>`)
	context.Println(`               load <domain> <filename>`)
	context.Println(`               save <domain> <architecture> <filename>`)
	context.Println(`               show <domain> <architecture>`)
	context.Println(`               delete <domain> <architecture>`)
	context.Println(`               execute <domain> <architecture>`)
}

//------------------------------------------------------------------------------
