package utils

import (
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


func GetChar () (rune, keyboard.Key, error) {
	return keyboard.GetKey()
}
