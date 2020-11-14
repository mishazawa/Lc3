package main

import (
  "os"
  "bufio"
  "fmt"
  "regexp"
  reg "github.com/mishazawa/Lc3/registers"
  mem "github.com/mishazawa/Lc3/memory"
  _ "github.com/mishazawa/Lc3/opcodes"
  _ "github.com/mishazawa/Lc3/cond"
)


func main () {
  registers, memory := initialize()

  /*
      1. Load one instruction from memory at the address of the PC register.
      2. Increment the PC register.
      3. Look at the opcode to determine which type of instruction it should perform.
      4. Perform the instruction using the parameters in the instruction.
      5. Go back to step 1.
  */

  load()

  running := true

  for {
    if (!running) { break }
    instruction := memory.Read(registers.Inc(reg.PC))
    fmt.Println(instruction >> 12)
    running = false
  }
}


func initialize () (*reg.Reg, *mem.Memory) {
  registers := reg.New()
  memory    := mem.New()
  // defaults
  registers.Write(reg.PC, 0x3000)
  return registers, memory
}

var HELP_MESSAGE string = `
  help                     - print this message
  load "... path to image" - load assembler image
`


func load () {
  COMMAND := regexp.MustCompile(`(load|help)\s?(.*)?$`)

  scanner := bufio.NewScanner(os.Stdin)

  load_loop:
  for {
    fmt.Printf("Lc3 > ")
    scanner.Scan()

    groups := COMMAND.FindStringSubmatch(scanner.Text())

    if len(groups) == 0 { continue load_loop }

    switch groups[1] {
    case "help":
      fmt.Println(HELP_MESSAGE)
    case "load":
      if len(groups[2]) == 0 {
        fmt.Println("[Error] Enter path to file.")
        continue load_loop
      }

      info, err := os.Stat(groups[2])

      if os.IsNotExist(err) {
        fmt.Println("[Error] File doesn't exist.")
      } else if info.IsDir() {
        fmt.Printf("[Error] %s is directory.\n", groups[2])
      } else {
        fmt.Printf("[Reading] %s\n", groups[2])
        break load_loop
      }
    }
  }
}
