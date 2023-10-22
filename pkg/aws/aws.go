package aws

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/yasuyuki0321/ape/pkg/utils"
)

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
	AWS_SESSION_TOKEN     = "AWS_SESSION_TOKEN"
)

func SetTempCredentials(creds *sts.AssumeRoleOutput) {
	os.Setenv(AWS_ACCESS_KEY_ID, *creds.Credentials.AccessKeyId)
	os.Setenv(AWS_SECRET_ACCESS_KEY, *creds.Credentials.SecretAccessKey)
	os.Setenv(AWS_SESSION_TOKEN, *creds.Credentials.SessionToken)
}

func ResetCredentials() {
	os.Unsetenv(AWS_ACCESS_KEY_ID)
	os.Unsetenv(AWS_SECRET_ACCESS_KEY)
	os.Unsetenv(AWS_SESSION_TOKEN)
}

func ExecuteAWSCLI(outputBuffer *bytes.Buffer, account utils.Account, command string) error {
	cmd := strings.Fields(command)
	execCmd := exec.Command(cmd[0], cmd[1:]...)

	var out bytes.Buffer
	execCmd.Stdout = &out

	if err := execCmd.Run(); err != nil {
		return err
	}
	utils.PrintHeader(outputBuffer, account, command)

	outputBuffer.WriteString(out.String())
	outputBuffer.WriteString("\n")

	return nil
}

func AssumeRole(ctx context.Context, roleArn string) (*sts.AssumeRoleOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(cfg)
	params := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: aws.String("MySessionName"),
	}

	return client.AssumeRole(ctx, params)
}
