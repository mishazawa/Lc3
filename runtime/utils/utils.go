package utils

import (
  "os"
  "github.com/eiannone/keyboard"
)

func InitKeyboard () {
  if err := keyboard.Open(); err != nil {
    panic(err)
  }
}

func CloseKeyboard () {
  keyboard.Close()
}

func GetChar () uint16 {
  char, key, err := keyboard.GetKey()

  if err != nil {
    panic(err)
  }

  if key == keyboard.KeyEsc {
    os.Exit(0)
  }

  return uint16(char)
}
