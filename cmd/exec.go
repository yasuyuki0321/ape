package cmd

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/yasuyuki0321/ape/pkg/aws"
	"github.com/yasuyuki0321/ape/pkg/utils"
)

var (
	target, command string
	skipPreview     bool
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute AWS command across multiple accounts",
	RunE:  executeCommand,
}

func executeCommand(cmd *cobra.Command, args []string) error {
	roleArns, err := utils.LoadAccountsFromConfig()
	if err != nil {
		return fmt.Errorf("error loading accounts from config: %w", err)
	}

	if target != "" {
		roleArns = utils.FilterAccounts(roleArns, target)
	}

	if !skipPreview && !utils.DisplayPreview(roleArns, command) {
		return fmt.Errorf("operation aborted")
	}

	return executeAWSCommands(roleArns)
}

func executeAWSCommands(roleArns []utils.Account) error {
	ctx := context.TODO()
	var wg sync.WaitGroup

	errChan := make(chan error, len(roleArns))

	for _, roleArn := range roleArns {
		wg.Add(1)

		go func(roleArn utils.Account) { // goroutineを開始
			defer wg.Done()
			var outputBuffer bytes.Buffer

			creds, err := aws.AssumeRole(ctx, roleArn.RoleArn)
			if err != nil {
				errChan <- fmt.Errorf("Failed to assume role %s: %s", roleArn, err)
				return
			}

			aws.SetTempCredentials(creds)
			if err := aws.ExecuteAWSCLI(&outputBuffer, roleArn, command); err != nil {
				errChan <- fmt.Errorf("Error executing AWS CLI for role %s: %s", roleArn, err)
				return
			}

			fmt.Print(outputBuffer.String())
			aws.ResetCredentials()
		}(roleArn)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute")
	execCmd.Flags().BoolVarP(&skipPreview, "skip-preview", "y", false, "Skip the preview and execute the command directly")
	execCmd.Flags().StringVarP(&target, "target", "t", "all", "Command target accounts")
}
