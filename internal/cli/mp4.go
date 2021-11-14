package cli

import (
	"fmt"
	"github.com/anmitsu/go-shlex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func ToMP4(outFilepath string, steps []string) error {
	ffmpegPath := "ffmpeg"
	var cmdArgs = []string{
		"-y",
		"-progress",
		"-",
		"-nostats",
		"-f",
		"image2pipe",
		"-vcodec",
		"png",
		"-framerate",
		strconv.Itoa(Fps),
		"-i",
		"-",
		"-vcodec",
		"libx264",
		"-framerate",
		strconv.Itoa(Fps),
		"-pix_fmt",
		"yuv420p",
		"-crf",
		"17",
		"-preset",
		"veryslow",
	}

	progressQuitChan := make(chan bool)
	compressionSpinQuitChan := make(chan bool)

	cmdArgs = append(cmdArgs, outFilepath)
	cmd := exec.Command(ffmpegPath, cmdArgs...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	totalFrames := EndAt - StartAt
	if Boomerang {
		totalFrames *= 2
	}
	go showClipProgress(cmdStdout, totalFrames, "Loading frames...", "frames", progressQuitChan)

	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	for _, file := range steps {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(cmdStdin, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	progressQuitChan <- true
	go showCompressionProgress(compressionSpinQuitChan)

	err = cmdStdin.Close()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	compressionSpinQuitChan <- true
	return nil
}

func SplitMP4(outFilepath string, tempDest string, totalFrames int) error {
	progressQuitChan := make(chan bool)

	ffmpegPath := "ffmpeg"
	var cmdArgs = fmt.Sprintf("-y -i %s -nostats -ss %%dms -t %%dms %s/%%d.mp4", outFilepath, tempDest)

	msPerFrame := 1000 / Fps
	totalDurationMs := totalFrames * msPerFrame
	endFirstHalf := totalDurationMs / 2
	startSecondHalf := (totalDurationMs - endFirstHalf) + 1

	firstHalfRunningCmdArgs := fmt.Sprintf(cmdArgs, 0, endFirstHalf, 1)
	secondHalfRunningCmdArgs := fmt.Sprintf(cmdArgs, startSecondHalf, totalDurationMs, 2)

	firstCmdArgsChunks, err := shlex.Split(firstHalfRunningCmdArgs, false)
	if err != nil {
		return err
	}
	secondCmdArgsChunks, err := shlex.Split(secondHalfRunningCmdArgs, false)
	if err != nil {
		return err
	}

	err = ExecWithProgress(ffmpegPath, firstCmdArgsChunks, "Splitting first half...", "frames", totalFrames/2, progressQuitChan)
	if err != nil {
		return err
	}

	err = ExecWithProgress(ffmpegPath, secondCmdArgsChunks, "Splitting second half...", "frames", totalFrames/2, progressQuitChan)
	if err != nil {
		return err
	}

	return nil
}

func ApplyCrossfade(outFilepath string, tempDest string, totalFrames int) error {
	progressQuitChan := make(chan bool)
	msPerFrame := 1000 / Fps
	totalDurationMs := totalFrames * msPerFrame
	ffmpegPath := "ffmpeg"
	fadeOffset := totalDurationMs/4
	fadeDuration := (totalDurationMs/2)-1000

	if FadeOffset != -1{
		fadeOffset = FadeOffset
	}

	if FadeDuration != -1{
		fadeDuration = FadeDuration
	}

	var cmdArgs = fmt.Sprintf("-y -i %s/2.mp4 -i %s/1.mp4 -nostats -filter_complex xfade=offset=%dms:duration=%dms -vcodec libx264 -crf 17 -pix_fmt yuv420p %s", tempDest, tempDest, fadeOffset, fadeDuration, outFilepath)

	chunkCmdArgs, err := shlex.Split(cmdArgs, false)
	if err != nil {
		return err
	}

	err = exec.Command(ffmpegPath, chunkCmdArgs...).Run()
	if err != nil {
		return err
	}

	err = ExecWithProgress(ffmpegPath, chunkCmdArgs, fmt.Sprintf("Applying crossfade... (Offset : %dms, Duration : %dms)", fadeOffset, fadeDuration), "frames", totalFrames/2, progressQuitChan)
	if err != nil {
		return err
	}

	return nil
}

func MakeCrossfade(outFilepath string, totalFrames int) error {
	tempDir, err := ioutil.TempDir("/tmp", "paprika_")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempDir)

	err = SplitMP4(outFilepath, tempDir, totalFrames)
	if err != nil {
		return err
	}

	err = ApplyCrossfade(outFilepath, tempDir, totalFrames)
	if err != nil {
		return err
	}

	return nil
}
