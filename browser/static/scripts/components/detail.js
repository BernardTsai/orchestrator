Vue.component( 'detail',
  {
    props:    ['model','view'],
    created: function() {
       window.addEventListener('keyup', this.handleKeys)
    },
    methods: {
      close: function() {
        this.view.focus = "";
      },
      handleKeys: function(keyup) {
        if (keyup.key === "Escape"){
          this.close();
        }
      }
    },
    template: `
      <div id="detail" v-if="view.focus !== ''">
        <div class="content">
          <div class="close" v-on:click="close()"><i class="fas fa-times-circle"/></div>

          <!-- node detail -->
          <nodeDetail
            v-if="view.focus === 'node'"
            v-bind:node="view.node"
            v-bind:view="view">
          </nodeDetail>

          <!-- version detail -->
          <versionDetail
            v-if="view.focus === 'version'"
            v-bind:version="view.version"
            v-bind:view="view">
          </versionDetail>
        </div>
      </div>`
  }
)
