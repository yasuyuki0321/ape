package config

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/yasuyuki0321/ape/pkg/aws"
	"github.com/yasuyuki0321/ape/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Account struct {
	Name    string `yaml:"name"`
	RoleArn string `yaml:"roleArn"`
}

type Config struct {
	Accounts []Account `yaml:"accounts"`
}

func readConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func writeConfig(configPath string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, os.ModePerm)
}

func AddAccount(configPath string, name, roleArn string) error {
	cfg, err := readConfig(configPath)
	if err != nil {
		return err
	}

	newAccount := Account{
		Name:    name,
		RoleArn: roleArn,
	}

	cfg.Accounts = append(cfg.Accounts, newAccount)
	return writeConfig(configPath, cfg)
}

func RemoveAccount(configPath string, name string) error {
	cfg, err := readConfig(configPath)
	if err != nil {
		return err
	}

	newAccounts := []Account{}
	for _, acc := range cfg.Accounts {
		if acc.Name != name {
			newAccounts = append(newAccounts, acc)
		}
	}

	cfg.Accounts = newAccounts
	return writeConfig(configPath, cfg)
}

func ListAccounts(configPath string) {
	cfg, err := readConfig(configPath)
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
			fmt.Printf("Account: %s - NG (Failed to execute AWS CLI: %s)\n", roleArn.Name, err)
		} else {
			fmt.Printf("Account: %s - OK\n", roleArn.Name)
		}

		aws.ResetCredentials()
	}
	return nil
}
