package main

import (
  "log"
  reg "github.com/mishazawa/Lc3/registers"
  mem "github.com/mishazawa/Lc3/memory"
  _ "github.com/mishazawa/Lc3/opcodes"
  _ "github.com/mishazawa/Lc3/cond"
)


func main () {
  //
  registers := reg.New()
  memory    := mem.New()
  // defaults
  registers.Write(reg.PC, 0x3000)

  /*
      1. Load one instruction from memory at the address of the PC register.
      2. Increment the PC register.
      3. Look at the opcode to determine which type of instruction it should perform.
      4. Perform the instruction using the parameters in the instruction.
      5. Go back to step 1.
  */

  // load

  running := true

  for {
    if (!running) { break }
    instruction := memory.Read(registers.Inc(reg.PC))
    log.Println(instruction >> 12)
    running = false
  }
}
