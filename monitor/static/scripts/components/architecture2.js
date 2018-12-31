Vue.component(
  'architecture2',
  {
    props: ['model', 'architecture', 'view'],
    template: `
      <div class="architecture">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }"></div>
        <div class="elements">
          <architecturetask2
            v-for="task in architecture.tasks"
            v-bind:model="model"
            v-bind:task="task"
            v-bind:view="view">
          </architecturetask2>
        </div>
      </div>`
  }
)
