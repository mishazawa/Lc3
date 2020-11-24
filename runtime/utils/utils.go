package utils

import (
	"context"
	"time"

	"github.com/eiannone/keyboard"
)

const shortDuration = 0 * time.Millisecond

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

func IsEsc (c rune, key keyboard.Key, err error) bool {
	return key == keyboard.KeyEsc
}

func Keypress () bool {
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel()

	select {
	case <-getCharChan():
		return true
	case <-ctx.Done():
		return false
	}
}

func getCharChan () chan rune {
	chanrune := make(chan rune)
  go func() {
  	r, _, _ := GetChar()
    chanrune <- r
  }()
  return chanrune
}
