package shell

import (
	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
)

//------------------------------------------------------------------------------

// ModelCommand executes the model related subcommands
func ModelCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		ModelUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		ModelUsage(true, context)
	case "reset":
		err := m.Reset()
		handleResult(context, err, "model could not be reset", "model has been reset")
	case "load":
		err := m.Load(context.Args[1])
		handleResult(context, err, "model could not be loaded", "model has been loaded")
	case "save":
		err := m.Save(context.Args[1])
		handleResult(context, err, "model could not be saved", "model has been saved")
	case "show":
		result, err := m.Show()
		handleResult(context, err, "model can not be displayed", result)
	default:
		ModelUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// ModelUsage describes how to make use of the model subcommand
func ModelUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  model reset`)
	context.Println(`        load <filename>`)
	context.Println(`        save <filename>`)
	context.Println(`        show`)
}

//------------------------------------------------------------------------------
