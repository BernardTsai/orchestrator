package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// DependencyCommand executes the template dependency related subcommands
func DependencyCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		DependencyUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		DependencyUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 4 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// list dependencies
		dependencies, _ := variant.ListDependencies()
		result, err := util.ConvertToJSON(dependencies)
		handleResult(context, err, "dependencies could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 8 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// create new dependency
		dependency, err := model.NewDependency(context.Args[4], context.Args[5], context.Args[6], context.Args[7])
		if err != nil {
			handleResult(context, err, "unable to create a new dependency", "")
			return
		}

		// add dependency to variant
		err = variant.AddDependency(dependency)
		handleResult(context, err, "unable to create dependency", "dependency has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 5 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// create new dependency
		dependency, _ := model.NewDependency("", "", "", "")

		// load version
		err = dependency.Load(context.Args[4])

		if err != nil {
			handleResult(context, err, "dependency could not be loaded", "")
		}

		// add dependency to version
		err = variant.AddDependency(dependency)
		handleResult(context, err, "unable to load dependency", "dependency has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 6 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// get dependency
		dependency, err := variant.GetDependency(context.Args[4])

		if err != nil {
			handleResult(context, err, "dependency can not be identified", "")
			return
		}

		// save dependency
		err = dependency.Save(context.Args[5])
		handleResult(context, err, "unable to save dependency", "dependency has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 5 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// get dependency
		dependency, err := variant.GetDependency(context.Args[4])

		if err != nil {
			handleResult(context, err, "dependency can not be identified", "")
			return
		}

		// execute the command
		result, err := dependency.Show()
		handleResult(context, err, "dependency can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 5 {
			DependencyUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// get variant
		variant, err := template.GetVariant(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant can not be identified", "")
			return
		}

		// execute command
		err = variant.DeleteDependency(context.Args[4])
		handleResult(context, err, "dependency can not be deleted", "template dependency has been deleted")
	default:
		DependencyUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// DependencyUsage describes how to make use of the subcommand
func DependencyUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  dependency list <domain> <template>`)
	context.Println(`             create <domain> <template> <variant> <dependency> <type> <component> <version>`)
	context.Println(`             load <domain> <template> <dependency> <filename>`)
	context.Println(`             save <domain> <template> <variant> <dependency> <filename>`)
	context.Println(`             show <domain> <template> <variant> <dependency> `)
	context.Println(`             delete <domain> <template> <variant> <dependency> `)
}

//------------------------------------------------------------------------------
