package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// TemplateCommand executes the template related subcommands
func TemplateCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		TemplateUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		TemplateUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 2 {
			TemplateUsage(true, context)
			return
		}

		// execute command
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		templates, _ := domain.ListTemplates()
		result, err := util.ConvertToJSON(templates)
		handleResult(context, err, "templates could not be listed", result)

	case "create":
		// check availability of arguments
		if len(context.Args) < 4 {
			TemplateUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new template
		template, err := model.NewTemplate(context.Args[2], context.Args[3])
		if err != nil {
			handleResult(context, err, "unable to create a new template", "")
			return
		}

		// add template to domain
		err = domain.AddTemplate(template)
		handleResult(context, err, "unable to create template", "template has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 3 {
			TemplateUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// create new template
		template, _ := model.NewTemplate("", "")

		// load template
		err = template.Load(context.Args[2])

		if err != nil {
			handleResult(context, err, "template could not be loaded", "")
		}

		// add template to domain
		err = domain.AddTemplate(template)
		handleResult(context, err, "unable to load template", "template has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 4 {
			TemplateUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine template
		template, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// save template
		err = template.Save(context.Args[3])
		handleResult(context, err, "unable to save template", "template has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 3 {
			TemplateUsage(true, context)
			return
		}

		// determine domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine template
		t, err := domain.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// execute the command
		result, err := t.Show()
		handleResult(context, err, "template can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 3 {
			TemplateUsage(true, context)
			return
		}

		// determine domain
		d, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine template
		_, err = d.GetTemplate(context.Args[2])

		if err != nil {
			handleResult(context, err, "template can not be identified", "")
			return
		}

		// execute command
		err = d.DeleteTemplate(context.Args[2])
		handleResult(context, err, "template can not be deleted", "template has been deleted")
	default:
		TemplateUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// TemplateUsage describes how to make use of the subcommand
func TemplateUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  template list <domain>`)
	context.Println(`           create <domain> <template> <type>`)
	context.Println(`           load <domain> <filename>`)
	context.Println(`           save <domain> <template> <filename>`)
	context.Println(`           show <domain> <template>`)
	context.Println(`           delete <domain> <template>`)
}

//------------------------------------------------------------------------------
