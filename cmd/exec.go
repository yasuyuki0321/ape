package cmd

import (
	"bytes"
	"context"
	"fmt"

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

	for _, roleArn := range roleArns {
		var outputBuffer bytes.Buffer
		creds, err := aws.AssumeRole(ctx, roleArn.RoleArn)
		if err != nil {
			fmt.Printf("Failed to assume role %s: %s\n", roleArn, err)
			continue
		}

		aws.SetTempCredentials(creds)
		if err := aws.ExecuteAWSCLI(&outputBuffer, roleArn, command); err != nil {
			fmt.Println("Error executing AWS CLI:", err)
		}

		fmt.Print(outputBuffer.String())
		aws.ResetCredentials()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute")
	execCmd.Flags().BoolVarP(&skipPreview, "skip-preview", "y", false, "Skip the preview and execute the command directly")
	execCmd.Flags().StringVarP(&target, "target", "t", "all", "Command target accounts")
}
