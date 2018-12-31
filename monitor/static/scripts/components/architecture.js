Vue.component(
  'architecture',
  {
    props: ['model', 'architecture', 'view'],
    template: `
      <div class="architecture" v-bind:title="architecture.data.name">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }">{{architecture.data.name}}</div>
        <div class="elements">
          <architecturetask
            v-for="task in architecture.tasks"
            v-bind:model="model"
            v-bind:task="task"
            v-bind:view="view">
          </architecturetask>
        </div>
      </div>`
  }
)
