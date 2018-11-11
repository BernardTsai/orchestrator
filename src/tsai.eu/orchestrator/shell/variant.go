package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// VariantCommand executes the variant related subcommands
func VariantCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		VariantUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		VariantUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 3 {
			VariantUsage(true, context)
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

		// list variants
		variants, _ := template.ListVariants()
		result, err := util.ConvertToJSON(variants)
		handleResult(context, err, "variants could not be listed", result)
	case "create":
		// check availability of arguments
		if len(context.Args) < 5 {
			VariantUsage(true, context)
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

		// create new variant
		variant, err := model.NewVariant(context.Args[3], context.Args[4])
		if err != nil {
			handleResult(context, err, "unable to create a new variant", "")
			return
		}

		// add variant to template
		err = template.AddVariant(variant)
		handleResult(context, err, "unable to create variant", "variant has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 4 {
			VariantUsage(true, context)
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

		// create new variant
		variant, _ := model.NewVariant("", "")

		// load variant
		err = variant.Load(context.Args[3])

		if err != nil {
			handleResult(context, err, "variant could not be loaded", "")
		}

		// add variant to template
		err = template.AddVariant(variant)
		handleResult(context, err, "unable to load variant", "variant has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 5 {
			VariantUsage(true, context)
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

		// save variant
		err = variant.Save(context.Args[4])
		handleResult(context, err, "unable to save variant", "variant has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 4 {
			VariantUsage(true, context)
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

		// execute the command
		result, err := variant.Show()
		handleResult(context, err, "variant can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 4 {
			VariantUsage(true, context)
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

		// execute command
		err = template.DeleteVariant(context.Args[3])
		handleResult(context, err, "variant can not be deleted", "variant has been deleted")
	default:
		VariantUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// VariantUsage describes how to make use of the subcommand
func VariantUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  templateVersion list <domain> <template>`)
	context.Println(`                  create <domain> <template> <variant> <configuration>`)
	context.Println(`                  load <domain> <template> <filename>`)
	context.Println(`                  save <domain> <template> <variant> <filename>`)
	context.Println(`                  show <domain> <template> <variant>`)
	context.Println(`                  delete <domain> <template> <variant>`)
}

//------------------------------------------------------------------------------
