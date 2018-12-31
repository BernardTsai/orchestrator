Vue.component(
  'usertask',
  {
    props: ['model', 'task', 'view'],
    template: `
      <div class="usertask" v-bind:title="task.data.name">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': (view.line - view.ygap) + 'px'  }">{{task.data}}</div>
      </div>`
  }
)
