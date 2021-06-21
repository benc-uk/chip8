//
// CHIP-8 - External audio sound (a single bleep!) used by emulator
// Ben C, June 2021
// Notes:
//

package sounds

import _ "embed"

// We use the embed package to embed the wav into the binary

//go:embed bleep.wav
var BleepWav []byte
