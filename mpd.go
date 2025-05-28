package main

import (
	"os"
	"strconv"
	"time"

	"github.com/fhs/gompd/v2/mpd"
)

// Global for simplicity
var conn *mpd.Client

func get_mpd_addr() (string, string) {
	// If MPD_HOST isn't defined, use the local socket
	addr := os.Getenv("MPD_HOST")
	if addr == "" {
		return "unix", "/var/run/mpd/socket"
	}

	// If it begines with a / then assume it's a local socket
	if addr[0] == '/' {
		return "unix", addr
	}

	port := os.Getenv("MPD_PORT")
	if len(port) == 0 {
		port = "6600"
	}
	return "tcp", addr + ":" + port
}

// Connect to MPD server
func connect_to_mpd() {
	if conn == nil {
		proto, addr := get_mpd_addr()
		c, err := mpd.Dial(proto, addr)
		if err != nil {
			die(err)
		}
		conn = c
	}
}

func get_status() mpd.Attrs {
	connect_to_mpd()
	status, err := conn.Status()
	if err != nil {
		die(err)
	}

	return status
}

func get_play() string {
	status := get_status()
	return status["state"]
}

func get_song() mpd.Attrs {
	connect_to_mpd()
	song, err := conn.CurrentSong()
	if err != nil {
		die(err)
	}

	return song
}

func close_mpd() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func process(f func() string) string {
	if !ignore_it() {
		connect_to_mpd()
		s := f()
		close_mpd()
		return ("--> " + s)
	}
	return "Ignored, wrong focus"
}

func play() string {
	state := get_play()

	// If there are no arguments, just tell the daemon to play
	// If we've stopped we need to start the find the current song
	// and start that, otherwise turn off pause
	if state == "stop" {
		song := get_song()
		p, _ := strconv.Atoi(song["Pos"])
		conn.Play(p)
	} else if state == "pause" {
		conn.Pause(false)
	} else {
		conn.Pause(true)
	}
	return get_play()
}

func pause() string {
	conn.Pause(true)
	return get_play()
}

// Act like a CD player; if we're within 3 seconds then skip back to
// previous song, else go to beginning of track
func prev() string {
	status := get_status()
	t, _ := strconv.ParseFloat(status["elapsed"], 64)
	if t > 3 {
		conn.SeekCur(0, false)
	} else {
		conn.Previous()
	}
	return get_play()
}

func next() string {
	conn.Next()
	return get_play()
}

func skip_forward() string {
	conn.SeekCur(time.Duration(30*time.Second), true)
	return get_play()
}

func skip_back() string {
	conn.SeekCur(time.Duration(-5*time.Second), true)
	return get_play()
}
