package antialexa

import (
	"bytes"
	"embed"
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

//go:embed Silent1s.mp3
var assets embed.FS

func RegisterComponent() {
	mAntialexa := systray.AddMenuItemCheckbox("AntiAlexa", "Stop annoying alexa connection notifications", true)
	var timer *time.Ticker

	go func() {
		for {
			<-mAntialexa.ClickedCh
			if mAntialexa.Checked() {
				mAntialexa.Uncheck()
				handleMenuToggle(false, timer)
			} else {
				mAntialexa.Check()
				handleMenuToggle(true, timer)
			}
		}
	}()

	// Handle for the state of the check when the program was launched
	handleMenuToggle(mAntialexa.Checked(), timer)
}

func handleMenuToggle(newCheckedState bool, timer *time.Ticker) {
	if newCheckedState {
		timer = time.NewTicker(time.Minute * 5)
		defer timer.Stop()
		// Run once right away
		executeTimerTick()
		go func() {
			for range timer.C {
				executeTimerTick()
			}
		}()
	} else {
		if timer != nil {
			timer.Stop()
		}
	}
}

func executeTimerTick() {
	fmt.Println("Playing sound")
	go playSilentSound()
}

func playSilentSound() {
	// Read the mp3 file into memory
	fileBytes, err := assets.ReadFile("Silent1s.mp3")
	if err != nil {
		panic("reading embed/Silent1s.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	samplingRate := 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
	// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
	// if err != nil{
	//     panic("player.Seek failed: " + err.Error())
	// }
	// println("Player is now at position:", newPos)
	// player.Play()

	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
