Vue.component(
  'version',
  {
    props: ['version', 'view'],
    methods:  {
          toggleExpand: function(event) {
            // toggle version class "expanded"
            element = event.target;
            parent  = element.parentElement.parentElement;
            parent.classList.toggle("expanded");
          },
          showDetails: function() {
            view.focus   = "version";
            view.version = this.version;
          }
        },
    template: `
      <!-- version node -->
      <div class="version" v-bind:class="'level' + version.level">

        <!-- title row -->
        <div class="title">
          <div class="name" v-on:click="toggleExpand($event)" v-on:click.shift="showDetails()">
            <span class="fa-layers fa-fw default">
              <i class="fas fa-hashtag"></i>
              <span class="fa-layers-counter" style="background:Tomato">{{Object.keys(version.instances).length}}</span>
            </span><span class="expanded">
              <i class="fas fa-hashtag"/>
            </span>&nbsp;V{{version.version}}
          </div>
          <div class="info">
            <span v-for="instance, instanceID in version.instances" class="instance" v-bind:class="[instance.state]">{{instance.instance}}</span>
          </div>
          <div class="buttons"><span v-on:click="showDetails()"><i class="fas fa-eye"/></div>
        </div>

        <!-- instances row -->
        <div class="instances">
          <instance
            v-for="instance, instanceID in version.instances"
            v-bind:instance="instance"
            v-bind:view="view"></instance>
        </div>
      </div>`
  }
)

//------------------------------------------------------------------------------

Vue.component(
  'versionDetail',
  {
    props: ['version', 'view'],
    methods:  {
          selectTab: function(event) {
            var tab  = event.target;
            var name = event.target.innerHTML;

            // update tab display
            var tabs = document.getElementById("versionDetail").getElementsByClassName("tab");
            for (var t of tabs) {
              t.classList.remove("selected")
            }
            tab.classList.add("selected");

            // activate appropriate pabe
            var panes = document.getElementById("versionDetail").getElementsByClassName("pane");
            for (var p of panes) {
              p.classList.remove("selected")
            }

            var pane = document.getElementById("versionDetail").getElementsByClassName(name + "Tab");
            pane[0].classList.add("selected");
          },
          content: function() {
            // response = fetch('/render/' + this.node.path + ":" + this.node.Version);
            return this.version;
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
      <!-- version detail -->
      <div id="versionDetail" class="versionDetail">
        <h2>Version</h2>

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
        <div class="pane GeneralTab selected">
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
        <div class="pane ServiceTab">
          <table>
            <tr>
              <th>Path:</th>
              <td>{{version.path}}</td>
            </tr>
          </table>
        </div>

        <!-- Instances Tab -->
        <div class="pane InstancesTab">
          <table>
            <tr v-for="(instance,key,index) in version.instances">
              <th><span v-if="index === 0">Instances:</span></th>
              <td>{{instance.instance}} ({{instance.state}})</td>
            </tr>
          </table>
        </div>

        <!-- Content Tab -->
        <div class="pane ContentTab">
          <h3>Content</h3>
        </div>

        <!-- Configuration Tab -->
        <div class="pane ConfigurationTab">
          <h3>Configuration</h3>
          <pre>{{configuration()}}</pre>
        </div>

        <!-- Dependencies Tab -->
        <div class="pane DependenciesTab">
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
        <div class="pane ContextTab">
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
