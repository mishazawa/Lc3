package opcodes

import (
  r    "github.com/mishazawa/Lc3/runtime/types"
  reg  "github.com/mishazawa/Lc3/runtime/registers"
  trap "github.com/mishazawa/Lc3/runtime/routines"
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
  Add values.

  Add Register:

  ADD R2 R0 R1 ; add the contents of R0 to R1 and store in R2.

  Add Immediate:

  ADD R0 R0 1 ; add 1 to R0 and store back in R0
*/
func Add (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  source      := (instruction >> 6) & 0x7

  mode        := (instruction >> 5) & 0x1  // mode 1 = immediate

  val := rt.ReadRegister(source)

  if mode == 1 {
    val += sign_extend(instruction & 0x1F, 5)
  } else {
    val += rt.ReadRegister(instruction & 0x1F)
  }

  rt.WriteRegister(destination, val)
  rt.UpdateRegisterFlags(destination)
}

/*
  Bitwise AND (register or immediate).
*/
func And (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  source      := (instruction >> 6) & 0x7
  mode        := (instruction >> 5) & 0x1  // mode 1 = immediate

  val := rt.ReadRegister(source)

  if mode == 1 {
    val &= sign_extend(instruction & 0x1F, 5)
  } else {
    val &= rt.ReadRegister(instruction & 0x1F)
  }

  rt.WriteRegister(destination, val)
  rt.UpdateRegisterFlags(destination)
}

/*
  Branch according to states of flags instruction[9:11].

  if neg | zero | positive then
    PC += instruction[0:8]
*/
func Branch (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  nzp := (instruction >> 9) & 0x7
  if (nzp & 0x4) == 1 || (nzp & 0x2) == 1 || (nzp & 0x1) == 1 {
    offset := sign_extend(instruction & 0x1ff, 9)
    rt.WriteRegister(reg.PC, rt.ReadRegister(reg.PC) + offset)
  }
}

/*
  Jump or return.

  PC <- instruction[9:11]
*/
func Jump (rt r.Runtime) {
  instruction := rt.ReadInstruction()
  rt.WriteRegister(reg.PC, rt.ReadRegister((instruction >> 6) & 0x7))
}

func JumpRegister (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  rt.WriteRegister(reg.R7, rt.ReadRegister(reg.PC))

  mode := (instruction >> 11) & 1

  if mode == 0 {
    value := (instruction >> 6) & 0x7
    rt.WriteRegister(reg.PC, value)
  } else {
    offset := sign_extend(instruction & 0x7ff, 11)
    rt.WriteRegister(reg.PC, rt.ReadRegister(reg.PC) + offset)
  }

}

func Load (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  offset      := sign_extend(instruction & 0x1ff, 9)

  rt.WriteRegister(destination, rt.ReadMemory(rt.ReadRegister(reg.PC) + offset))
  rt.UpdateRegisterFlags(destination)
}

/*
  Load value from memory to register. Decode memory address that
  store memory address to actual value.
*/
func LoadIndirect (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  offset      := sign_extend(instruction & 0x1ff, 9)

  /*
    Value of PC register + encoded offset [9 bit] points to memory
    location in near segment of memory (max memory address = 16 bit)
    that point to another memory location (0 - 16 bit).
  */

  val := rt.ReadMemory(rt.ReadMemory(rt.ReadRegister(reg.PC) + offset))

  rt.WriteRegister(destination, val)
  rt.UpdateRegisterFlags(destination)
}

func LoadRegister (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  base        := (instruction >> 6) & 0x7
  offset      := sign_extend(instruction & 0x2f, 6)

  rt.WriteRegister(destination, rt.ReadMemory(base + offset))
  rt.UpdateRegisterFlags(destination)
}

func LoadEffectiveAddress (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  offset      := sign_extend(instruction & 0x1ff, 9)

  rt.WriteRegister(destination, rt.ReadRegister(reg.PC) + offset)
  rt.UpdateRegisterFlags(destination)
}

func Not (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  destination := (instruction >> 9) & 0x7
  source      := (instruction >> 6) & 0x7

  rt.WriteRegister(destination, ^rt.ReadRegister(source))
  rt.UpdateRegisterFlags(destination)
}

func Store (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  source := (instruction >> 9) & 0x7
  offset := sign_extend(instruction & 0x1ff, 9)
  rt.WriteMemory(rt.ReadRegister(reg.PC) + offset, rt.ReadRegister(source))
}

func StoreIndirect (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  source := (instruction >> 9) & 0x7
  offset := sign_extend(instruction & 0x1ff, 9)
  rt.WriteMemory(rt.ReadMemory(rt.ReadRegister(reg.PC) + offset),  rt.ReadRegister(source))
}

func StoreRegister (rt r.Runtime) {
  instruction := rt.ReadInstruction()

  source := (instruction >> 9) & 0x7
  base   := (instruction >> 6) & 0x7
  offset := sign_extend(instruction & 0x2f, 6)
  rt.WriteMemory(rt.ReadMemory(rt.ReadRegister(base) + offset),  rt.ReadRegister(source))
}

func Trap (rt r.Runtime) {
  instruction := rt.ReadInstruction()
  switch instruction & 0xff {
  case trap.GETC:
    rt.WriteRegister(reg.R0, trap.Getc())
  case trap.OUT:
    trap.Out(rt.ReadRegister(reg.R0))
  case trap.PUTS:
    pointer := rt.ReadRegister(reg.R0)
    trap.Puts(pointer, rt.ReadMemory)
  case trap.IN:
    char := trap.Getc()
    rt.WriteRegister(reg.R0, char)
    trap.Out(char)
  case trap.PUTSP:
    pointer := rt.ReadRegister(reg.R0)
    trap.Puts(pointer, rt.ReadMemory)
  case trap.HALT:
    rt.Stop(0)
  }
}
