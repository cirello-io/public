// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"time"

	term "github.com/nsf/termbox-go"
)

func main() {
	if err := term.Init(); err != nil {
		log.Fatal(err)
	}

	defer term.Close()
	kb := make(chan string)
	go func() {
		defer close(kb)
		observeKeyboard(kb)
	}()

	tick := time.Tick(1 * time.Second / 60) // 60hz frame rate
	var currentElement string
	for {
		select {
		case event, ok := <-kb:
			if !ok {
				return
			}
			currentElement = event
		case <-tick:
			term.Clear(term.ColorDefault, term.ColorDefault)
			for i := 0; i < len(currentElement); i++ {
				term.SetCell(i, 0, rune(currentElement[i]), term.ColorDefault, term.ColorDefault)
			}
			term.Sync()
		}
	}
}

func observeKeyboard(kb chan string) {
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyF1:
				kb <- ("F1 pressed")
			case term.KeyF2:
				kb <- ("F2 pressed")
			case term.KeyF3:
				kb <- ("F3 pressed")
			case term.KeyF4:
				kb <- ("F4 pressed")
			case term.KeyF5:
				kb <- ("F5 pressed")
			case term.KeyF6:
				kb <- ("F6 pressed")
			case term.KeyF7:
				kb <- ("F7 pressed")
			case term.KeyF8:
				kb <- ("F8 pressed")
			case term.KeyF9:
				kb <- ("F9 pressed")
			case term.KeyF10:
				kb <- ("F10 pressed")
			case term.KeyF11:
				kb <- ("F11 pressed")
			case term.KeyF12:
				kb <- ("F12 pressed")
			case term.KeyInsert:
				kb <- ("Insert pressed")
			case term.KeyDelete:
				kb <- ("Delete pressed")
			case term.KeyHome:
				kb <- ("Home pressed")
			case term.KeyEnd:
				kb <- ("End pressed")
			case term.KeyPgup:
				kb <- ("Page Up pressed")
			case term.KeyPgdn:
				kb <- ("Page Down pressed")
			case term.KeyArrowUp:
				kb <- ("Arrow Up pressed")
			case term.KeyArrowDown:
				kb <- ("Arrow Down pressed")
			case term.KeyArrowLeft:
				kb <- ("Arrow Left pressed")
			case term.KeyArrowRight:
				kb <- ("Arrow Right pressed")
			case term.KeySpace:
				kb <- ("Space pressed")
			case term.KeyBackspace:
				kb <- ("Backspace pressed")
			case term.KeyEnter:
				kb <- ("Enter pressed")
			case term.KeyTab:
				kb <- ("Tab pressed")
			default:
				// we only want to read a single character or one key pressed event
				kb <- fmt.Sprint("ASCII : ", ev.Ch)

			}
		case term.EventError:
			log.Fatal(ev.Err)
		}
	}
}
