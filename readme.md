Quick Start
===========

Prequisites
-----------

The following prerequisites need to be met:

- golang installed (https://golang.org/)
- git installed ()

Installation
------------

A) Retrieve the code

Clone the repository from GitHub:

```
git clone https://github.com/bernardtsai/orchestrator
```

B) Setup the environment

Source the setup script to adjust the required paths for GO and install all
required 3rd party libraries:

```
./setup.sh
```

Usage
-----

Start the orchestrator:

```
./start.sh
```

The main commands can be listed by entering "help" onto the command line:

```
>>> help

Commands:
  architecture      architecture commands
  clear             clear the screen
  comment           comment
  component         component commands
  dependency        dependency commands
  domain            domain commands
  event             event commands
  exit              exit the program
  help              display help
  instance          instance commands
  model             model commands
  service           service commands
  setup             setup commands
  task              task commands
  template          template commands
  usage             usage command
  variant           variant commands
>>>
```

The detailed explanation
