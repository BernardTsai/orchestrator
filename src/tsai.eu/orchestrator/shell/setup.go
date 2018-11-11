package shell

import (
	"strconv"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// SetupCommand executes the architecture setup related subcommands
func SetupCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		SetupUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		SetupUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 4 {
			SetupUsage(true, context)
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

		// list setups
		setups, _ := service.ListSetups()
		result, err := util.ConvertToJSON(setups)
		handleResult(context, err, "setups could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 8 {
			SetupUsage(true, context)
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

		// create new setup (name, version, state, size)
		size, err := strconv.Atoi(context.Args[7])
		if err != nil {
			handleResult(context, err, "invalid siez", "")
			return
		}

		setup, err := model.NewSetup(context.Args[4], context.Args[5], context.Args[6], size)
		if err != nil {
			handleResult(context, err, "unable to create a new setup", "")
			return
		}

		// add setup to service
		err = service.AddSetup(setup)
		handleResult(context, err, "unable to create setup", "setup has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 5 {
			SetupUsage(true, context)
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

		// create new setup
		setup, _ := model.NewSetup("", "", "", 0)

		// load version
		err = setup.Load(context.Args[4])

		if err != nil {
			handleResult(context, err, "setup could not be loaded", "")
		}

		// add setup to version
		err = service.AddSetup(setup)
		handleResult(context, err, "unable to load setup", "setup has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 6 {
			SetupUsage(true, context)
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

		// get setup
		setup, err := service.GetSetup(context.Args[4])

		if err != nil {
			handleResult(context, err, "setup can not be identified", "")
			return
		}

		// save setup
		err = setup.Save(context.Args[5])
		handleResult(context, err, "unable to save setup", "setup has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 5 {
			SetupUsage(true, context)
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

		// get setup
		setup, err := service.GetSetup(context.Args[4])

		if err != nil {
			handleResult(context, err, "setup can not be identified", "")
			return
		}

		// execute the command
		result, err := setup.Show()
		handleResult(context, err, "setup can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 5 {
			SetupUsage(true, context)
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

		// execute command
		err = service.DeleteSetup(context.Args[4])
		handleResult(context, err, "setup can not be deleted", "architecture setup has been deleted")
	default:
		SetupUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// SetupUsage describes how to make use of the subcommand
func SetupUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  setup list <domain> <architecture>`)
	context.Println(`             create <domain> <architecture> <service> <setup> <version> <state> <size>`)
	context.Println(`             load <domain> <architecture> <setup> <filename>`)
	context.Println(`             save <domain> <architecture> <service> <setup> <filename>`)
	context.Println(`             show <domain> <architecture> <service> <setup> `)
	context.Println(`             delete <domain> <architecture> <service> <setup> `)
}

//------------------------------------------------------------------------------
