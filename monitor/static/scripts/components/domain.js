Vue.component(
  'domain',
  {
    props: ['model', 'view'],
    methods:  {
      toggleDetails: function(event) {
        // toggle node class "details"
        element = event.target;
        while (element != null && !element.classList.contains("domain")) {
          element = element.parentElement
        }

        if (element != null) {
          element.classList.toggle("details");
        }
      }
    },
    template: `
      <!-- domain node -->
      <div class="domain" v-bind:class="{expanded: model.expanded, details: model.details}">

        <!-- title row -->
        <div class="title">
          <div class="name">
            <i class="fas fa-sitemap"></i>&nbsp;Domain
          </div>
          <div class="info" v-on:click="toggleDetails($event)">&nbsp;</div>
          <div class="buttons"><span v-on:click="toggleDetails($event)"><i class="fas fa-eye"/></div>
        </div>

        <!-- details row -->
        <div class="details">
        ...
        </div>

        <!-- children row -->
        <div class="children">

          <!-- top level nodes -->
          <node
            v-for="child in model.children"
            v-bind:node="child"
            v-bind:view="view"></node>

        </div>
        <!-- end of children row -->

      </div>
      <!-- end of domain -->`
  }
)
