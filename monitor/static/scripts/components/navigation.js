Vue.component(
  'navigation',
  {
    props: ['model', 'view'],
    template: `
      <div id="navigation">

        <!-- domain node -->
        <domain v-bind:model="model" v-bind:view="view"></domain>

      </div>`
  }
)
