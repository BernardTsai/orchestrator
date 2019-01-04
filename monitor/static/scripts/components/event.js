Vue.component(
  'event',
  {
    props: ['model', 'event', 'view'],
    methods: {
      x: function(event) {
        return (((event.time-view.min)/view.range)*100) + '%'
      },
      y: function(event) {
        src = document.getElementById("task-" + event.source)
        dst = document.getElementById("task-" + event.task)

        y1 = src.offsetTop
        y2 = dst.offsetTop

        return (Math.min(y1,y2) + view.line/2) +'px'
      },
      h: function(event) {
        src = document.getElementById("task-" + event.source)
        dst = document.getElementById("task-" + event.task)

        y1 = src.offsetTop
        y2 = dst.offsetTop

        return (Math.abs(y2-y1)) + 'px'
      },
      orientation: function(event) {
        src = document.getElementById("task-" + event.source)
        dst = document.getElementById("task-" + event.task)

        y1 = src.offsetTop
        y2 = dst.offsetTop

        return y1 == y2 ? 'same' :
               y2 <  y1 ? 'up' : 'down'
      }
    },
    template: `
      <div class="event"
        v-bind:id="event.uuid"
        v-bind:class="{
          executionEvent:   event.type         == 'execution',
          completionEvent:  event.type         == 'completion',
          terminationEvent: event.type         == 'termination',
          failureEvent:     event.type         == 'failure',
          timeoutEvent:     event.type         == 'timeout',
          unknownEvent:     event.type         == 'unknown',
          orientationUp:    orientation(event) == 'up',
          orientationDown:  orientation(event) == 'down',
          orientationSame:  orientation(event) == 'same'
        }"
        v-bind:style="{ 'margin-left':x(event), top:y(event), height: h(event) }"
        v-bind:title="event.type">
        <div class="shaft"></div>
        <div class="base"></div>
        <div class="arrow"></div>
      </div>`
  }
)
