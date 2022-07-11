package cmd

import (
	"github.com/foxdex/ftx-site/pkg/log"
	"github.com/spf13/cobra"
)

var ApiCmd = &cobra.Command{
	Use:  "api",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		setup()
		resourceRelease()
	},
}

func setup() {
	log.InitLog()
}

func resourceRelease() {

}
