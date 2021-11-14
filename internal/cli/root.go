package cli

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	RootCmd = &cobra.Command{
		Use:   "paprika",
		Short: "Paprika",
		Long:  "Paprika",
	}
)

func init() {
	log.SetFlags(0)
}
