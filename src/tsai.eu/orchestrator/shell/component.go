package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ComponentCommand executes the component related subcommands
func ComponentCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ComponentUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ComponentUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 2 {
			ComponentUsage(true, context)
			return
		}

		// execute command
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		components, _ := domain.ListComponents()
		result, err := util.ConvertToJSON(components)
		handleResult(context, err, "components could not be listed", result)

	case "create":
		// check availability of arguments
		if len(context.Args) < 4 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new component
		component, err := model.NewComponent(context.Args[2], context.Args[3])
		if err != nil {
			handleResult(context, err, "unable to create a new component", "")
			return
		}

		// add component to domain
		err = domain.AddComponent(component)
		handleResult(context, err, "unable to create component", "component has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 3 {
			ComponentUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new component
		component, _ := model.NewComponent("", "")

		// load component
		err = component.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "component could not be loaded", "")
		}

		// add component to domain
		err = domain.AddComponent(component)
		handleResult(context, err, "unable to load component", "component has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 4 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine component
		component, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// save component
		err = component.Save(context.Args[3])
		handleResult(context, err, "unable to save component", "component has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 3 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine component
		t, err := domain.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute the command
		result, err := t.Show()
		handleResult(context, err, "component can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 3 {
			ComponentUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine component
		_, err = d.GetComponent(context.Args[2])

		if err != nil {
			handleResult(context, err, "component can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteComponent(context.Args[2])
		handleResult(context, err, "component can not be deleted", "component has been deleted")
	default:
		ComponentUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ComponentUsage describes how to make use of the subcommand
func ComponentUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  component list <domain>`)
	context.Println(`            create <domain> <component> <type>`)
	context.Println(`            load <domain> <filename>`)
	context.Println(`            save <domain> <component> <filename>`)
	context.Println(`            show <domain> <component>`)
	context.Println(`            delete <domain> <component>`)
}

//------------------------------------------------------------------------------
