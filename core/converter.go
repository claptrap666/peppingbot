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

//Start xx
func (fc *FileConvertor) Start() error {
	aw, err := mjpeg.New("test.avi", 200, 100, 2)
	fc.aw = aw
	if err != nil {
		return err
	}
	for data := range fc.src {
		aw.AddFrame(data.Bytes())
	}
	return nil
}

//Stop xxx
func (fc *FileConvertor) Stop() error {
	return fc.aw.Close()
}
