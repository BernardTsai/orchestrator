var model;

const loadModel = async () => {
  url      = "http://localhost:8080/data/";
  response = await fetch(url);
  text     = await response.text();
  yaml     = jsyaml.safeLoad(text);
  model    = convertModel(yaml);
}

// convertModel converts yaml into the required model structure for VUE
function convertModel(yaml) {
  // create top level domain node
  var result = {
    name: yaml.name,
    expanded: true,
    details: false,
    children: {}
  };

  // add children to top level domain node
  for (var node of yaml.children) {
    key = node.name;
    result.children[key] = convertNode(node, true, 1);
  }

  // finished
  return result;
}

// convertNode converts a single node into the required model structure for VUE
function convertNode(node, expanded, level) {
  // create node
  var result = {
    domain:    node.component.domain,
    component: node.component.component,
    path:      node.component.path,
    versions:  {},
    children:  {},
    expanded:  expanded,
    details:   false,
    level:     level
  };

  // construct versions
  var versions = {};

  for (const [key, instance] of Object.entries(node.instances)) {
    version = instance.version;

    // version if needed
    if ( !(version in versions) ) {
      versions[version] = {};
    }

    // add instance to version
    versions[version][key] = instance;
  }

  // add versions to node
  for (const [key, version] of Object.entries(versions)) {
    result.versions[key] = convertVersion(node, key, version, level+1)
  }

  // add children to node
  for (const [key, childNode] of Object.entries(node.children)) {
    result.children[key] = convertNode(childNode, false, level+1)
  }

  // finished
  return result;
}

// convertVersion converts a version into the required model structure for VUE
function convertVersion(node, key, version, level) {
  // create version
  var result = {
    domain:    node.component.domain,
    component: node.component.component,
    version:   key,
    state:     "unknown",
    path:      node.component.path,
    instances: {},
    expanded:  false,
    details:   false,
    level:     level
  }

  // add instances to version
  for (const [key, instance] of Object.entries(version)) {
    result.instances[key] = convertInstance(instance, level+1)
  }

  // determine state of version
  for (const [key, instance] of Object.entries(result.instances)) {
    // do we know anythin about the instance
    if (result.state == "unknown") {
      result.state = instance.state;
    }
    // check for inactive or active state
    if ((result.state == "initial" || result.state == "inactive") &&
        (instance.state == "inactive" || instance.state == "active")){
      result.state = instance.state;
    }
    // check for failure
    if (instance.state == "failure"){
      result.state = instance.state;
    }
  }

  // finished
  return result;
}

// convertInstance converts an instance into the required model structure for VUE
function convertInstance(instance, level) {
  // create instance
  var result = {
    domain:        instance.domain,
    component:     instance.component,
    version:       instance.version,
    instance:      instance.instance,
    state:         instance.state,
    path:          instance.path,
    endpoint:      instance.endpoint,
    configuration: instance.configuration,
    dependencies:  instance.dependencies,
    expanded:      false,
    details:       false,
    level:         level
  };

  // finished
  return result;
}
