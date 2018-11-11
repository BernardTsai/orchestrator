package shell

import (
	"strings"

	ishell "gopkg.in/abiosoft/ishell.v2"
	"tsai.eu/orchestrator/model"
	"tsai.eu/orchestrator/util"
)

//------------------------------------------------------------------------------

// Run executes the main shell event loop
func Run(m *model.Model) {
	// create new shell which by default includes 'exit', 'help' and
	// 'clear' commands
	shell := ishell.New()

	// register a function for the "model" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "usage",
		Help: "usage command",
		Func: func(c *ishell.Context) {
			ModelUsage(true, c)
			DomainUsage(false, c)
			TemplateUsage(false, c)
			VariantUsage(false, c)
			DependencyUsage(false, c)
			ArchitectureUsage(false, c)
			ServiceUsage(false, c)
			SetupUsage(false, c)
			ComponentUsage(false, c)
			InstanceUsage(false, c)
			TaskUsage(false, c)
			EventUsage(false, c)
		},
	})

	// register a function for the "model" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "model",
		Help: "model commands",
		Func: func(c *ishell.Context) { ModelCommand(c, m) },
	})

	// register a function for the "domain" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "domain",
		Help: "domain commands",
		Func: func(c *ishell.Context) { DomainCommand(c, m) },
	})

	// register a function for the "template" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "template",
		Help: "template commands",
		Func: func(c *ishell.Context) { TemplateCommand(c, m) },
	})

	// register a function for the "templateVersion" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "variant",
		Help: "variant commands",
		Func: func(c *ishell.Context) { VariantCommand(c, m) },
	})

	// register a function for the "dependency" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "dependency",
		Help: "dependency commands",
		Func: func(c *ishell.Context) { DependencyCommand(c, m) },
	})

	// register a function for the "architecture" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "architecture",
		Help: "architecture commands",
		Func: func(c *ishell.Context) { ArchitectureCommand(c, m) },
	})

	// register a function for the "service" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "service",
		Help: "service commands",
		Func: func(c *ishell.Context) { ServiceCommand(c, m) },
	})

	// register a function for the "setup" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "setup",
		Help: "setup commands",
		Func: func(c *ishell.Context) { SetupCommand(c, m) },
	})

	// register a function for the "component" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "component",
		Help: "component commands",
		Func: func(c *ishell.Context) { ComponentCommand(c, m) },
	})

	// register a function for the "instance" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "instance",
		Help: "instance commands",
		Func: func(c *ishell.Context) { InstanceCommand(c, m) },
	})

	// register a function for the "task" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "task",
		Help: "task commands",
		Func: func(c *ishell.Context) { TaskCommand(c, m) },
	})

	// register a function for the "event" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "event",
		Help: "event commands",
		Func: func(c *ishell.Context) { EventCommand(c, m) },
	})

	// register a function for "#" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "comment",
		Help: "comment",
		Func: func(c *ishell.Context) {
			c.Println(strings.Join(c.Args, " "))
		},
	})

	// run shell
	shell.Run()
}

//------------------------------------------------------------------------------

// handleResult reports error information if present or display success message
func handleResult(context *ishell.Context, err error, fail string, success string) {
	if err != nil {
		if util.Debug() {
			context.Printf("%s\n%+v\n	", fail, err)
		} else {
			context.Printf("%s\n ", fail)
		}
	} else {
		context.Println(success)
	}
}

//------------------------------------------------------------------------------
