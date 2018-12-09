Vue.component(
  'node',
  {
    props: ['node', 'view'],
    methods:  {
          toggleExpand: function(event) {
            // toggle node class "expanded"
            element = event.target;
            parent  = element.parentElement.parentElement;
            parent.classList.toggle("expanded");
          },
          showDetails: function(event) {
            view.focus = "node";
            view.node  = this.node;
          }
        },
    template: `
      <!-- standard node -->
      <div class="node" v-bind:class="'level' + node.level">

        <!-- title row -->
        <div class="title">
          <div class="name" v-on:click="toggleExpand($event)">
            <i class="fas fa-folder default"/><i class="fas fa-folder-open expanded"/>&nbsp;{{node.component}}
          </div>
          <div class="info">
            <span v-for="version, versionNumber in node.versions" class="version" v-bind:class="[version.state]">V{{version.version}}&nbsp;</span>
          </div>
          <div class="buttons"><span v-on:click="showDetails($event)"><i class="fas fa-eye"/></span></div>
        </div>

        <!-- children row -->
        <div class="children">
          <node
            v-for="child in node.children"
            v-bind:node="child"
            v-bind:view="view"></node>
        </div>
      </div>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'nodeDetail',
  {
    props: ['node', 'view'],
    data: function() {
      for (var key in this.node.versions) {
        return {versionID: key, content: "undefined"}
      }

      return {versionID: ""}
    },
    computed: {
      versions: function() {
        return Object.keys(this.node.versions)
      },
      version: function() {
        return this.node.versions[this.versionID]
      }
    },
    mounted: function(){
      for (var key in this.node.versions) {
        this.view.version = this.node.versions[key];
        break;
      }
    },
    methods:  {
          selectTab: function(event) {
            var tab  = event.target;
            var name = event.target.innerHTML;

            // update tab display
            var tabs = document.getElementById("nodeDetail").getElementsByClassName("tab");
            for (var t of tabs) {
              t.classList.remove("selected")
            }
            tab.classList.add("selected");

            // activate appropriate pabe
            var panes = document.getElementById("nodeDetail").getElementsByClassName("pane");
            for (var p of panes) {
              p.classList.remove("selected")
            }

            var pane = document.getElementById("nodeDetail").getElementsByClassName(name + "Tab");
            pane[0].classList.add("selected");

            // update content
            node = this.node.path.replace("/tmp/data/", "")
            url = '/render/' + node + ":" + this.version.version;
            loadData(url).then(content => {this.content=content})
          },
          configuration: function() {
            result = "";

            for (var key in this.version.instances) {
              instance = this.version.instances[key];
              result = instance.configuration;
              break;
            }

            return jsyaml.safeDump(result);
          },
          dependencies: function() {
            result = "";

            for (var key in this.version.instances) {
              instance = this.version.instances[key];
              result = instance.dependencies;
              break;
            }

            return result;
          }
        },
    template: `
      <!-- node detail -->
      <div id="nodeDetail" class="nodeDetail">
        <div class="header">
          <div class="nameLabel">Node:&nbsp;</div>
          <div class="name">{{node.component}}</div>
          <div class="seperator"/>
          <div class="versionLabel">Version:&nbsp;</div>
          <select class="versionSelector" v-model="versionID">
            <option v-for="versionNumber in versions">{{versionNumber}}</option>
          </select>
          <div class="seperator"/>
        </div>

        <!-- tabs -->
        <div class="tabs">
          <div v-on:click="selectTab($event)" class="tab general selected">General</div>
          <div v-on:click="selectTab($event)" class="tab instance">Instances</div>
          <div v-on:click="selectTab($event)" class="tab service">Service</div>
          <div v-on:click="selectTab($event)" class="tab content">Content</div>
          <div v-on:click="selectTab($event)" class="tab configuration">Configuration</div>
          <div v-on:click="selectTab($event)" class="tab dependencies">Dependencies</div>
          <div v-on:click="selectTab($event)" class="tab context">Context</div>
        </div>

        <!-- General Tab -->
        <div class="pane GeneralTab selected" v-if="version">
          <table>
            <tr>
              <th>Domain:</th>
              <td>{{version.domain}}</td>
            </tr>
            <tr>
              <th>Component:</th>
              <td>{{version.component}}</td>
            </tr>
            <tr>
              <th>Version:</th>
              <td>{{version.version}}</td>
            </tr>
            <tr>
              <th>State:</th>
              <td>{{version.state}}</td>
            </tr>
          </table>
        </div>

        <!-- Service Tab -->
        <div class="pane ServiceTab" v-if="version">
          <table>
            <tr>
              <th>Path:</th>
              <td>{{version.path}}</td>
            </tr>
          </table>
        </div>

        <!-- Instances Tab -->
        <div class="pane InstancesTab" v-if="version">
          <table>
            <tr v-for="(instance,key,index) in version.instances">
              <th><span v-if="index === 0">Instances:</span></th>
              <td>{{instance.instance}} ({{instance.state}})</td>
            </tr>
          </table>
        </div>

        <!-- Content Tab -->
        <div class="pane ContentTab" v-if="version">
          <h3>Content</h3>
          <pre>{{content}}</pre>
        </div>

        <!-- Configuration Tab -->
        <div class="pane ConfigurationTab" v-if="version">
          <h3>Configuration</h3>
          <pre>{{configuration()}}</pre>
        </div>

        <!-- Dependencies Tab -->
        <div class="pane DependenciesTab" v-if="version">
          <table>
            <tr>
              <th>Reference</th>
              <th>Component</th>
              <th>Version</th>
              <th>State</th>
              <th class="expand">Endpoint</th>
            </tr>
            <tr v-for="(dependency,key,index) in dependencies()" v-if="dependency.type==='service'">
              <td>{{key}}</td>
              <td>{{dependency.component}}</td>
              <td>{{dependency.version}}</td>
              <td>{{dependency.state}}</td>
              <td>{{dependency.endpoint}}</td>
            </tr>
          </table>
        </div>

        <!-- Context Tab -->
        <div class="pane ContextTab" v-if="version">
          <table>
            <tr>
              <th>Reference</th>
              <th>Component</th>
              <th>Version</th>
              <th>State</th>
              <th class="expand">Endpoint</th>
            </tr>
            <tr v-for="(dependency,key,index) in dependencies()" v-if="dependency.type==='context'">
              <td>{{key}}</td>
              <td>{{dependency.component}}</td>
              <td>{{dependency.version}}</td>
              <td>{{dependency.state}}</td>
              <td>{{dependency.endpoint}}</td>
            </tr>
          </table>
        </div>

      </div>`
  }
)
