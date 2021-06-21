//
function createRunningEmulator() {
  removeFrame();
  const emuFrame = createEmuFrame();
  const pgmName = document.querySelector("#program").value;

  const theme = document.querySelector("#theme").value;
  fgcolour = theme.split(",")[0];
  bgcolour = theme.split(",")[1];

  emuFrame.setAttribute("data-pgm-name", "roms/" + pgmName);
  emuFrame.setAttribute("data-speed", document.querySelector("#speed").value);
  emuFrame.setAttribute("data-fgcolour", fgcolour);
  emuFrame.setAttribute("data-bgcolour", bgcolour);

  // By appending the iframe we start everything
  document.querySelector("#wrapper").appendChild(emuFrame);
}

//
function createStoppedEmulator() {
  removeFrame();
  document.querySelector("#wrapper").appendChild(createEmuFrame());
}

//
function removeFrame() {
  oldFrame = document.querySelector("#emulator");
  if (oldFrame) {
    document.querySelector("#wrapper").removeChild(oldFrame);
  }
}

//
function createEmuFrame() {
  var emuFrame = document.createElement("iframe");
  emuFrame.src = "emulator.html";
  emuFrame.id = "emulator";
  emuFrame.width = 1024;
  emuFrame.height = 512;
  return emuFrame;
}
