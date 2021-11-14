package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"

	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
)

func showGifProgress(total int, description string, iterationName string, loadedFrameChan chan int, quit chan bool) {
	var err error
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetItsString(iterationName),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]━[reset]",
			SaucerHead:    "[cyan]╸[reset]",
			SaucerPadding: "━",
			BarStart:      "",
			BarEnd:        "",
		}))

	err = bar.Set(0)
	if err != nil {
		log.Fatalln(err)
	}
	spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	//spin.Start()
	for {
		select {
		case frameLoaded := <-loadedFrameChan:
			err = bar.Add(frameLoaded)
			if err != nil {
				log.Fatalln(err)
			}
			spin.Prefix = fmt.Sprintf("%s ", bar.String())
		case <-quit:
			spin.Stop()
			err = bar.Finish()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
}

func showClipProgress(cmdStdout io.ReadCloser, total int, description string, iterationName string, quit chan bool) {
	var chunks []string
	var stats = make(map[string]string)
	var statName string
	var statVal string
	stdout := make(chan string)

	bar := progressbar.NewOptions(total,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetItsString(iterationName),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]━[reset]",
			SaucerHead:    "[cyan]╸[reset]",
			SaucerPadding: "━",
			BarStart:      "",
			BarEnd:        "",
		}))
	spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	go func() {
		scanner := bufio.NewScanner(cmdStdout)
		for scanner.Scan() {
			stdout <- scanner.Text()
		}
	}()

	spin.Start()
	for {
		select {
		case line := <-stdout:
			if strings.Contains(line, "=") {
				chunks = strings.Split(line, "=")
				statName, statVal = chunks[0], chunks[1]
				stats[statName] = statVal

				currentFrame, err := strconv.Atoi(stats["frame"])
				if err != nil {
					log.Fatalln(err)
				}

				err = bar.Set(currentFrame)
				if err != nil {
					log.Fatalln(err)
				}
				spin.Prefix = fmt.Sprintf("%s ", bar.String())
			} else {
				fmt.Println(line)
			}
		case <-quit:
			spin.Stop()
			err := bar.Finish()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
}

func showCompressionProgress(quit chan bool) {
	spin := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	spin.Prefix = "Compressing video... "
	spin.Start()
	<-quit
	spin.Stop()
}

func ExecWithProgress(ffmpegPath string, cmdArgs []string, description string, iterationName string, totalFrames int, progressQuitChan chan bool) error {
	cmd := exec.Command(ffmpegPath, cmdArgs...)
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	go showClipProgress(cmdStdout, totalFrames, description, iterationName, progressQuitChan)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	progressQuitChan <- true

	return nil
}
