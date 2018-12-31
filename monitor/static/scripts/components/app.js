Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <tasks v-bind:model="model" v-bind:view="view"></tasks>
        <sidebar v-bind:model="model" v-bind:view="view"></sidebar>
        <div id="header"   v-bind:style="{ height: view.header + 'px' }">Monitor V1.0.0</div>
        <div id="context"  v-bind:style="{ top: view.header + 'px', width: view.sidebar + 'px', height: view.title + 'px' }"></div>
        <div id="timeline" v-bind:style="{ top: view.header + 'px', left: view.sidebar + 'px', height: view.title + 'px' }"></div>
      </div>`
  }
)
