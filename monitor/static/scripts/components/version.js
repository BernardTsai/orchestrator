Vue.component(
  'version',
  {
    props: ['model', 'version', 'view'],
    template: `
      <div class="version" v-bind:style="{
        'padding-left':   view.indent + 'px',
        'width':          (2*view.indent) + 'px',
        'line-height':    view.line + 'px',
        'font-size':      (view.line-4) + 'px',
        'left':           view.indent + 'px',
        'top':            (version.y * (view.line+view.task)) + 'px',
        'height':         (version.h * (view.line+view.task)) + 'px'
      }">{{version.n}}</div>`
  }
)
