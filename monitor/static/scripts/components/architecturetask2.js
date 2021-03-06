Vue.component(
  'architecturetask2',
  {
    props: ['model', 'task', 'view'],
    template: `
      <div class="architecturetask2"
        v-bind:id="'task-' + task.data.uuid"
        v-bind:title="task.data.type + ':' + task.data.uuid">
        <div class="background"
          v-bind:style="{
            'height':           (view.line) + 'px',
            'line-height':      (view.line) + 'px',
            'font-size':        (view.line - view.ygap) + 'px',
            'margin-bottom':    (view.ygap) + 'px',
            'margin-left':      (((task.start-view.min)/view.range)*100) + '%',
            'border-radius':    (view.line/2) + 'px',
            'width':            (((task.end-task.start)/view.range)*100) + '%'
          }"
        ></div>
        <div class="title"
          v-bind:style="{
            'height':           (view.line) + 'px',
            'line-height':      (view.line) + 'px',
            'font-size':        (view.line - view.ygap) + 'px',
            'margin-top':       (-view.line-view.ygap) + 'px',
            'margin-bottom':    (view.ygap) + 'px',
            'margin-left':      (((task.first-view.min)/view.range)*100) + '%',
            'border-radius':    (view.line/2) + 'px',
            'width':            (((task.last-task.first)/view.range)*100) + '%',
            'color':            task.data.status == 0 ? 'black'     : 'white',
            'background-color': task.data.status == 0 ? 'lightgrey' : 'blue'
          }"
        >{{task.data.type + " (" + task.data.status + ")"}}</div>
      </div>`
  }
)
