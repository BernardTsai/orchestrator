Vue.component(
  'timeline',
  {
    props: ['model', 'view'],
    data: {

    },
    template: `
      <div id="timeline" v-bind:style="{ top: view.header + 'px', left: view.sidebar + 'px', height: view.title + 'px' }">
        <div id="min"
          v-bind:style="{
            'height':      view.line + 'px',
            'line-height': view.line + 'px',
            'font-size':   (view.line-4) + 'px'}
        ">{{view.min}}</div>
        <div id="max"
          v-bind:style="{
            'height':      view.line + 'px',
            'line-height': view.line + 'px',
            'font-size':   (view.line-4) + 'px'}
        ">{{view.max}}</div>
        <div id="curr"
          v-if="view.min <= view.curr && view.curr <= view.max"
          v-bind:style="{
            'left':        ((view.curr - view.min) / (view.range) * 100) + '%',
            'height':      view.line + 'px',
            'line-height': view.line + 'px',
            'font-size':   (view.line-4) + 'px'}
        ">{{view.curr}}</div>
      </div>`
  }
)
