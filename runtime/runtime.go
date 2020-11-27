package runtime

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	m "github.com/mishazawa/Lc3/runtime/memory"
	r "github.com/mishazawa/Lc3/runtime/registers"

	op "github.com/mishazawa/Lc3/runtime/opcodes"

	utils "github.com/mishazawa/Lc3/runtime/utils"
)

type runtime struct {
	memory      *m.Memory
	registers   *r.Registers
	running     bool
	instruction uint16
	code        int
	ioBuffer    uint16
}

func Boot() *runtime {
	utils.InitKeyboard()
	return &runtime{m.New(), r.New(), false, 0, 0, 0}
}

func (runtime *runtime) Load(file *os.File) error {
	/*
	   1. read origin 2 bytes
	   2. put it to PC register
	   3. load rest of instructions to memory from origin point
	*/

	origin := make([]byte, 2)

	_, err := file.Read(origin)

	if err != nil {
		return err
	}

	runtime.registers.Write(r.PC, binary.BigEndian.Uint16(origin))

	pointer := runtime.registers.Read(r.PC)

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
		runtime.memory.Write(pointer, code)
		pointer += 1
	}

	return nil
}

func (runtime *runtime) Run() int {
	defer utils.CloseKeyboard()

	runtime.running = true
	/*
	   1. Load one instruction from memory at the address of the PC register.
	   2. Increment the PC register.
	   3. Look at the opcode to determine which type of instruction it should perform.
	   4. Perform the instruction using the parameters in the instruction.
	   5. Go back to step 1.
	*/
	for {
		if !runtime.running {
			break
		}

		/* Program counter should contain address of NEXT instruction */
		runtime.registers.Inc(r.PC)

		/* Decrement PC to execute CURRENT instruction */
		runtime.instruction = runtime.memory.Read(runtime.registers.Read(r.PC) - 1)

		switch runtime.instruction >> 12 {
		case op.LD:
			op.Load(runtime)
		case op.LDI:
			op.LoadIndirect(runtime)
		case op.LDR:
			op.LoadRegister(runtime)
		case op.LEA:
			op.LoadEffectiveAddress(runtime)
		case op.ADD:
			op.Add(runtime)
		case op.AND:
			op.And(runtime)
		case op.NOT:
			op.Not(runtime)
		case op.BR:
			op.Branch(runtime)
		case op.JMP:
			op.Jump(runtime)
		case op.JSR:
			op.JumpRegister(runtime)
		case op.ST:
			op.Store(runtime)
		case op.STI:
			op.StoreIndirect(runtime)
		case op.STR:
			op.StoreRegister(runtime)
		case op.TRAP:
			op.Trap(runtime)
		default:
			panic(fmt.Sprintf("%016b not implemented.\n", runtime.instruction))
		}

	}
	return 0
}

func (runtime *runtime) Stop(stopCode int) {
	runtime.code = stopCode
	runtime.running = false
}

/* interface runtime */

func (runtime *runtime) ReadMemory(addr uint16) uint16 {
	return runtime.memory.Read(addr)
}

func (runtime *runtime) ReadRegister(addr uint16) uint16 {
	return runtime.registers.Read(addr)
}

func (runtime *runtime) WriteMemory(addr, val uint16) {
	runtime.memory.Write(addr, val)
}

func (runtime *runtime) WriteRegister(addr, val uint16) {
	runtime.registers.Write(addr, val)
}

func (runtime *runtime) UpdateRegisterFlags(addr uint16) {
	runtime.registers.UpdateFlags(addr)
}

func (runtime *runtime) ReadInstruction() uint16 {
	return runtime.instruction
}
