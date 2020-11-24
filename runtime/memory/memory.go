package memory

import (
  utils "github.com/mishazawa/Lc3/runtime/utils"
)

const (
  KBSR = 0xFE00 /* keyboard status */
  KBDR = 0xFE02 /* keyboard data */
)

type Memory struct {
	layout []uint16
}

func New () *Memory {
	return &Memory {
		layout: make([]uint16, ^uint16(0)),
	}
}

func (m *Memory) Read (addr uint16) uint16 {
	if addr == KBSR {
		if c := utils.Keypress(); c != 0 {
			m.Write(KBSR, 1 << 15)
			m.Write(KBDR, c & 0x0f)
		} else {
			m.Write(KBSR, 0x00)
			m.Write(KBDR, 0x00)
		}
	}
	return m.layout[addr]
}

func (m *Memory) Write (addr uint16, value uint16) {
	m.layout[addr] = value
}
