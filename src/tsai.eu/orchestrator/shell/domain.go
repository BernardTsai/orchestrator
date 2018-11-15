package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// DomainCommand executes the domain related subcommands
func DomainCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		DomainUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		DomainUsage(true, context)
	case "list":
		domains, _ := m.ListDomains()
		result, err := util.ConvertToYAML(domains)
		handleResult(context, err, "domains could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		domain, _ := model.NewDomain(context.Args[1])
		err := m.AddDomain(domain)
		handleResult(context, err, "unable to create domain", "domain has been created")
	case "show":
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be displayed", "")
			return
		}

		result, err := d.Show()
		handleResult(context, err, "domain can not be displayed", result)
	case "load":
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		err := m.Load(context.Args[1])
		handleResult(context, err, "domain could not be loaded", "domain has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) < 3 {
			DomainUsage(true, context)
			return
		}

		// execute command
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "unknown domain", "")
			return
		}

		err = d.Save(context.Args[2])
		handleResult(context, err, "domain could not be saved", "domain has been saved")
	case "delete":
		// check availability of arguments
		if len(context.Args) < 2 {
			DomainUsage(true, context)
			return
		}

		// execute command
		err := m.DeleteDomain(context.Args[1])
		handleResult(context, err, "domain can not be deleted", "domain has been deleted")
	default:
		DomainUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// DomainUsage describes how to make use of the domain subcommand
func DomainUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println("  domain list")
	context.Println("         create <domain>")
	context.Println("         show <domain>")
	context.Println("         load <domain> <filename>")
	context.Println("         save <domain> <filename>")
	context.Println("         delete <domain>")
}

//------------------------------------------------------------------------------
