package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
)

//RtmpConverter push to rtmp server with configured address
type RtmpConverter struct {
	Src  chan *image.RGBA
	Done chan bool
}

//Start push to livego server
func (mc *RtmpConverter) Start(screenid int) error {
	t := template.Must(template.New("keyurl").Parse(Config.Rtmp.key))
	keyurl := bytes.NewBufferString("")
	t.Execute(keyurl, &struct{ channel string }{channel: fmt.Sprintf("screen%d", screenid)})
	resp, err := http.DefaultClient.Get(keyurl.String())
	if err != nil {
		panic(err)
	}
	s, _ := ioutil.ReadAll(resp.Body)
	keyobj := struct{ data string }{}
	json.Unmarshal(s, keyobj)
	t = template.Must(template.New("push").Parse(Config.Rtmp.push))
	pushurl := bytes.NewBufferString("")
	t.Execute(pushurl, &struct{ channelkey string }{channelkey: keyobj.data})
	cmdArgs := fmt.Sprintf("%s %s %s %s %d %s %s",
		Config.Ffmpeg.bin,
		" -y -an -f image2pipe -vcodec png -pix_fmt bgr24 -s",
		fmt.Sprintf("%dx%d", Config.Screen.Width, Config.Screen.Height),
		"-r",
		Config.Screen.FPS,
		"-i - -c:v libx264 -pix_fmt yuv420p -preset ultrafast -f flv",
		pushurl,
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
		case <-mc.Done:
			return cmdIn.Close()
		case img := <-mc.Src:
			data := bytes.NewBufferString("")
			png.Encode(data, img)
			cmdIn.Write(data.Bytes())
		}
	}
}
