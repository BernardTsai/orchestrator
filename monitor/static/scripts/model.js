var model;

const loadModel = async () => {
  domain       = "Test"
  architecture = "architecture_1.0.0"
  url          = "http://localhost:8081/data";
  response     = await fetch(url);
  text         = await response.text();
  yaml         = jsyaml.safeLoad(text);
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

  console.log(result)

  // success
  return result;
}

//------------------------------------------------------------------------------

// convertData2 converts yaml into the required model structure for VUE
function convertData2(domain, architecture) {
  // initialize result
  result = {
    lanes:  [],
    tasks:  {},
    events: {},
    min: 0,
    max: 0
  }

  console.log(domain.tasks)

  //----- construct service tree ----
  // loop over all architecture services
  services = {}
  for (var service in architecture.services) {
    services[service] = {}

    architectureService = architecture.services[service]

    // loop over all service setups
    for (var setup in architectureService.setups){
      services[service][setup] = []

      architectureServiceSetup = architectureService.setups[setup]

      // add instances
      for (var n=0; n<architectureServiceSetup.size; n++) {
        services[service][setup].push("Instance-" + n)
      }
    }
  }

  // determine component services
  for (var component in domain.components) {
    // TODO
  }

  //----- construct lanes -----
  serviceKeys = Object.keys(services).sort()
  lane = 0
  for (var serviceKey of serviceKeys) {

    result.lanes.push( {type: "component", component: serviceKey, version: null, instance: null, lane: lane} )
    lane = lane +1

    service = services[serviceKey]

    // loop over versions
    versionKeys = Object.keys(service).sort()
    for (var versionKey of versionKeys) {

      result.lanes.push( {type: "version", component: serviceKey, version: versionKey, instance: null, lane: lane} )
      lane = lane +1

      version = service[versionKey]

      // loop over instances
      for (var instanceKey of version) {

        result.lanes.push( {type: "instance", component: serviceKey, version: versionKey, instance: instanceKey, lane: lane} )
        lane = lane +1
      }
    }
  }

  //----- construct tasks -----
  for (var taskUUID in domain.tasks) {
    task = domain.tasks[taskUUID]
    console.log(task)
    result.tasks[task.UUID] = {
      uuid: task.UUID,
      min:  0,
      max:  0
    }
  }

  // success
  return result;
}
