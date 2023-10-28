/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage config.yaml file",
	Run:   executeConfig,
}

func executeConfig(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(".config.yaml")

	if configBytes, err := os.ReadFile(configPath); err == nil {
		fmt.Println(string(configBytes))
	} else {
		fmt.Printf("Error reading config file: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.PersistentFlags().String("list", "", "A help for foo")
}
