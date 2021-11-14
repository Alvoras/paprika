package cli

import (
	"errors"
	"fmt"
	"github.com/alvoras/paprika/internal/fancy"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	clipCmd = &cobra.Command{
		Use:   "clip",
		Args:  cobra.ExactArgs(1),
		Short: "Create a clip from a 'steps' directory",
		Long:  "This subcommand is used create a clip from a 'steps' directory",
		Run:   createClip,
	}

	Fps        = 24
	Cut        string
	Last       int
	First      int
	StartAt    int
	EndAt      int
	Boomerang  bool
	Crossfade  bool
	Blend      bool
	ConvertGif bool
	GifDelay   int
	FadeOffset int
	FadeDuration int
	//duration    time.Duration
	OutFilepath string
)

func init() {
	RootCmd.AddCommand(clipCmd)
	clipCmd.Flags().IntVarP(&Fps, "fps", "F", 24, "Frame per second")
	clipCmd.Flags().StringVarP(&OutFilepath, "out", "o", "", "Path to save the clip to")
	clipCmd.Flags().IntVarP(&Last, "last", "l", 0, "Use the last N frames")
	clipCmd.Flags().IntVarP(&First, "first", "f", 0, "Use the first N frames")
	clipCmd.Flags().StringVarP(&Cut, "cut", "c", "", "Shorthand for the --start-at and --end-at combo. Syntax : '--cut <start>:<end>'. Example to make a clip from the 10th to the 20th frames : '--cut 10:20'")
	clipCmd.Flags().IntVarP(&StartAt, "start-at", "S", 0, "Start at the specified frame")
	clipCmd.Flags().IntVarP(&EndAt, "end-at", "E", 0, "End at the specified frame")
	clipCmd.Flags().BoolVarP(&Boomerang, "boomerang", "b", false, "Apply a boomerang effect")
	clipCmd.Flags().BoolVarP(&Crossfade, "xfade", "x", false, "Apply a crossfade effect")
	clipCmd.Flags().IntVarP(&FadeOffset, "xfade-offset", "O", -1, "Specify the time offset in milliseconds at which the crossfade starts. Default is a quarter of the total duration")
	clipCmd.Flags().IntVarP(&FadeDuration, "xfade-duration", "D", -1, "Specify the duration in milliseconds of the crossfade effect. Default is half the total duration")
	//clipCmd.Flags().BoolVarP(&Blend, "blend", "B", false, "Blend the sequence images for looping videos")
	clipCmd.Flags().BoolVarP(&ConvertGif, "gif", "g", false, "Convert to gif instead of mp4")
	clipCmd.Flags().IntVarP(&GifDelay, "gif-delay", "d", 4, "Delay to apply for each frame (gif only)")
	//clipCmd.Flags().DurationVarP(&duration, "duration", "d", 0*time.Second, "Desired clip duration. Override 'Fps' setting")
}

func createClip(_ *cobra.Command, args []string) {
	root, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatalln("Failed to resolve", root)
	}

	var steps []string
	maxFrames := getStepsCount(root)

	if EndAt == 0 {
		EndAt = maxFrames
	}

	if Cut != "" {
		chunks := strings.Split(Cut, ":")
		StartAt, err = strconv.Atoi(chunks[0])
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				log.Fatalln("--cut syntax error. Use numbers only (eg. '--cut 10:20')")
			} else {
				log.Fatalln("--cut format error. See --help for details")
			}
		}

		EndAt, err = strconv.Atoi(chunks[1])
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				log.Fatalln("--cut syntax error. Use numbers only (eg. '--cut 10:20')")
			} else {
				log.Fatalln("--cut format error. See --help for details")
			}
		}
	}

	frameArgs := []int{
		First,
		Last,
		StartAt,
		EndAt,
	}

	for _, arg := range frameArgs {
		if arg != 0 {
			if arg > maxFrames {
				log.Fatalf("Frame number must be below the max number of images (expected less than %d, found %d)\n", maxFrames, arg)
			}

			if arg < 0 {
				log.Fatalf("Frame number must be above 0 (found %d)", arg)
			}
		}
	}

	if StartAt > EndAt {
		log.Fatalln("--cut format error. First frame number cannot be bigger than Last frame number")
	}

	if First != 0 {
		EndAt = First
	}

	if Last != 0 {
		StartAt = maxFrames - Last
	}

	ext := "mp4"
	if ConvertGif {
		ext = "gif"
	}

	if Crossfade && ConvertGif {
		log.Fatalln("Crossfade effect is only compatible with mp4 export")
	}

	if len(OutFilepath) == 0 {
		OutFilepath, err = filepath.Abs(fmt.Sprintf("%s.%s", filepath.Base(root), ext))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if !strings.HasSuffix(OutFilepath, ext) {
			OutFilepath += ext
		}
	}

	for i := StartAt; i < EndAt; i++ {
		steps = append(steps, path.Join(root, fmt.Sprintf("%04d.png", i)))
	}

	if Blend {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(steps), func(i, j int) { steps[i], steps[j] = steps[j], steps[i] })
	}

	if Boomerang {
		reversed := make([]string, len(steps))
		copy(reversed, steps)

		for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
			reversed[i], reversed[j] = reversed[j], reversed[i]
		}

		// Remove Last frame of the original array to avoid duplicate frame
		steps = steps[:len(steps)-1]
		steps = append(steps, reversed...)
	}

	if ConvertGif {
		ToGif(OutFilepath, steps)
	} else {
		err = ToMP4(OutFilepath, steps)
		if err != nil {
			log.Fatalln(err)
		}

		if Crossfade {
			totalFrames := EndAt - StartAt

			err = MakeCrossfade(OutFilepath, totalFrames)
			if err != nil {
				log.Fatalln(err)
			}
		}

		fmt.Println()
		log.Printf("✨ [%s] Clip created\n", fancy.Bold(OutFilepath))
	}
}

func getStepsCount(root string) int {
	var steps []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".png") && !strings.Contains(info.Name(), "progress") {
			steps = append(steps, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	return len(steps)
}
