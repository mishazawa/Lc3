package memory

type Memory struct {
	layout []uint16
}

func New () *Memory {
	return &Memory {
		layout: make([]uint16, ^uint16(0)),
	}
}

func (m *Memory) Read (addr uint16) uint16 {
	return m.layout[addr]
}

func (m *Memory) Write (addr uint16, value uint16) {
	m.layout[addr] = value
}
