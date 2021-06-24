//
// Tell the emulator to start running with a postMessage
//
function run() {
  // This removes focus e.g. from the select rom dropdown
  document.activeElement.blur()

  const theme = document.querySelector('#theme').value
  const emuFrame = document.querySelector('#emulator')

  emuFrame.contentWindow.postMessage(
    {
      programName: document.querySelector('#program').value,
      speed: document.querySelector('#speed').value,
      fgColour: theme.split(',')[0],
      bgColour: theme.split(',')[1],
    },
    '*'
  )
}

//
// Simply hard reset the emulator by reloading the frame
//
function reset() {
  const emuFrame = document.querySelector('#emulator')
  emuFrame.contentWindow.location.reload()
}
