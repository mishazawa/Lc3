package types

import (
	"os"
)

type Runtime interface {
  ReadMemory          (uint16)         uint16
  ReadRegister        (uint16)         uint16

  WriteMemory         (uint16, uint16)
  WriteRegister       (uint16, uint16)

  UpdateRegisterFlags (uint16)
  ReadInstruction     ()               uint16

  Load                (*os.File)       error
  Run                 ()               int
  Stop                (int)
}
