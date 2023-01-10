package main

import (
	"fmt"
	"os"
	"strings"
	"webpfex/webpfex"
)

var HELP string = strings.TrimSpace(`
Usage: webpfex extract AWEBP OUTDIR
       webpfex convert AWEBP OUTMP4

webpfex extracts frames from an animated WEBP or convert them to MP4. Relies on 
webpmux and ffmpeg.
`)

func main() {
	args := os.Args[1:]
	if len(args) == 3 {
		switch args[0] {
		case "extract":
			webp := args[1]
			outdir := args[2]

			err := webpfex.ExtractWebpFramesAsPng(webp, outdir)
			if err != nil {
				panic(err)
			}
		case "convert":
			webp := args[1]
			out := args[2]

			err := webpfex.ConvertWebpToMp4(webp, out)
			if err != nil {
				panic(err)
			}
		}

		return
	}

	fmt.Println(HELP)
}
