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
	src chan *bytes.Buffer
	aw  mjpeg.AviWriter
}

//Init xxxxx
func (fc *FileConvertor) Init(filename string) error {
	fc.src = Images
	var err error
	fc.aw, err = mjpeg.New(filename, 640, 480, FPS)
	return err
}

//Start xx
func (fc *FileConvertor) Start() error {
	for {
		select {
		case <-Done:
			return fc.aw.Close()
		case data := <-fc.src:
			fc.aw.AddFrame(data.Bytes())
		}

	}
}
