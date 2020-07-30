package core

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"strings"
	"time"
)

//Converter xxx
type Converter interface {
	Start() error
	Stop() error
}

//FileConvertor xxx
type FileConvertor struct {
	Src  chan *image.RGBA
	Done chan bool
}

//Start xx
func (fc *FileConvertor) Start() error {
	cmdArgs := fmt.Sprintf("%s %s %s %s %d %s %s",
		Config.Ffmpeg.bin,
		"-y -an -f image2pipe -vcodec png -pix_fmt bgr24 -s",
		fmt.Sprintf("%dx%d", Config.Screen.Width, Config.Screen.Height),
		"-r",
		Config.Screen.FPS,
		"-i - -c:v libx264 -pix_fmt yuv420p -preset ultrafast -f flv",
		fmt.Sprintf("%s/screenshot-%d.flv", Config.Flv.dir, time.Now().Unix()),
	)
	fmt.Printf("cmdargs:%s\n", cmdArgs)
	list := strings.Split(cmdArgs, " ")
	cmd := exec.Command(list[0], list[1:]...)
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	defer cmdIn.Close()
	for {
		select {
		case <-fc.Done:
			return cmdIn.Close()
		case img := <-fc.Src:
			data := bytes.NewBufferString("")
			png.Encode(data, img)
			cmdIn.Write(data.Bytes())
		}
	}
}

//for debug
func saveToFile(img *image.RGBA) {
	f, err := os.Create(fmt.Sprintf("sceenshot-%d.png", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
