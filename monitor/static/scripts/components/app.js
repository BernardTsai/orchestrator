Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <div id="header"   v-bind:style="{ height: view.header + 'px' }">Monitor V1.0.0</div>
        <div id="context"  v-bind:style="{ top: view.header + 'px', width: view.sidebar + 'px', height: view.title + 'px' }"></div>
        <div id="timeline" v-bind:style="{ top: view.header + 'px', left: view.sidebar + 'px', height: view.title + 'px' }"></div>
        <div id="sidebar"  v-bind:style="{ top: (view.header + view.title) + 'px', width: view.sidebar + 'px' }"></div>
        <div id="events"   v-bind:style="{ top: (view.header + view.title) + 'px', left: view.sidebar + 'px' }"></div>
        <!-- navigation v-bind:model="model" v-bind:view="view"></navigation>
        <detail v-bind:model="model" v-bind:view="view"></detail -->
      </div>`
  }
)
