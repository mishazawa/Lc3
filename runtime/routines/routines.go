package routines

import (
	"fmt"
)

const (
  GETC  = 0x20 /* get character from keyboard, not echoed onto the terminal */
  OUT   = 0x21 /* output a character */
  PUTS  = 0x22 /* output a word string */
  IN    = 0x23 /* get character from keyboard, echoed onto the terminal */
  PUTSP = 0x24 /* output a byte string */
  HALT  = 0x25 /* halt the program */
)

func Puts (pointer uint16, memoryRead func(uint16) uint16) {
  for {
    val := memoryRead(pointer)
    fmt.Printf("%c", val)
    if val == 0 {
      break
    }
    pointer += 1
  }
}
