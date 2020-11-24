package routines

import (
	"fmt"
  utils "github.com/mishazawa/Lc3/runtime/utils"
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
    Out(val)
    if val == 0 {
      break
    }
    pointer += 1
  }
}

func Putsp (pointer uint16, memoryRead func(uint16) uint16) {
  for {
    val := memoryRead(pointer)
    Out(val & 0xff)
    Out((val >> 8) & 0xff)
    if val == 0 {
      break
    }
    pointer += 1
  }
}

func Getc () uint16 {
  rune, _, _ := utils.GetChar()
  fmt.Println("Getc", uint16(rune))
  panic("asd")
  return uint16(rune)

}

func Out (val uint16) {
  fmt.Printf("%c", val & 0xff)
}
