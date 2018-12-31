Vue.component(
  'sidebar',
  {
    props: ['model', 'view'],
    template: `
      <div id="sidebar" v-bind:style="{ top: (view.header + view.title) + 'px', width: view.sidebar + 'px' }">

        <!-- User related components -->
        <div id="User">
          <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }">Users</div>
          <div class="elements">
            <usertask
              v-for="task in view.tree.tasks"
              v-bind:model="model"
              v-bind:task="task"
              v-bind:view="view">
            </usertask>
          </div>
        </div>

        <!-- Architecture related components -->
        <div id="Architectures">
          <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }">Architectures</div>
          <div class="elements">
            <architecture
              v-for="architecture in view.tree.architectures"
              v-bind:model="model"
              v-bind:architecture="architecture"
              v-bind:view="view">
            </architecture>
          </div>
        </div>

        <!-- Component related components -->
        <div id="Components">
          <div class="title" v-bind:style="{ height: (view.line + view.ygap) + 'px', 'line-height': (view.line + view.ygap) + 'px', 'font-size': view.line + 'px'  }">Components</div>
          <div class="elements">
            <component
              v-for="component in view.tree.components"
              v-bind:model="model"
              v-bind:component="component"
              v-bind:view="view">
            </component>
          </div>
        </div>`
  }
)
