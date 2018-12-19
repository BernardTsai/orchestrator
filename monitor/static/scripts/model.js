var model;

const loadModel = async () => {
  domain       = "Test"
  architecture = "architecture_1.0.0"
  url          = "http://localhost:8081/data";
  response     = await fetch(url);
  text         = await response.text();
  yaml         = jsyaml.safeLoad(text);
  console.log(yaml)
  model        = convertData(yaml, domain, architecture);
}

//------------------------------------------------------------------------------

// convertData converts yaml into the required model structure for VUE
function convertData(yaml, domainName, architectureName) {
  var result

  // add children to top level domain node
  domain = yaml.domains[domainName]
  if (domain != null) {

    architecture = domain.architectures[architectureName]
    if (architecture != null) {
      result = convertData2(domain, architecture)
    }
  }

  // success
  return result;
}

//------------------------------------------------------------------------------

// convertData2 converts yaml into the required model structure for VUE
function convertData2(domain, architecture) {
  // initialize result
  result = {
    architecture: {x:0, y:0, w:0, h:0 },
    components:   {},
    events:       {},
    min:          0,
    max:          0
  }

  //----- construct component/version/instance tree ----

  // loop over all architecture services
  for (var serviceName in architecture.services) {
    architectureService = architecture.services[serviceName]

    result.components[serviceName] = {
      x:        0,
      y:        0,
      w:        0,
      h:        0,
      versions: {},
      tasks:    {}
    }

    // loop over all service setups
    for (var setupName in architectureService.setups){
      architectureServiceSetup = architectureService.setups[setupName]

      result.components[serviceName].versions[setupName] = {
        x:         0,
        y:         0,
        w:         0,
        h:         0,
        instances: {}
      }

      // add instances
      for (var n=0; n<architectureServiceSetup.size; n++) {
        result.components[serviceName].versions[setupName].instances["Instance-" + n] = {
          x:     0,
          y:     0,
          w:     0,
          h:     0,
          tasks: {}
        }
      }
    }
  }

  // loop over all components
  for (var componentName in domain.components) {
    component = domain.components[componentName]

    // add component if needed
    if (! componentName in result.components) {
      result.components[componentName] = {
        x:        0,
        y:        0,
        w:        0,
        h:        0,
        versions: {},
        tasks:    {}
      }
    }

    // loop over all instances
    for (var instanceName in component.instances) {
      instance    = component.instances[instanceName]
      versionName = instance.version

      // add version if needed
      if (! versionName in result.components.versions) {
        result.components[componentName].versions[versionName] = {
          x:         0,
          y:         0,
          w:         0,
          h:         0,
          instances: {}
        }
      }

      // find last defined but not instantiated instance and remove
      for (var instanceName2 in result.components[componentName].versions[versionName].instances) {
        if (instanceName2.startsWith("Instance")) {
          delete result.components[componentName].versions[versionName].instances[instanceName2]
          break
        }
      }

      // add instance
      result.components[componentName].versions[versionName].instances[instanceName] = {
        x:     0,
        y:     0,
        w:     0,
        h:     0,
        tasks: {}
      }
    }
  }

  // add tasks to components and instances
  for (var taskName in domain.tasks) {
    task = domain.tasks[taskName]

    componentName = task.component
    versionName   = task.version
    instanceName  = task.instance

    if (componentName == "") {
      result.architecture = {
        x: 0,
        y: 0,
        w: 0,
        h: 0
      }
    } else if (instanceName == "") {
      result.components[componentName].tasks[taskName] = {
        x: 0,
        y: 0,
        w: 0,
        h: 0
      }
    } else {
      result.components[componentName].versions[versionName].instances[instanceName].tasks[taskName] = {
        x: 0,
        y: 0,
        w: 0,
        h: 0
      }
    }
  }

  // calculate vertical dimensions of components/versions/instances and tasks
  var lane = 1
  for (var componentName in result.components) {
    component = result.components[componentName]

    component.y = lane++

    for (var taskName in component.tasks) {
      task = component.tasks[taskName]

      task.y = lane++
    }

    for (var versionName in component.versions) {
      version = component.versions[versionName]

      version.y = lane

      for (var instanceName in version.instances) {
        instance = version.instances[instanceName]

        instance.y = lane++

        for (var taskName in instance.tasks) {
          task = instance.tasks[taskName]

          task.y = lane++
        }
      }

      version.h = lane - version.y
    }
  }

  // determine extensions
  for (var eventName in domain.events) {
    event = domain.events[eventName]

    // update min/max
    if (result.min == 0) {
      result.min = event.time
    } else {
      result.min = Math.min(result.min, event.time)
    }
    result.max = Math.max(result.max, event.time)
  }

  // add events
  for (var eventName in domain.events) {
    event = domain.events[eventName]

    // update min/max
    result.min = Math.min(result.min, event.time)
    result.max = Math.max(result.max, event.time)

    // determine tasks
    task1 = getTaskDimension(domain, result.architecture, result.components, event.source)
    task2 = getTaskDimension(domain, result.architecture, result.components, event.task)

    // add event
    result.events[eventName] = {
      x: event.time - result.min,
      y: task1.y,
      w: 0,
      h: task2.y - task1.y,
      t: event.type
    }

    // update source task dimensions
    if (event.source != "") {
      if (task1.x == 0) {
        task1.x = event.time - result.min
        task1.w  = 0
      } else {
        x1 = Math.min(task1.x, event.time - result.min)
        x2 = Math.max(task1.x + task1.w, event.time - result.min)

        task1.x = x1
        task1.w = x2 - x1
      }
    }

    // update task dimensions
    if (task2.x == 0) {
      task2.x = event.time - result.min
      task2.w  = 0
    } else {
      x1 = Math.min(task2.x, event.time - result.min)
      x2 = Math.max(task2.x + task2.w, event.time - result.min)

      task2.x = x1
      task2.w = x2 - x1
    }
  }

  // success
  return result;
}

//------------------------------------------------------------------------------

// getTaskDimension retrieves task
function getTaskDimension(domain, architecture, components, taskName) {
  if (taskName == "") {
    return {
      x: 0,
      y: 0,
      w: 0,
      h: 0
    }
  }

  // determine task dimensions
  task = domain.tasks[taskName]

  componentName = task.component
  versionName   = task.version
  instanceName  = task.instance

  // architecture task
  if (componentName == "") {
    return architecture
  }

  // component task
  if (instanceName == "") {
    return components[componentName].tasks[taskName]
  }

  // instance task
  return components[componentName].versions[versionName].instances[instanceName].tasks[taskName]
}

//------------------------------------------------------------------------------
