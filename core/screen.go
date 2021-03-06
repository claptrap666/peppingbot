package core

import (
	"context"
	"fmt"
	"image"
	"time"

	screenshot "github.com/kbinani/screenshot"
)

//Shooter xxx
type Shooter struct {
	Images chan *image.RGBA
	Done   chan bool
}

//CaptureScreen xxx
func CaptureScreen(rect image.Rectangle) *image.RGBA {
	img, err := screenshot.CaptureRect(rect)
	errN := 0
	for err != nil {
		fmt.Printf("CaptureScreen Error: %v", err)
		img, err = screenshot.CaptureRect(rect)
		time.Sleep(time.Duration(10)*time.Millisecond + time.Duration(errN*50)*time.Millisecond)
		if errN < 20 {
			errN++
		}
	}
	return img
}

//Start xxx
func (st *Shooter) Start(displayindex int) {
	bounds := screenshot.GetDisplayBounds(displayindex)
	stdInterval := float64(1.0 / float64(Config.Screen.FPS))
	timeToSleep := stdInterval
	s := time.Now()
	for {
		select {
		case <-st.Done:
			return
		default:
		}
		img := CaptureScreen(bounds)
		st.Images <- img
		ss := time.Now()
		interval := ss.Sub(s)
		if interval.Seconds() < stdInterval {
			timeToSleep += float64(Config.Screen.Alpha) * timeToSleep / float64(100)
		} else {
			timeToSleep -= float64(Config.Screen.Alpha) * timeToSleep / float64(100)
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

//Stop xxx
func (st *Shooter) Stop(ctx context.Context) error {
	st.Done <- true
	return nil
}
