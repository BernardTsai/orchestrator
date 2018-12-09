function main() {
  var app = new Vue({
    el:   '#app',
    data: {
      model: model
    },
    template: `<app v-bind:model="model"></app>`
  })
}

window.onload = loadModel().then(main);
