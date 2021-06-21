// Size for CHIP8 not Super-CHIP8
const CHIP8_DISPLAY_WIDTH = 64;
const VERSION = "1.0.0";

// NOTE. We get our params from the iframe data attributes, hey it's not elegant but it works
// IMO it's better than postMessage

//
// Start here on document body loading
//
function bodyLoaded() {
  // We need this so we can use the Fn keys in the emulator
  const body = document.querySelector("body");
  body.onkeydown = function (e) {
    if (!e.metaKey) {
      e.preventDefault();
    }
  };

  // We use pgm name to determine if we boot & run emulator or just show a banner
  const pgmName = window.frameElement.getAttribute("data-pgm-name");
  if (pgmName != null) {
    runEmulator(pgmName);
  } else {
    var pre = document.createElement("pre");
    pre.innerText = `
       ██████╗  ██████╗      ██████╗██╗  ██╗██╗██████╗        █████╗
      ██╔════╝ ██╔═══██╗    ██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
      ██║  ███╗██║   ██║    ██║     ███████║██║██████╔╝█████╗╚█████╔╝
      ██║   ██║██║   ██║    ██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
      ╚██████╔╝╚██████╔╝    ╚██████╗██║  ██║██║██║           ╚█████╔╝
       ╚═════╝  ╚═════╝      ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝

    +++ Version ${VERSION}
    +++ Enabling WASM data & protocol sync buffers
    +++ Mem check ~ 4096 bytes: OK
    +++ System cold boot completed with 0 warnings
      > █`;
    pre.setAttribute("unselectable", "on");
    pre.setAttribute("onselectstart", "return false");
    document.body.append(pre);
  }
}

//
// Run the CHIP emulator for given program name
//
async function runEmulator(pgmName) {
  let errored = false;

  // Calc pixelsize to request based on our iframe size
  const pixelSize = Math.floor(window.frameElement.width / CHIP8_DISPLAY_WIDTH);
  const speed = window.frameElement.getAttribute("data-speed");
  const fgColour = window.frameElement.getAttribute("data-fgcolour");
  const bgColour = window.frameElement.getAttribute("data-bgcolour");

  // Load CHIP-8 WASM
  const go = new Go();

  // Kinda janky but through main is our only way of working with ebiten wasm
  go.argv = [pgmName, false, speed, pixelSize, fgColour, bgColour];
  const wasm = await WebAssembly.instantiateStreaming(
    fetch("chip8.wasm"),
    go.importObject
  );

  // We override the wasm_exec.js exit function to trap errors
  // Doesn't seem to be a way to get error message so we use the exit codes
  go.exit = (code) => {
    if (code !== 0) {
      var pre = document.createElement("pre");
      reason = "General error";
      // Codes above 100 are HTTP errors from cmd/chip8wasm/main.go
      if (code > 100) {
        reason = "Failed to download program: " + pgmName;
      }
      // Other codes come from system errors trapped by pkg/emulator/emulator.go
      if (code == 51) {
        reason = "Memory address out of bounds";
      }
      if (code == 52) {
        reason = "Invalid Opcode";
      }
      if (code == 53) {
        reason = "Out of memory";
      }
      pre.innerText = `+++ Guru Meditation!\n+++ CODE: ${code}\n+++ REASON: ${reason}`;
      document.body.append(pre);
      var canvas = document.querySelector("canvas");
      document.body.removeChild(canvas);
    }
  };

  // Only way to trap anything, try catch does NOT work
  window.addEventListener("error", function (event) {
    if (!errored) {
      errored = true;
    }
  });

  // FINALLY we actually start the WASM main function running, phew!
  go.run(wasm.instance);
}
