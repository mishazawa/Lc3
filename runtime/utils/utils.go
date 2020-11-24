package utils

import (
	"context"
	"time"

	"github.com/eiannone/keyboard"
)

const shortDuration = 2000 * time.Millisecond

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

func Keypress () uint16 {
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel()

	select {
	case v := <-getCharChan():
		return uint16(v)
	case <-ctx.Done():
		return 0
	}
}

func getCharChan () chan rune {
	chanrune := make(chan rune)
  go func() {
  	r, k, err := GetChar()
  	if IsEsc(r, k, err) {
  		panic("Esc")
  	}
    chanrune <- r
  }()
  return chanrune
}
