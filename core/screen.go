package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	screenshot "github.com/vova616/screenshot"
)

//Images xxx
var Images chan *bytes.Buffer

// Done xxx
var Done chan bool

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
	Images = make(chan *bytes.Buffer, 10)
	stdInterval := float64(1.0 / float64(FPS))
	timeToSleep := stdInterval
	s := time.Now()
	Done = make(chan bool, 1)
	for {
		select {
		case <-Done:
			return
		default:
		}
		img := CaptureScreenMust()
		sio := bytes.NewBufferString("")
		err := jpeg.Encode(sio, img, &jpeg.Options{Quality: Quality})
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

//Stopshot xxx
func Stopshot(ctx context.Context) error {
	Done <- true
	return nil
}
