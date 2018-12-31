Vue.component(
  'tasks',
  {
    props: ['model', 'view'],
    template: `
    <div id="tasks" v-bind:style="{ top: (view.header + view.title) + 'px', left: view.sidebar + 'px' }">

      <!-- User related tasks -->
      <div id="User">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }"></div>
        <div class="elements">
          <usertask2
            v-for="task in view.tree.tasks"
            v-bind:model="model"
            v-bind:task="task"
            v-bind:view="view">
          </usertask2>
        </div>
      </div>

      <!-- Architecture related tasks -->
      <div id="Architectures">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }"></div>
        <div class="elements">
          <architecture2
            v-for="architecture in view.tree.architectures"
            v-bind:model="model"
            v-bind:architecture="architecture"
            v-bind:view="view">
          </architecture2>
        </div>
      </div>

      <!-- Component related components -->
      <div id="Components">
        <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }"></div>
        <div class="elements">
          <component2
            v-for="component in view.tree.components"
            v-bind:model="model"
            v-bind:component="component"
            v-bind:view="view">
          </component2>
        </div>
      </div>`
  }
)
