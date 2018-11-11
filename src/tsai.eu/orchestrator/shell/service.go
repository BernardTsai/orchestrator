package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// ServiceCommand executes the service related subcommands
func ServiceCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ServiceUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ServiceUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 3 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// list services
		services, _ := architecture.ListServices()
		result, err := util.ConvertToJSON(services)
		handleResult(context, err, "services could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 4 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// create new service
		service, err := model.NewService(context.Args[3])
		if err != nil {
			handleResult(context, err, "unable to create a new service", "")
			return
		}

		// add service to architecture
		err = architecture.AddService(service)
		handleResult(context, err, "unable to create service", "service has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 4 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// create new service
		service, _ := model.NewService("")

		// load service
		err = service.Load(context.Args[3])

		if err != nil {
			handleResult(context, err, "service could not be loaded", "")
		}

		// add service to architecture
		err = architecture.AddService(service)
		handleResult(context, err, "unable to load service", "service has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 5 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// get service
		service, err := architecture.GetService(context.Args[3])

		if err != nil {
			handleResult(context, err, "service can not be identified", "")
			return
		}

		// save service
		err = service.Save(context.Args[4])
		handleResult(context, err, "unable to save service", "service has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 4 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// get service
		service, err := architecture.GetService(context.Args[3])

		if err != nil {
			handleResult(context, err, "service can not be identified", "")
			return
		}

		// execute the command
		result, err := service.Show()
		handleResult(context, err, "service can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 4 {
			ServiceUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get architecture
		architecture, err := domain.GetArchitecture(context.Args[2])

		if err != nil {
			handleResult(context, err, "architecture can not be identified", "")
			return
		}

		// execute command
		err = architecture.DeleteService(context.Args[3])
		handleResult(context, err, "service can not be deleted", "service has been deleted")
	default:
		ServiceUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ServiceUsage describes how to make use of the subcommand
func ServiceUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  variant list <domain> <architecture>`)
	context.Println(`          create <domain> <architecture> <service>`)
	context.Println(`          load <domain> <architecture> <filename>`)
	context.Println(`          save <domain> <architecture> <service> <filename>`)
	context.Println(`          show <domain> <architecture> <service>`)
	context.Println(`          delete <domain> <architecture> <service>`)
}

//------------------------------------------------------------------------------
