package cli

import (
	"github.com/alvoras/paprika/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alvoras/paprika/internal/fancy"

	"github.com/spf13/cobra"
)

var (
	extractCmd = &cobra.Command{
		Use:   "extract",
		Args:  cobra.ExactArgs(1),
		Short: "Extract an image every N step and create a new directory with them",
		Long:  "This subcommand is used create a clip from the steps of a Dreamstation run",
		Run:   extract,
	}

	stepOffset int
	outputDir  string
)

func init() {
	RootCmd.AddCommand(extractCmd)
	extractCmd.Flags().IntVarP(&stepOffset, "step", "s", 50, "Save an image every N pictures")
	extractCmd.Flags().StringVarP(&outputDir, "out", "o", "./extracted", "Output directory")
}

func extract(_ *cobra.Command, args []string) {
	log.Printf("Extracting to %s\n", fancy.Bold(outputDir))

	var err error
	root := args[0]
	var steps = make(map[string][]string)
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		log.Fatalln(err)
	}

	err = filepath.Walk(root, func(currentPath string, info os.FileInfo, err error) error {
		if info.IsDir() && currentPath != root {
			// Compatible with the output format of the Dreamstation, with and without the seed
			if !strings.Contains(currentPath, "it_") && !strings.HasSuffix(currentPath, "it") {
				return filepath.SkipDir
			}
		} else {
			if strings.HasSuffix(info.Name(), ".png") && !strings.Contains("progress", info.Name()) {
				absPath, err := filepath.Abs(currentPath)
				if err != nil {
					log.Fatalln(err)
				}

				dirPath := filepath.Base(filepath.Dir(absPath))
				steps[dirPath] = append(steps[dirPath], absPath)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalln("Failed to walk a directory :", err)
	}

	rootLine := "│"
	count := 0

	log.Println()
	log.Println(fancy.Bold(outputDir))

	for dirPath, files := range steps {
		dirRootLine := "├"

		if count == (len(steps) - 1) {
			rootLine = " "
			dirRootLine = "└"
		}

		log.Printf(fancy.Bold("%s── %s"), dirRootLine, dirPath)

		dstDir := filepath.Join(outputDir, dirPath)
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}

		for idx, file := range files {
			if idx == len(files)-1 || (idx+1)%stepOffset == 0 {
				dst := filepath.Join(dstDir, filepath.Base(file))
				err = utils.CopyFile(file, dst)
				if err != nil {
					log.Fatalln(err)
				}
				treePrefix := rootLine + "   ├──"
				if idx+stepOffset > len(files) {
					treePrefix = rootLine + "   └──"
				}
				log.Printf("%s%s ✨\n", treePrefix, filepath.Base(file))
			}
		}

		count++
	}

}
