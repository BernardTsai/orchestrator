rearrange view Components


<!-- architecture lane -->
<architecture v-bind:model="model" v-bind:view="view"></architecture>

<!-- component lanes -->
<comp
  v-for="component in view.components"
  v-bind:model="model"
  v-bind:component="component"
  v-bind:view="view">
  {{component}}
</comp>

<!-- version lanes -->
<version
  v-for="version in view.versions"
  v-bind:model="model"
  v-bind:version="version"
  v-bind:view="view"></version>

<!-- instance lanes -->
<instance
  v-for="instance in view.instances"
  v-bind:model="model"
  v-bind:instance="instance"
  v-bind:view="view"></instance>
</div>
