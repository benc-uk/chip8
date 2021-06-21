# Go CHIP-8

This is (yet another) CHIP-8 emulator written in Go. Yes there are many others but this one is mine.

What is CHIP-8? It's a simple 8-bit virtual machine developed in the 1970s, which found new life in the 90s on HP graphing calculators. It's become the defacto "hello world" for emulator developers to get started with.

https://en.wikipedia.org/wiki/CHIP-8

### Screenshots
<img src="./docs/screens/astro.jpg" width="500">
<img src="./docs/screens/rps.jpg" width="500">
<img src="./docs/screens/blinky.jpg" width="500">

Goals:

- Write my first emulator
- Learn about virtual machine development (no, I mean _real_ [virtual machines](https://wiki.c2.com/?VirtualMachine), not that ugly IaaS stuff)
- Have fun

Use cases & key features:

- WASM support, playable in browser as well as Linux and Windows binaries
- Configurable speed
- Builtin debugger

Supporting technologies and libraries:

- [Ebiten](https://github.com/hajimehoshi/ebiten) 

## Status

![](https://img.shields.io/github/license/benc-uk/chip8)
![](https://img.shields.io/github/last-commit/benc-uk/chip8)
![](https://img.shields.io/github/release/benc-uk/chip8)
![](https://img.shields.io/github/checks-status/benc-uk/chip8/main)
![](https://img.shields.io/github/workflow/status/benc-uk/chip8/CI%20Build?label=ci-build)
![](https://img.shields.io/github/workflow/status/benc-uk/chip8/Release%20Binaries?label=release)

# Getting Started

## Try web version

Online WASM version is playable here with a small set of roms/programs provided

https://code.benco.io/chip8/web/

## Installing & Running

Download from [releases](https://github.com/benc-uk/chip8/releases), unzip/untar and run :)

```
chip8 <flags> <program-filename>

  -debug
        Enable debug, lots of output very slow
  -scale int
        Size of pixels, default results in a 640x320 window (default 10)
  -speed int
        Speed of the emulator in cycles per tick (default 12)
  -bg-colour string
        Colour of background (default "#000000")  
  -fg-colour string
        Colour of foreground pixels (default "#22DD22")              
```

## Developing & building locally

Run `make build` and then `./bin/chip8`

```text
$ make
build                🔨 Run a local build
help                 💬 This help message :)
lint-fix             📝 Lint & format, will try to fix errors and modify code
lint                 🔍 Lint & format, will not fix but sets exit code on error
run                  🏃‍ Run application, used for local development
test                 🤡 Run those unit tests and pretend you have no bugs
```

# Repository Structure

A brief description of the top-level directories of this project are:

```text
/cmd        - Main apps, both standalone and WASM versions
/docs       - Docs, not much here
/pkg        - Go packages and modules, most of the code is here
/roms       - A handful of CHIP8 ROMs and programs for testing
/web        - Web / WASM version
```

# Known Issues

The project is NOT FINISHED

# Change Log

See [complete change log](./CHANGELOG.md)

# License

This project uses the MIT software license. See [full license file](./LICENSE)

# Acknowledgements

- https://tobiasvl.github.io/blog/write-a-chip-8-emulator/
- http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
- https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
- https://github.com/massung/CHIP-8
