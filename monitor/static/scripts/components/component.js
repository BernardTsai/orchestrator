Vue.component(
  'component',
  {
    props: ['model', 'component', 'view'],
    template: `
      <div class="component" v-bind:title="component.data.name">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }">{{component.data.name}}</div>
        <div class="elements">
          <componenttask
            v-for="task in component.tasks"
            v-bind:model="model"
            v-bind:task="task"
            v-bind:view="view">
          </componenttask>
        </div>
      </div>`
  }
)
