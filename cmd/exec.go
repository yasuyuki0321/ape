package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/yasuyuki0321/ape/pkg/aws"
	"github.com/yasuyuki0321/ape/pkg/utils"

	"github.com/spf13/cobra"
)

var target_account, command string
var skipPreview bool

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute aws command across multiple accounts",
	Run:   executeCommand,
}

func executeCommand(cmd *cobra.Command, args []string) {

	roleArns, err := utils.LoadAccountsFromConfig()
	if err != nil {
		fmt.Println("Error loading accounts from config:", err)
		return
	}

	if !skipPreview {
		if !utils.DisplayPreview(roleArns, command) {
			fmt.Println("Operation aborted.")
			return
		}
	}

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

	fmt.Println("finish")
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute")
	execCmd.Flags().BoolVarP(&skipPreview, "skip-preview", "y", false, "skip the preview and execute the command directly")
}
