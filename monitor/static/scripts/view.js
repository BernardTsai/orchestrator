var view = {
  header:      60,
  title:       60,
  sidebar:    200,
  version:     80,
  indent:       8,
  line:        12,
  task:         4,
  xgap:         4,
  ygap:         4,
  showEvents:   false,
  screen: {
    width:  document.body.clientWidth,
    height: document.body.clientHeight
  },
  tree: {
    tasks:         {
      "": { data: "User", received: {}, triggered: {} }
    },
    architectures: {},
    components:    {}
  },
  domain:        null,
  architectures: {},
  components:    {},
  versions:      {},
  instances:     {},
  tasks:         {
    "": {x: 0, y: 0, w: 0, h: 0, c:0, n: "User", d: null}
  },
  events:        {},
  min:           0,
  max:           0,
  curr:          0,
  range:         0,
  lanes: {
    "user":           {x: 0, y: 0, w: 0, h: 0, n: "user", d: null},
    "architectures":  {x: 0, y: 0, w: 0, h: 0, n: "architectures", d: null},
    "components":     {x: 0, y: 0, w: 0, h: 0, n: "components", d: null}
  }
}
