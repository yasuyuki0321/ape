package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/yasuyuki0321/ape/pkg/aws"
	"github.com/yasuyuki0321/ape/pkg/utils"
	"gopkg.in/yaml.v3"
)

const roleArnPattern = `^arn:aws:iam::\d{12}:role/.+$`

type Account struct {
	Name    string `yaml:"name"`
	RoleArn string `yaml:"roleArn"`
}

type Config struct {
	Accounts []Account `yaml:"accounts"`
}

func readConfig(configFile string) (*Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func writeConfig(configFile string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, os.ModePerm)
}

func isValidRoleArn(roleArn string) bool {
	matched, _ := regexp.MatchString(roleArnPattern, roleArn)
	return matched
}

func AddAccount(configFile string, name, roleArn string) error {

	// configファイルがなければ作成
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		_, err = os.Create(configFile)
		if err != nil {
			return err
		}
	}

	cfg, err := readConfig(configFile)
	if err != nil {
		return err
	}

	// アカウント名が既に存在するかを確認
	for _, acc := range cfg.Accounts {
		if acc.Name == name {
			return fmt.Errorf("account name '%s' already exists", name)
		}
	}

	// 接続テストを行う
	testAccount := utils.Account{Name: name, RoleArn: roleArn}
	if err := CheckAWSConnections([]utils.Account{testAccount}); err != nil {
		return fmt.Errorf("failed to test connection for account %s with roleArn %s: %w", name, roleArn, err)
	}

	// roleの形式が正しいか確認
	if !isValidRoleArn(roleArn) {
		return fmt.Errorf("provided roleArn is not valid: %s", roleArn)
	}

	newAccount := Account{
		Name:    name,
		RoleArn: roleArn,
	}

	// 新しいアカウントの追加
	cfg.Accounts = append(cfg.Accounts, newAccount)
	return writeConfig(configFile, cfg)
}

func RemoveAccount(configFile string, name string) error {
	cfg, err := readConfig(configFile)
	if err != nil {
		return err
	}

	exists := false
	newAccounts := []Account{}

	for _, acc := range cfg.Accounts {
		if acc.Name == name {
			exists = true
			continue
		}
		newAccounts = append(newAccounts, acc)
	}

	if !exists {
		return fmt.Errorf("account with name '%s' does not exist", name)
	}

	cfg.Accounts = newAccounts
	return writeConfig(configFile, cfg)
}

func ListAccounts(configFile string) {
	cfg, err := readConfig(configFile)
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println("Accounts:")
	for _, acc := range cfg.Accounts {
		fmt.Printf("Name: %s, RoleArn: %s\n", acc.Name, acc.RoleArn)
	}
}

func CheckAWSConnections(roleArns []utils.Account) error {
	ctx := context.TODO()

	for _, roleArn := range roleArns {
		var outputBuffer bytes.Buffer

		creds, err := aws.AssumeRole(ctx, roleArn.RoleArn)
		if err != nil {
			fmt.Printf("Account: %s - NG (Failed to assume role: %s)\n", roleArn.Name, err)
			continue
		}

		aws.SetTempCredentials(creds)
		cmd := "aws sts get-caller-identity"
		err = aws.ExecuteAWSCLI(&outputBuffer, roleArn, cmd)

		if err != nil {
			fmt.Printf("Account: %s(%s) - NG (Failed to execute AWS CLI: %s)\n", roleArn.Name, utils.GetAccountID(roleArn.RoleArn), err)
		} else {
			fmt.Printf("Account: %s(%s) - OK\n", roleArn.Name, utils.GetAccountID(roleArn.RoleArn))
		}

		aws.ResetCredentials()
	}
	return nil
}
