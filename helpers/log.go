package helpers

import (
	"github.com/spf13/cobra"
	"log"
)

func Debug(cmd *cobra.Command, args ...any) {
	v, _ := cmd.Flags().GetBool("verbose")
	if v {
		log.Println(args...)
	}
}
