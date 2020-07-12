package core

import (
	"bytes"

	mjpeg "gitlab.eazytec-cloud.com/zhanglv/peepingbot/core/mjpeg"
)

//Converter xxx
type Converter interface {
	Start() error
	Stop() error
}

//FileConvertor xxx
type FileConvertor struct {
	Src  chan *bytes.Buffer
	aw   mjpeg.AviWriter
	Done chan bool
}

//Init xxxxx
func (fc *FileConvertor) Init(filename string) error {
	var err error
	fc.aw, err = mjpeg.New(filename, 640, 480, Config.FPS)
	return err
}

//Start xx
func (fc *FileConvertor) Start() error {
	for {
		select {
		case <-fc.Done:
			return fc.aw.Close()
		case data := <-fc.Src:
			fc.aw.AddFrame(data.Bytes())
		}

	}
}
