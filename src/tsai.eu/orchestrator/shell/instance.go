package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// InstanceCommand executes the instance related subcommands
func InstanceCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		InstanceUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		InstanceUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 3 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// list instances
		instances, _ := component.ListInstances()
		result, err := util.ConvertToJSON(instances)
		handleResult(context, err, "instances could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 4 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// create new instance
		instance, err := model.NewInstance(context.Args[3])
		if err != nil {
			handleResult(context, err, "unable to create a new instance", "")
			return
		}

		// add instance to component
		err = component.AddInstance(instance)
		handleResult(context, err, "unable to create instance", "instance has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 4 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// create new instance
		instance, _ := model.NewInstance("")

		// load instance
		err = instance.Load(context.Args[3])

		if err != nil {
			handleResult(context, err, "instance could not be loaded", "")
		}

		// add instance to component
		err = component.AddInstance(instance)
		handleResult(context, err, "unable to load instance", "instance has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 5 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// get instance
		instance, err := component.GetInstance(context.Args[3])

		if err != nil {
			handleResult(context, err, "instance can not be identified", "")
			return
		}

		// save instance
		err = instance.Save(context.Args[4])
		handleResult(context, err, "unable to save instance", "instance has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 4 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// get instance
		instance, err := component.GetInstance(context.Args[3])

		if err != nil {
			handleResult(context, err, "instance can not be identified", "")
			return
		}

		// execute the command
		result, err := instance.Show()
		handleResult(context, err, "instance can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 4 {
			InstanceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute command
		err = component.DeleteInstance(context.Args[3])
		handleResult(context, err, "instance can not be deleted", "instance has been deleted")
	default:
		InstanceUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// InstanceUsage describes how to make use of the subcommand
func InstanceUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  instance list <domain> <component>`)
	context.Println(`           create <domain> <component> <instance>`)
	context.Println(`           load <domain> <component> <filename>`)
	context.Println(`           save <domain> <component> <instance> <filename>`)
	context.Println(`           show <domain> <component> <instance>`)
	context.Println(`           delete <domain> <component> <instance>`)
}

//------------------------------------------------------------------------------
