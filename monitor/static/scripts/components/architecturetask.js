Vue.component(
  'architecturetask',
  {
    props: ['model', 'task', 'view'],
    template: `
      <div class="architecturetask" v-bind:title="task.data.uuid">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': (view.line - view.ygap) + 'px'  }">{{task.data.uuid}}</div>
      </div>`
  }
)
