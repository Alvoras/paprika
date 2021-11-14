package cli

import (
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"log"
	"os"

	"github.com/alvoras/paprika/internal/fancy"
	"github.com/andybons/gogif"
)

func ToGif(outFilepath string, steps []string) {
	var quitChan = make(chan bool)
	var loadedFrameChan = make(chan int)
	gifColors := 256
	gifFile := &gif.GIF{}

	f, err := os.OpenFile(outFilepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	go showGifProgress(len(steps), "Loading frames...", "frames", loadedFrameChan, quitChan)
	for _, file := range steps {
		stepFile, err := os.Open(file)
		if err != nil {
			log.Fatalln(err)
		}

		frame, err := png.Decode(stepFile)
		if err != nil {
			log.Fatalln(err)
		}
		stepFile.Close()

		bounds := frame.Bounds()
		palettedFrame := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: gifColors}
		quantizer.Quantize(palettedFrame, bounds, frame, image.Point{})

		gifFile.Image = append(gifFile.Image, palettedFrame)
		gifFile.Delay = append(gifFile.Delay, GifDelay)
		loadedFrameChan <- 1
	}
	quitChan <- true

	err = gif.EncodeAll(f, gifFile)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("✨ [%s] Clip created\n", fancy.Bold(outFilepath))

}
