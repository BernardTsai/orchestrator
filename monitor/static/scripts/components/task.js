Vue.component(
  'task',
  {
    props: ['model', 'task', 'view'],
    template: `
      <div class="task"
        v-bind:title="task.n"
        v-bind:style="{
          'top':            (task.y * (view.line+view.task)) + 'px',
          'height':         view.line + 'px',
          'left':           (task.x/(view.max-view.min)*1280) + 'px',
          'width':          (task.w/(view.max-view.min)*1280) + 'px',
          'line-height':    view.line + 'px',
          'font-size':      (view.line-4) + 'px'
        }">{{task.n}}</div>`
  }
)
