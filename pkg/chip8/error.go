//
// CHIP-8 - Error codes and custom error type
// Ben C, June 2021
// Notes:
//

package chip8

const errorCodeOther = 50
const errorCodeAddress = 51
const errorBadOpcode = 52

type SystemError struct {
	code   int
	reason string
}

func (se SystemError) Error() string {
	return se.reason
}

func (se SystemError) Code() int {
	return se.code
}
