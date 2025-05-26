package main

import (
	"github.com/holoplot/go-evdev"
)

func main() {
	dev, err := evdev.Open(flirc_dev)
	if err != nil {
		die(err)
	}

	// Now loop forever reading from the Flirc
	for {
		e, err := dev.ReadOne()
		if err != nil {
			die(err)
		}

		// Value 1 == down, 0 == up, 2 == hold
		if e.Type == evdev.EV_KEY && e.Value == 1 {
			switch e.CodeName() {
				case "KEY_PLAYPAUSE":
					process(play)
				case "KEY_X":		// stop button
					process(play)
				case "KEY_RIGHTBRACE":	// ch- == previous
					process(prev)
				case "KEY_LEFTBRACE":	// ch+ == next
					process(next)
				case "KEY_COMMA", "KEY_R":	// skip back
					process(skip_back)
				case "KEY_DOT", "KEY_F":	// skip forward
					process(skip_forward)
			}
		}
	}
}
