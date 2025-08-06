// Touch to mouse event converter
function touchHandler(event) {
  let touch = event.changedTouches[0];
  let simulatedEvent;

  switch (event.type) {
    case "touchstart":
      simulatedEvent = new MouseEvent("mousedown", {
        bubbles: true,
        clientX: touch.clientX,
        clientY: touch.clientY
      });
      break;
    case "touchend":
      simulatedEvent = new MouseEvent("mouseup", {
        bubbles: true,
        clientX: touch.clientX,
        clientY: touch.clientY
      });
      break;
    case "touchmove":
      simulatedEvent = new MouseEvent("mousemove", {
        bubbles: true,
        clientX: touch.clientX,
        clientY: touch.clientY
      });
      break;
    default:
      return;
  }

  touch.target.dispatchEvent(simulatedEvent);
  event.preventDefault();
}

function initTouchToMouse(canvas) {
  canvas.addEventListener("touchstart", touchHandler, true);
  canvas.addEventListener("touchend", touchHandler, true);
  canvas.addEventListener("touchmove", touchHandler, true);
}

function activeTouchesToMouseInDevice(){
  // Activar touch -> mouse only in devices with touch screen.
  const canvas = document.querySelector("canvas");
  if ('ontouchstart' in window && canvas) {
    initTouchToMouse(canvas);
  } 
}

// Polyfill for browsers that don't support instantiateStreaming
function polyfillIinstantiateStreamingSupport(){
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }
}
