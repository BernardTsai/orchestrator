package shell

import (
	"errors"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// TaskCommand executes the task related subcommands
func TaskCommand(context *ishell.Context, m *model.Model) {
	// check if the action has been defined
	if len(context.Args) < 1 {
		TaskUsage(true, context)
		return
	}

	// determine the required action
	action := context.Args[0]

	// handle required action
	switch action {
	case "?":
		TaskUsage(true, context)
	case "list":
		// check availability of arguments
		if len(context.Args) != 2 {
			TaskUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// list tasks
		tasks, _ := domain.ListTasks()
		result, err := util.ConvertToJSON(tasks)
		handleResult(context, err, "tasks could not be listed", result)

	case "create":
		// check availability of arguments
		if len(context.Args) < 4 {
			TaskUsage(true, context)
			return
		}

		// get domain
		_, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// TODO: create new task and add it to domain

		// add template to domain
		err = errors.New("not implemented")
		handleResult(context, err, "unable to create task", "task has been created")
	case "load":
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// get domain
		_, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// TODO: create new task, load from file and add it to domain

		// add template to domain
		err = errors.New("not implemented")
		handleResult(context, err, "unable to load task", "task has been loaded")
	case "save":
		// check availability of arguments
		if len(context.Args) != 4 {
			TaskUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// determine task
		task, err := domain.GetTask(context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// save template
		err = task.Save(context.Args[3])
		handleResult(context, err, "unable to save template", "template has been saved")
	case "show":
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// get task
		task, err := domain.GetTask(context.Args[2])

		if err != nil {
			handleResult(context, err, "task can not be identified", "")
			return
		}

		// execute the command
		result, err := task.Show()
		handleResult(context, err, "task can not be displayed", result)
	case "delete":
		// check availability of arguments
		if len(context.Args) != 3 {
			TaskUsage(true, context)
			return
		}

		// get domain
		domain, err := m.GetDomain(context.Args[1])

		if err != nil {
			handleResult(context, err, "domain can not be identified", "")
			return
		}

		// execute command
		err = domain.DeleteTask(context.Args[2])
		handleResult(context, err, "task can not be deleted", "task has been deleted")
	default:
		TaskUsage(true, context)
	}
}

//------------------------------------------------------------------------------

// TaskUsage describes how to make use of the subcommand
func TaskUsage(header bool, context *ishell.Context) {
	if header {
		context.Println("usage:")
	}
	context.Println(`  task list <domain>`)
	context.Println(`       create <domain> <task> ... (tbd)`)
	context.Println(`       load <domain> <filename>`)
	context.Println(`       save <domain> <task> <filename>`)
	context.Println(`       show <domain> <task>`)
	context.Println(`       delete <domain> <task>`)
}

//------------------------------------------------------------------------------
