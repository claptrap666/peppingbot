package core

//Config all configs
var Config *config = &config{
	Screen: screen{},
	Rtmp:   rtmp{},
	Ffmpeg: ffmpeg{},
	Flv:    flv{},
}

type config struct {
	Screen screen
	Rtmp   rtmp
	Ffmpeg ffmpeg
	Flv    flv
}

type screen struct {
	Quality      int
	FPS          int32
	Alpha        int
	Left         int
	Top          int
	Width        int
	Height       int
	ResizeWidth  int
	ResizeHeight int
	ToSBS        bool
	Cursor       bool
	FullScreen   bool
	Convert      int
	WindowID     int64
}

type ffmpeg struct {
	bin string
}

type rtmp struct {
	push string
	key  string
}

type flv struct {
	dir string
}
