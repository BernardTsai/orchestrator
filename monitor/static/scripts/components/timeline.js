Vue.component(
  'timeline',
  {
    props: ['model', 'view'],
    computed: {
      timeline: function() {
        prefix = ""
        min    = String(view.min)
        max    = String(view.max)
        pos    = 0
        while ( min.charAt(pos) == max.charAt(pos) ) {
          prefix = prefix + min.charAt(pos)
          pos = pos + 1
        }
        min = Number(min.substr(pos))
        max = Number(max.substr(pos))
        rng = max - min

        pow = Math.floor( Math.log10(rng) - 1 )
        com = Math.pow(10, pow)

        // capture parameters
        view.timeline.min = min
        view.timeline.max = max
        view.timeline.rng = rng
        view.timeline.pow = pow
        view.timeline.com = com

        // small scale
        view.timeline.smallScale.min   = Math.ceil( min / com ) * com
        view.timeline.smallScale.max   = Math.floor( max / com ) * com
        view.timeline.smallScale.step  = com
        view.timeline.smallScale.array = []

        for ( loc =  view.timeline.smallScale.min;
              loc <= max;
              loc += com ) {
          view.timeline.smallScale.array.push( loc )
        }

        // medium scale
        view.timeline.mediumScale.min   = Math.ceil( min / (com*5) ) * (com*5)
        view.timeline.mediumScale.max   = Math.floor( max / (com*5) ) * (com*5)
        view.timeline.mediumScale.step  = com*5
        view.timeline.mediumScale.array = []

        for ( loc =  view.timeline.mediumScale.min;
              loc <= max;
              loc += (com*5) ) {
          view.timeline.mediumScale.array.push( loc )
        }

        // large scale
        view.timeline.largeScale.min   = Math.ceil( min / (com*10) ) * (com*10)
        view.timeline.largeScale.max   = Math.floor( max / (com*10) ) * (com*10)
        view.timeline.largeScale.step  = com*10
        view.timeline.largeScale.array = []

        for ( loc =  view.timeline.largeScale.min;
              loc <= max;
              loc += (com*10) ) {
          view.timeline.largeScale.array.push( loc )
        }

        return view.timeline

      }
    },
    template: `
      <div id="timeline" v-bind:style="{ top: view.header + 'px', left: view.sidebar + 'px', height: view.title + 'px' }">
        <div id="curr"
          v-if="view.min <= view.curr && view.curr <= view.max"
          v-bind:style="{
            'left':        ((view.curr - view.min) / (view.range) * 100) + '%',
            'height':      view.line + 'px',
            'line-height': view.line + 'px',
            'font-size':   (view.line-4) + 'px'}
        ">
          <div class="time">{{view.curr}}</div>
          <div class="pointer"></div>
        </div>

        <div class="large"
          v-for="loc in timeline.largeScale.array"
          v-bind:style="{
            'left':  ((loc - timeline.min) / (timeline.rng) * 100) + '%'}
        "></div>

        <div class="medium"
          v-for="loc in timeline.mediumScale.array"
          v-bind:style="{
            'left':  ((loc - timeline.min) / (timeline.rng) * 100) + '%'}
        "></div>

        <div class="small"
          v-for="loc in timeline.smallScale.array"
          v-bind:style="{
            'left':  ((loc - timeline.min) / (timeline.rng) * 100) + '%'}
        "></div>

      </div>`
  }
)
