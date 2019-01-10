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

    task.start = view.max + view.max
    task.first = view.max + view.max
    task.last  = 0
    task.end   = 0

    //----- loop over all triggered events -----
    for (var eventName in task.triggered) {
      event = task.triggered[eventName]

      task.first = Math.min(task.first, event.time)
      task.last  = Math.max(task.last,  event.time)
    } // triggered loop

    //----- loop over all received events -----
    for (var eventName in task.received) {
      event = task.received[eventName]

      task.first = Math.min(task.first, event.time)
      task.last  = Math.max(task.last,  event.time)

      if (event.type == "execution") {
        task.start = Math.min(task.start, event.time)
      } else {
        task.end = Math.max(task.end, event.time)
      }
    } // received loop

    // clean up endpoints
    if (task.start > task.first) { task.start = view.min }
    if (task.start > view.max)   { task.start = view.min }

    if (task.last > task.end) { task.end = view.max }
    if (task.last < view.min) { task.end = view.max }
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
