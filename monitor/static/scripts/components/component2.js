Vue.component(
  'component2',
  {
    props: ['model', 'component', 'view'],
    template: `
      <div class="component2">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }"></div>
        <div class="elements">
          <componenttask2
            v-for="task in component.tasks"
            v-bind:model="model"
            v-bind:task="task"
            v-bind:view="view">
          </componenttask2>
        </div>
      </div>`
  }
)
