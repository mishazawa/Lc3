package opcodes

import (
  reg "github.com/mishazawa/Lc3/registers"
  mem "github.com/mishazawa/Lc3/memory"
)

const (
  BR = iota /* branch */
  ADD       /* add  */
  LD        /* load */
  ST        /* store */
  JSR       /* jump register */
  AND       /* bitwise and */
  LDR       /* load register */
  STR       /* store register */
  RTI       /* unused */
  NOT       /* bitwise not */
  LDI       /* load indirect */
  STI       /* store indirect */
  JMP       /* jump */
  RES       /* reserved (unused) */
  LEA       /* load effective address */
  TRAP      /* execute trap */
)

func sign_extend (x uint16, bits int) uint16 {
  // check is 1 in MSB
  if x >> (bits - 1) & 1 == 1 {
    x |= 0xffff << bits
  }
  return x
}

/*
  Load value from memory to register. Decode memory address that
  store memory address to actual value.
*/
func LoadIndirect (instruction uint16, registers *reg.Reg, memory *mem.Memory) {
  destination := (instruction >> 9) & 0x7
  offset      := sign_extend(instruction & 0x1ff, 9)

  /*
    Value of PC register + encoded offset [9 bit] points to memory
    location in near segment of memory (max memory address = 16 bit)
    that point to another memory location (0 - 16 bit).
  */

  val := memory.Read(memory.Read(registers.Read(reg.PC) + offset))

  registers.Write(destination, val)
  registers.UpdateFlags(destination)
}

/*
  Add values.

  Add Register:

  ADD R2 R0 R1 ; add the contents of R0 to R1 and store in R2.

  Add Immediate:

  ADD R0 R0 1 ; add 1 to R0 and store back in R0
*/
func Add (instruction uint16, registers *reg.Reg) {

  // [(0 0 0 0 instruction),  0 (0 0 0 destination), (0 0 0 source) (0 mode),  0 0 0 0]
  destination := (instruction >> 9) & 0x7
  source      := (instruction >> 6) & 0x7

  mode        := (instruction >> 5) & 0x1  // mode 1 = immediate

  val := registers.Read(source)

  if mode == 1 {
    val += sign_extend(instruction & 0x1F, 5)
  } else {
    val += registers.Read(instruction & 0x1F)
  }

  registers.Write(destination, val)
  registers.UpdateFlags(destination)
}

/*
  Bitwise AND (register or immediate).
*/
func And (instruction uint16, registers *reg.Reg) {
  destination := (instruction >> 9) & 0x7
  source      := (instruction >> 6) & 0x7

  mode        := (instruction >> 5) & 0x1  // mode 1 = immediate

  val := registers.Read(source)

  if mode == 1 {
    val &= sign_extend(instruction & 0x1F, 5)
  } else {
    val &= registers.Read(instruction & 0x1F)
  }

  registers.Write(destination, val)
  registers.UpdateFlags(destination)
}


func Branch (instruction uint16, registers *reg.Reg, memory *mem.Memory) {
  nzp := (instruction >> 9) & 0x7
  if (nzp & 0x4) == 1 || (nzp & 0x2) == 1 || (nzp & 0x1) == 1 {
    offset := sign_extend(instruction & 0x1ff, 9)
    registers.Write(reg.PC, registers.Read(reg.PC) + offset)
  }
}