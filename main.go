package main

import (
  "os"
  "io"
  "flag"
  "bufio"
  "fmt"
  "regexp"
  "encoding/binary"

  reg "github.com/mishazawa/Lc3/registers"
  mem "github.com/mishazawa/Lc3/memory"
  op "github.com/mishazawa/Lc3/opcodes"
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

  err := load(memory, registers)

  if err != nil {
    panic(err)
  }

  running := true

  for {
    if (!running) { break }

    instruction := memory.Read(registers.Read(reg.PC))
    registers.Inc(reg.PC)

    switch instruction >> 12 {
    case op.LD:
      op.Load(instruction, registers, memory)
    case op.LDI:
      op.LoadIndirect(instruction, registers, memory)
    case op.LDR:
      op.LoadRegister(instruction, registers, memory)
    case op.LEA:
      op.LoadEffectiveAddress(instruction, registers)
    case op.ADD:
      op.Add(instruction, registers)
    case op.AND:
      op.And(instruction, registers)
    case op.NOT:
      op.Not(instruction, registers)
    case op.BR:
      op.Branch(instruction, registers, memory)
    case op.JMP:
      op.Jump(instruction, registers)
    case op.JSR:
      op.JumpRegister(instruction, registers)
    case op.ST:
      op.Store(instruction, registers, memory)
    case op.TRAP:
      running = op.Trap(instruction, registers, memory)
    case op.RTI:
      panic("RTI not implemented.")
    case op.RES:
    default:
      panic(fmt.Sprintf("\b not implemented.\n", instruction >> 12))
    }
  }
}


func initialize () (*reg.Reg, *mem.Memory) {
  registers := reg.New()
  memory    := mem.New()
  return registers, memory
}

var HELP_MESSAGE string = `
  help                     - print this message
  load "... path to image" - load assembler image
`


func load (memory *mem.Memory, registers *reg.Reg) error {
  var err error
  var mess string

  file := flag.String("exec", "", "Path to asm file")
  flag.Parse()

  if len(*file) != 0 {
    err, mess = loadFile(*file, memory, registers)
    if len(mess) != 0 {
      fmt.Println(mess)
      return nil
    }
    return err
  }

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
      err, mess = loadFile(groups[2], memory, registers)
      if len(mess) != 0 {
        fmt.Println(mess)
      } else {
        break load_loop
      }
    }
  }

  return err
}


func readImageToMemory (filename string, memory *mem.Memory, registers *reg.Reg) error {
  file, err := os.Open(filename)

  if err != nil {
    return err
  }

  defer file.Close()

  /*
    1. read origin 2 bytes
    2. put it to PC register
    3. load rest of instructions to memory from origin point
  */

  origin := make([]byte, 2)

  _, err = file.Read(origin)

  if err != nil {
    return err
  }
  registers.Write(reg.PC, binary.BigEndian.Uint16(origin))

  pointer := registers.Read(reg.PC)

  for {
    data := make([]byte, 2)
    _, err := file.Read(data)

    if err == io.EOF {
      return nil
    }

    if err != nil {
      return err
    }

    code := binary.BigEndian.Uint16(data)
    memory.Write(pointer, code)
    pointer += 1
  }

  return nil
}

func swap16(x uint16) uint16 {
  return (x << 8) | (x >> 8)
}

func loadFile (path string, memory *mem.Memory, registers *reg.Reg) (error, string) {
  if len(path) == 0 {
    return nil, "[Error] Enter path to file."
  }

  info, err := os.Stat(path)

  if os.IsNotExist(err) {
    return nil, "[Error] File doesn't exist."
  } else if info.IsDir() {
    return nil, fmt.Sprintf("[Error] %s is directory.", path)
  } else {
    err := readImageToMemory(path, memory, registers)
    return err, ""
  }
}
