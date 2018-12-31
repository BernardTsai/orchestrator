Vue.component(
  'instance',
  {
    props: ['model', 'instance', 'view'],
    template: `
      <div class="instance" v-bind:style="{
        'padding-left':   view.indent + 'px',
        'width':          (view.sidebar-3*view.indent) + 'px',
        'height':         (view.line+view.task) + 'px',
        'line-height':    (view.line+view.task) + 'px',
        'font-size':      (view.line-4) + 'px',
        'top':            (instance.y * (view.line+view.task)) + 'px',
        'left':           (3*view.indent) + 'px'
      }">{{instance.n}}</div>`
  }
)
