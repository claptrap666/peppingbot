/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/spf13/cobra"
	"gitlab.eazytec-cloud.com/zhanglv/peepingbot/core"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "stream screen as rtmp",
	Run: func(cmd *cobra.Command, args []string) {
		n := screenshot.NumActiveDisplays()
		var done [10]chan bool
		var images [10]chan *image.RGBA
		var shooters [10]*core.Shooter
		var converters [10]*core.RtmpConverter
		for i := 0; i < n; i++ {
			done[i] = make(chan bool, 1)
			images[i] = make(chan *image.RGBA, 10)
			shooters[i] = &core.Shooter{
				Images: images[i],
				Done:   done[i],
			}
			converters[i] = &core.RtmpConverter{
				Src:  images[i],
				Done: done[i],
			}
			go func(index int) {
				shooters[index].Start(index)
			}(i)
			go func(index int) {
				converters[index].Start(index)
			}(i)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case <-sigs:
				fmt.Println("begin shutdown......")
				for i := 0; i < n; i++ {
					done[i] <- true
				}
				time.Sleep(5 * time.Second)
				os.Exit(0)
			default:
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
