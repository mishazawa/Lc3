package main

import (
  "log"
  "github.com/mishazawa/Lc3/registers"
  _ "github.com/mishazawa/Lc3/opcodes"
  _ "github.com/mishazawa/Lc3/cond"
)


var mem []uint16 = make([]uint16, ^uint16(0))
var reg []uint16 = make([]uint16, registers.COUNT)

func main () {
  /*
      1. Load one instruction from memory at the address of the PC register.
      2. Increment the PC register.
      3. Look at the opcode to determine which type of instruction it should perform.
      4. Perform the instruction using the parameters in the instruction.
      5. Go back to step 1.
  */

  // load

  log.Println("Mem: ", len(memory))
  log.Println("Registers: ", len(reg))
}
