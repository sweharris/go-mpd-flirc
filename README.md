## go-mpd-flirc

I have a [FLIRC USB receiver](https://flirc.tv/products/flirc-usb-receiver)
connected to my Media PC running Linux.  This kinda does in hardware
what LIRC does in software; you can map remote control buttons to keys and
it looks just like a HID device to the OS.  So, for example, if I have a
terminal window open and press "1" on the remote then a "1" is entered
just as if I'd typed it on a keyboard.

Using this I wanted to have minimal control over MPD, similar to how
iTunes used to have remote control.

The main challenge is that there's no foreground window process with
focus so this has to run in the background.  Which then raises the challenge
of making sure it doesn't do anything if Kodi is already running.

So I make use of two libraries:

* github.com/fhs/gompd/v2  - this talks to MPD
* github.com/holoplot/go-evdev  - this reads the FLIRC `/dev/input` entry

I also adapted https://github.com/DeedleFake/ptt-fix/blob/v0.7.5/internal/xdo/xdo.go so I could see what window was active (essentially the equivalent
of `xdotool getwindowfocus getwindowname`)

So the code logic goes something like:

* Open the FLIRC
* Read a key
* If it's a key we care about, use xdo to get the current active window
* If the window is not Kodi (or empty) then process the key
* Otherwise ignore it.

The file `device.go` contains the `/dev/input` entry for the device.

The `main()` loop contains the remote keys we care about.  This is very
minimal; play/pause/stop, next/previous track, skip forwards/backwards.

On my Debian machine I can start this at login time with an entry in
`$HOME/.config/autostart/mpdflirc.desktop` containing the following

```
[Desktop Entry]
Name=MPDFlirc
Comment=MPDFlirc
Exec=/home/sweh/bin/mpd_flirc
Terminal=false
Type=Application
Encoding=UTF-8
Categories=Application;
```

In my case this is a simple wrapper script which lets me restart it whenever
I wish, so if I make changes (eg add new keys) I don't need to logout/login
again, I can just run this script.

```
#!/bin/bash

pkill mpd-flirc

export DISPLAY=:0.0 
setsid /home/sweh/src/go-mpd-flirc/mpd-flirc &
```
