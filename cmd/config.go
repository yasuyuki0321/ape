package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/yasuyuki0321/ape/pkg/config"
	"github.com/yasuyuki0321/ape/pkg/utils"
)

const configFile = ".config.yaml"

var name, roleArn string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage config.yaml file",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Run: func(cmd *cobra.Command, args []string) {
		config.ListAccounts(configFile)
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an account",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		roleArn, _ := cmd.Flags().GetString("roleArn")

		if err := config.AddAccount(configFile, name, roleArn); err != nil {
			fmt.Printf("Error adding account: %s\n", err)
		}
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an account",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		if err := config.RemoveAccount(configFile, name); err != nil {
			fmt.Printf("Error removing account: %s\n", err)
		}
	},
}

var testConnectionCmd = &cobra.Command{
	Use:   "test",
	Short: "Test connection to AWS accounts",
	RunE:  testConnection,
}

func testConnection(cmd *cobra.Command, args []string) error {
	roleArns, err := utils.LoadAccountsFromConfig()
	if err != nil {
		return fmt.Errorf("error loading accounts from config: %w", err)
	}

	if target != "all" {
		roleArns = utils.FilterAccounts(roleArns, target)
	}

	if len(roleArns) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no accounts found for target: %s\n", target)
		// cobraの仕様でエラーを返すとUsageの情報が表示されることを回避するためnilを返す
		return nil
	}

	var wg sync.WaitGroup
	errorsCh := make(chan error, len(roleArns))

	for _, roleArn := range roleArns {
		wg.Add(1)
		go func(roleArn utils.Account) {
			defer wg.Done()
			if err := config.CheckAWSConnections([]utils.Account{roleArn}); err != nil {
				errorsCh <- err
			}
		}(roleArn)
	}
	wg.Wait()
	close(errorsCh)

	for err := range errorsCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configListCmd, configAddCmd, configRemoveCmd, testConnectionCmd)

	configAddCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the account to be added")
	configAddCmd.Flags().StringVarP(&roleArn, "roleArn", "r", "", "Role ARN of the account to be added")
	configAddCmd.MarkFlagRequired("name")
	configAddCmd.MarkFlagRequired("roleArn")

	configRemoveCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the account to be removed")
	configRemoveCmd.MarkFlagRequired("name")

	testConnectionCmd.Flags().StringVarP(&target, "target", "t", "all", "Command target accounts (default: all)")
}
