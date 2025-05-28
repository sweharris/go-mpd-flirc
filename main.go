package main

import (
	"log"

	"github.com/holoplot/go-evdev"
)

func main() {
	log.Println("Starting...")
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
				log.Println("Play/pause", process(play))
			case "KEY_X": // stop button
				log.Println("Stop", process(pause))
			case "KEY_RIGHTBRACE": // ch- == previous
				log.Println("ch-", process(prev))
			case "KEY_LEFTBRACE": // ch+ == next
				log.Println("ch+", process(next))
			case "KEY_COMMA", "KEY_R": // skip back
				log.Println("replay", process(skip_back))
			case "KEY_DOT", "KEY_F": // skip forward
				log.Println("skip", process(skip_forward))
			default:
				log.Println("Ignored", e.CodeName())
			}
		}
	}
}
