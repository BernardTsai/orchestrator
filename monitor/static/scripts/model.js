var model;

const loadModel = async () => {
  domainName   = "Test"
  url          = "http://localhost:8081/data";
  response     = await fetch(url);
  text         = await response.text();
  model        = jsyaml.safeLoad(text);

  adjustView(domainName);
}

//------------------------------------------------------------------------------

// constructTree constructs tree hierarchy of components
function constructTree(domain) {

  // loop over all architectures
  for (var architectureName in domain.architectures) {
    view.tree.architectures[architectureName] = {
      data:  domain.architectures[architectureName],
      tasks: {}
    }
  }

  // loop over all components
  for (var componentName in domain.components) {
    view.tree.components[componentName] = {
      data:     domain.components[componentName],
      tasks:    {},
      versions: {}
    }
  }

  // loop over all tasks with references to components
  for (var taskName in domain.tasks) {
    task = domain.tasks[taskName]

    // add component if needed
    if ((task.component != "") && (view.tree.components[task.component] == null)) {
      view.tree.components[task.component] = {
        data:     { name: task.component },
        tasks:    {},
        versions: {}
      }
    }
  }

  // loop over all versions
  for (var componentName in domain.components) {
    component = domain.components[componentName]

    for (var versionName in component.versions) {

      view.tree.components[componentName].versions[versionName] = {
        data:      component.versions[versionName],
        instances: {}
      }
    }
  }

  // loop over all instances
  for (var componentName in domain.components) {
    component = domain.components[componentName]

    for (var versionName in component.versions) {
      version = component.versions[versionName]

      for (var instanceName in version.instances) {

        view.tree.components[componentName].versions[versionName].instances[instanceName] = {
          data:  version.instances[instanceName],
          tasks: {}
        }
      }
    }
  }

  // loop over all tasks
  for (var taskName in domain.tasks) {
    task = domain.tasks[taskName]

    // add user tasks
    if (!task || task.architecture == "") {
      view.tree.tasks[taskName] = {
        data:      domain.tasks[taskName],
        received:  {},
        triggered: {}
      }
      continue
    }

    // add architecture task
    if (task.component == "") {
      view.tree.architectures[task.architecture].tasks[taskName] = {
        data:      domain.tasks[taskName],
        received:  {},
        triggered: {}
      }
      continue
    }

    // add component task
    if (task.version == "") {
      view.tree.components[task.component].tasks[taskName] = {
        data:      domain.tasks[taskName],
        received:  {},
        triggered: {}
      }
      continue
    }

    // add instance task
    if (task.instance != "") {
      view.tree.components[task.component].versions[task.version].instances[task.instance].tasks[taskName] = {
        data:      domain.tasks[taskName],
        received:  {},
        triggered: {}
      }
      continue
    }
  }

  // loop over all events
  for (var eventName in domain.events) {
    event = domain.events[eventName]

    // determine tasks
    task1 = findTask(domain, event.source)
    task2 = findTask(domain, event.task)

    // add events to tasks
    task1.triggered[eventName] = event
    task2.received[eventName]  = event

    // add event to view
    view.events[eventName] = event
  }
}

//------------------------------------------------------------------------------

// findTask determines the task in the tree via the provided name
function findTask(domain, taskName) {
  // determine task data
  task = domain.tasks[taskName]

  //  user tasks
  if (!task || task.architecture == "") {
    return view.tree.tasks[taskName]
  }

  // architecture tasks
  if (task.component == "") {
    return view.tree.architectures[task.architecture].tasks[taskName]
  }

  // component tasks
  if (task.version == "") {
    return view.tree.components[task.component].tasks[taskName]
  }

  // instance tasks
  if (task.instance != "") {
    return view.tree.components[task.component].versions[task.version].instances[task.instance].tasks[taskName]
  }

  console.log("Unknown task: '" + taskName + "'" )
  return null
}

//------------------------------------------------------------------------------

// calculateEventDimensions determines dimensions of all events
function calculateEventDimensions(domain) {

  // loop over all events
  for (var eventName in domain.events) {
    event = domain.events[eventName]

    // update min, max and range
    if (view.min == 0) {
      view.min = event.time
    } else {
      view.min = Math.min(view.min, event.time)
    }
    view.max = Math.max(view.max, event.time)

    view.range = view.max - view.min
  }
}

//------------------------------------------------------------------------------

// calculateTaskDimensions determines dimensions of all tasks
function calculateTaskDimensions(domain) {
  // loop over all tasks
  for (var taskName in domain.tasks) {
    task = findTask(domain, taskName)

    task.x = 0
    task.w = view.max
    task.c = view.max

    //----- loop over all triggered events -----
    for (var eventName in task.triggered) {
      event = task.triggered[eventName]

      if (task.x == 0) {
        x1 = event.time
        x2 = event.time
        x3 = view.max
      } else {
        x1 = Math.min(task.x,        event.time)
        x2 = Math.max(task.x+task.w, event.time)
        x3 = task.c
      }

      task.x = x1
      task.w = x2 - x1
      task.c = x3
    } // triggered loop

    //----- loop over all received events -----
    for (var eventName in task.received) {
      event = task.received[eventName]

      if (event.type == "execution") {
        if (task.x == 0) {
          x1 = event.time
          x2 = event.time
          x3 = view.max
        } else {
          x1 = Math.min(task.x,        event.time)
          x2 = Math.max(task.x+task.w, event.time)
          x3 = Math.min(task.c,        event.time)
        }
      } else {
        if (task.x == 0) {
          x1 = event.time
          x2 = event.time
          x3 = event.time
        } else {
          x1 = Math.min(task.x,        event.time)
          x2 = Math.max(task.x+task.w, event.time)
          x3 = Math.min(task.c,        event.time)
        }
      }

      task.x = x1
      task.w = x2 - x1
      task.c = x3
    } // received loop

  } // task loop

  // loop over all tasks
  for (var taskName in domain.tasks) {
    task = findTask(domain, taskName)

    if (task.x == 0)
    {
      task.x = view.min
      task.w = view.range
      task.c = view.max
    }
  } // task loop
}

//------------------------------------------------------------------------------

// adjustView converts model into the required structure for VUE
function adjustView(domainName) {

  // determine domain
  view.domain = model.domains[domainName]
  if (view.domain == null) {
    console.log("Unknown domain")
    return
  }

  // construct tree
  constructTree(view.domain)

  // calculate dimensions
  calculateEventDimensions(view.domain)
  calculateTaskDimensions(view.domain)
}

//------------------------------------------------------------------------------
