package core

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	ljpeg "github.com/pixiv/go-libjpeg/jpeg"
	screenshot "github.com/vova616/screenshot"
)

//Images xxx
var Images chan *bytes.Buffer

// Done xxx
var Done bool = false

//CaptureScreenMust xxx
func CaptureScreenMust() *image.RGBA {
	img, err := screenshot.CaptureScreen()
	errN := 0
	for err != nil {
		fmt.Printf("CaptureWindowMust Error: %v", err)
		img, err = screenshot.CaptureScreen()
		time.Sleep(time.Duration(10)*time.Millisecond + time.Duration(errN*50)*time.Millisecond)
		if errN < 20 {
			errN++
		}
	}
	return img
}

//StartShot xxx
func StartShot() {

	fc := &FileConvertor{
		src: Images,
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fc.Stop()
	stdInterval := float64(1.0 / float64(FPS))
	timeToSleep := stdInterval
	s := time.Now()
	for {
		if Done {
			return
		}
		img := CaptureScreenMust()
		sio := bytes.NewBufferString("")
		err := ljpeg.Encode(sio, img, &ljpeg.EncoderOptions{Quality: Quality})
		if err == nil {
			Images <- sio
		}
		ss := time.Now()
		interval := ss.Sub(s)
		if interval.Seconds() < stdInterval {
			timeToSleep += float64(Alpha) * timeToSleep / float64(100)
		} else {
			timeToSleep -= float64(Alpha) * timeToSleep / float64(100)
		}
		if timeToSleep < float64(0) {
			timeToSleep = float64(0)
		}
		s = ss
		sleepTime, _ := time.ParseDuration(fmt.Sprintf("%fs", timeToSleep))
		time.Sleep(sleepTime)
		fmt.Printf("sleep: %v", sleepTime)
	}

}
