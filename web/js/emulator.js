const CHIP8_DISPLAY_WIDTH = 128
const VERSION = '1.0.0'

// NOTE. We get our params from the iframe data attributes, hey it's not elegant but it works
// IMO it's better than postMessage

//
// Start here on document body loading
//
async function main() {
  // We need this so we can use the Fn keys in the emulator
  document.body.onkeydown = function (e) {
    if (!e.metaKey) {
      e.preventDefault()
    }
  }

  // Start the emulator on a message from parent frame
  window.addEventListener(
    'message',
    async (msg) => {
      let errored = false

      // Remove old WASM canvas if we have one
      const canvas = document.querySelector('canvas')
      if (canvas) {
        document.body.removeChild(canvas)
      }

      const go = new Go()

      // We override the wasm_exec.js exit function to trap errors
      // Doesn't seem to be a way to get error message so we use the exit codes
      go.exit = (code) => {
        if (code !== 0) {
          var pre = document.createElement('pre')
          reason = 'General error'
          // Codes above 100 are HTTP errors from cmd/chip8wasm/main.go
          if (code > 100) {
            reason = 'Failed to download program: ' + msg.data.programName
          }
          // Other codes come from system errors trapped by pkg/emulator/emulator.go
          if (code == 51) {
            reason = 'Memory address out of bounds'
          }
          if (code == 52) {
            reason = 'Invalid Opcode'
          }
          if (code == 53) {
            reason = 'Out of memory'
          }
          pre.innerText = `+++ Guru Meditation!\n+++ CODE: ${code}\n+++ REASON: ${reason}`
          document.body.append(pre)
          var canvas = document.querySelector('canvas')
          document.body.removeChild(canvas)
        }
      }

      // Only way to trap anything, try catch does NOT work
      window.addEventListener('error', function (event) {
        if (!errored) {
          errored = true
        }
      })

      // Remove boot screen if it's there
      const bootScreen = document.querySelector('#bootScreen')
      if (bootScreen) {
        document.body.removeChild(bootScreen)
      }

      const wasm = await WebAssembly.instantiateStreaming(fetch('chip8.wasm'), go.importObject)
      const pixelSize = Math.floor(window.frameElement.width / CHIP8_DISPLAY_WIDTH)

      go.argv = ['roms/' + msg.data.programName, false, msg.data.speed, pixelSize, msg.data.fgColour, msg.data.bgColour]
      go.run(wasm.instance)
    },
    false
  )
}
