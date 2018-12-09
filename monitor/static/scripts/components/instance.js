Vue.component(
  'instance',
  {
    props: ['instance','view'],
    methods:  {
      toggleDetails: function(event) {
        // toggle node class "details"
        element = event.target;
        while (element != null && !element.classList.contains("instance")) {
          element = element.parentElement
        }

        if (element != null) {
          element.classList.toggle("details");
        }
      }
    },
    template: `
      <div class="instance" v-bind:class="'level' + instance.level">

        <!-- title row -->
        <div class="title">
          <div class="name"><i class="far fa-file"/>&nbsp;{{instance.instance}}
          </div>
          <div class="info">
            <span class="status" v-bind:class="[instance.state]">{{instance.state}}</span>
          </div>
          <div class="buttons"><span v-on:click="toggleDetails($event)"><i class="fas fa-eye"/></div>
        </div>

        <!-- details row -->
        <div class="details">
        ...
        </div>
      </div>`
  }
)
