Vue.component( 'app',
  {
    props:    ['model', 'view'],
    template: `
      <div>
        <div id="header">Node Browser V1.0.0</div>
        <navigation v-bind:model="model" v-bind:view="view"></navigation>
        <detail v-bind:model="model" v-bind:view="view"></detail>
      </div>`
  }
)
