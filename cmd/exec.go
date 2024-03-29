package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
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
	Short: "Execute awscli command across multiple accounts",
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
		fmt.Fprintf(os.Stderr, "operation aborted\n")
		return nil
	}

	return executeAWSCommands(roleArns)
}

func executeAWSCommands(roleArns []utils.Account) error {
	ctx := context.TODO()
	var wg sync.WaitGroup

	errChan := make(chan error, len(roleArns))

	// セマフォの設定: 同時に実行したいgoroutineの数を決定
	semaphore := make(chan struct{}, 5) // 例: 同時に10個のgoroutineを実行

	for _, roleArn := range roleArns {
		wg.Add(1)

		// セマフォを取得
		semaphore <- struct{}{}

		go func(roleArn utils.Account) {
			defer wg.Done()
			defer func() { <-semaphore }() // セマフォを解放

			var outputBuffer bytes.Buffer

			creds, err := aws.AssumeRole(ctx, roleArn.RoleArn)
			if err != nil {
				errChan <- fmt.Errorf("failed to assume role %s: %s", roleArn, err)
				return
			}

			aws.SetTempCredentials(creds)
			if err := aws.ExecuteAWSCLI(&outputBuffer, roleArn, command); err != nil {
				errChan <- fmt.Errorf("error executing AWS CLI for role %s: %s", roleArn, err)
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
	execCmd.MarkFlagRequired("command")
	execCmd.Flags().BoolVarP(&skipPreview, "skip-preview", "y", false, "Skip the preview and execute the command directly")
	execCmd.Flags().StringVarP(&target, "target", "t", "all", "Command target accounts")
}
