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
		if utils.Keypress() {
			m.Write(KBSR, 1 << 15)
			c, _, _ := utils.GetChar()
			m.Write(KBDR, uint16(c))
		} else {
			m.Write(KBSR, 0)
		}
	}
	return m.layout[addr]
}

func (m *Memory) Write (addr uint16, value uint16) {
	m.layout[addr] = value
}
