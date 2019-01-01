function startGestures() {
  var tasksElement = document.getElementById('tasks');

  // create a simple instance
  // by default, it only adds horizontal recognizers
  var mc = new Hammer(tasksElement);

  mc.get('pan').set({ direction: Hammer.DIRECTION_ALL });

  // listen to events...
  // tie in the handler that will be called
  mc.on("pan", panTasks);
}


// poor choice here, but to keep it simple
// setting up a few vars to keep track of things.
// at issue is these values need to be encapsulated
// in some scope other than global.
var lastPosX = 0;
var lastPosY = 0;
var isDragging = false;

function panTasks(ev) {

  // for convience, let's get a reference to our object
  var elem = document.getElementById('tasks');

  // DRAG STARTED
  // here, let's snag the current position
  // and keep track of the fact that we're dragging
  if ( ! isDragging ) {
    isDragging = true;
    lastPosX = elem.offsetLeft;
    lastPosY = elem.offsetTop;
  }

  // we simply need to determine where the x,y of this
  // object is relative to where it's "last" known position is
  // NOTE:
  //    deltaX and deltaY are cumulative
  // Thus we need to always calculate 'real x and y' relative
  // to the "lastPosX/Y"
  var posX = Math.min(ev.deltaX + lastPosX, view.sidebar);
  var posY = Math.min(ev.deltaY + lastPosY, view.header + view.title);

  posX = Math.max(posX, view.sidebar             + (view.screen.width  - view.sidebar)               - elem.clientWidth)
  posY = Math.max(posY, view.header + view.title + (view.screen.height - (view.header + view.title)) - elem.clientHeight)

  // move our element to that position
  elem.style.left = posX + "px";
  elem.style.top = posY + "px";

  var elem2 = document.getElementById('sidebar');
  elem2.style.top = posY + "px";


  // DRAG ENDED
  // this is where we simply forget we are dragging
  if (ev.isFinal) {
    isDragging = false;
  }
}
