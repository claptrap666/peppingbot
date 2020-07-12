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
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/spf13/cobra"
	"gitlab.eazytec-cloud.com/zhanglv/peepingbot/core"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "start to record screen stop to afile",
	Run: func(cmd *cobra.Command, args []string) {
		core.Config.FPS = 30
		core.Config.Alpha = 15
		core.Config.Quality = 75
		n := screenshot.NumActiveDisplays()
		var done [10]chan bool
		var images [10]chan *bytes.Buffer
		var shooters [10]*core.Shooter
		var converters [10]*core.FileConvertor
		for i := 0; i < n; i++ {
			done[i] = make(chan bool, 1)
			images[i] = make(chan *bytes.Buffer, 10)
			shooters[i] = &core.Shooter{
				Images: images[i],
				Done:   done[i],
			}
			converters[i] = &core.FileConvertor{
				Src:  images[i],
				Done: done[i],
			}
			converters[i].Init(fmt.Sprintf("test%d.avi", i))
			go func(index int) {
				shooters[index].Start(index)
			}(i)
			go func(index int) {
				converters[index].Start()
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
	rootCmd.AddCommand(recordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
