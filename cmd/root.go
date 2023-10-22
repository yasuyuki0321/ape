package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ape",
	Short: "aws parallel executer",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute: %v", err)
	}
}
