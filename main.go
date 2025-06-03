package main

import (
	"log"
	"os/exec"
	"syscall"

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
			// These next three keys are local kludges
			case "KEY_1": // start Kodi
				log.Println("kodi", process(kodi))
			case "KEY_2": // start Cantata
				log.Println("cantata", process(cantata))
			case "KEY_3": // Toggle Cantata info
				log.Println("cantata info", process(cantata_info))
			default:
				log.Println("Ignored", e.CodeName())
			}
		}
	}
}

func run_cmd(cmd, ret string) string {
	child := exec.Command(cmd)
	child.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err := child.Start()
	if err != nil {
		return err.Error()
	}
	return ret
}

// This will force pause any music and start Kodi
func kodi() string {
	conn.Pause(true)
	return run_cmd("kodi", "Kodi started")
}

// This will start cantata
func cantata() string {
	return run_cmd("cantata", "Cantata started")
}

// This will toggle cantata info screen
func cantata_info() string {
	return run_cmd("cantata_info", "Info toggled")
}
