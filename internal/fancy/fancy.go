package fancy

import (
	"github.com/fatih/color"
)

var (
	Bold = color.New(color.Bold, color.FgWhite).SprintFunc()
	Underline = color.New(color.Underline, color.FgWhite).SprintFunc()
	//failure = color.New()
)
