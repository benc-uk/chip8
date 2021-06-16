//
// Console - A colored console / stdout helper
// Ben C, June 2021
// Notes:
//

package console

import (
	"fmt"
	"os"
	"strings"
)

// Debug outputs strings plus newline in magenta
func Debug(s string) {
	fmt.Printf("\033[1;35m%s\033[0m\n", s)
}

// Debugf outputs formatted strings in magenta
func Debugf(f string, a ...interface{}) {
	fmt.Printf("\033[1;35m"+f+"\033[0m", a...)
}

// Info outputs strings plus newline in blue
func Info(s string) {
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}
	fmt.Printf("\033[1;34m%s\033[0m\n", s)
}

// Infof outputs formatted strings in blue
func Infof(f string, a ...interface{}) {
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}
	fmt.Printf("\033[1;34m"+f+"\033[0m", a...)
}

// Error outputs strings plus newline in blue
func Error(s string) {
	fmt.Printf("\033[1;31m%s\033[0m\n", s)
}

// Errorf outputs formatted strings in red
func Errorf(f string, a ...interface{}) {
	fmt.Printf("\033[1;31m"+f+"\033[0m", a...)
}

// Warning outputs strings plus newline in blue
func Warning(s string) {
	fmt.Printf("\033[1;33m%s\033[0m\n", s)
}

// Warningf outputs formatted strings in yellow
func Warningf(f string, a ...interface{}) {
	fmt.Printf("\033[1;33m"+f+"\033[0m", a...)
}

// Success outputs strings plus newline in blue
func Success(s string) {
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}
	fmt.Printf("\033[1;32m%s\033[0m\n", s)
}

// Successf outputs formatted strings in green
func Successf(f string, a ...interface{}) {
	if strings.HasSuffix(os.Args[0], ".test") {
		return
	}
	fmt.Printf("\033[1;32m"+f+"\033[0m", a...)
}
