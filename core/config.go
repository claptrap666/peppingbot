package core

//Config all configs
var Config config

type config struct {
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
