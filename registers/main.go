package registers

const (
  R0 = iota
  R1
  R2
  R3
  R4
  R5
  R6
  R7
  PC   /* program counter */
  COND
  COUNT
)

type Reg struct {
  layout [COUNT]uint16
}

func New () *Reg {
  return &Reg{}
}

func (r *Reg) Read (addr uint16) uint16 {
  return r.layout[addr]
}

func (r *Reg) Write (addr uint16, val uint16) {
  r.layout[addr] = val
}

func (r *Reg) Inc (addr uint16) uint16 {
  r.layout[addr] += 1
  return r.layout[addr]
}
