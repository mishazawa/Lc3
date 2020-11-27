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
		kbsr  := m.layout[KBSR]
		ready := (kbsr & 0x8000) == 0;

    button := utils.GetChar()

		if ready && button != 0 {
			m.Write(KBSR, kbsr   | 0x8000)
			m.Write(KBDR, button & 0x00ff)
		} else {
			m.Write(KBSR, 0x0)
		}
	}
	return m.layout[addr]
}

func (m *Memory) Write (addr uint16, value uint16) {
	m.layout[addr] = value
}
