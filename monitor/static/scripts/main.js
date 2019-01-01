function main() {
  // create the application
  var app = new Vue({
    el:   '#app',
    data: {
      model: model,
      view:  view
    },
    template: `<app v-bind:model="model" v-bind:view="view"></app>`
  })

  // register gesture handlers (e.g. pan, ...)
  startGestures()
}

window.onload = loadModel().then(main);
