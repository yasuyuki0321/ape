package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Account struct {
	Name    string
	RoleArn string
}

func GetHomePath(path string) string {
	if len(path) < 2 || path[:2] != "~/" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, path[2:])
}

func PrintHeader(outputBuffer *bytes.Buffer, account Account, command string) {
	accountId := strings.Split(account.RoleArn, ":")[4]

	outputBuffer.WriteString(fmt.Sprintln(strings.Repeat("-", 10)))
	outputBuffer.WriteString(fmt.Sprintf("Time: %v\n", time.Now().Format("2006-01-02 15:04:05")))
	outputBuffer.WriteString(fmt.Sprintf("Account: %s (%s)\n", account.Name, accountId))
	outputBuffer.WriteString(fmt.Sprintf("Command: %s\n", command))
	outputBuffer.WriteString(fmt.Sprintln(strings.Repeat("-", 10)))
}

func DisplayPreview(accounts []Account, command string) bool {
	fmt.Println("Accounts:")

	for _, account := range accounts {
		fmt.Printf("Name: %s command: %s\n", account.Name, command)
	}

	fmt.Print("\nDo you want to continue? [y/N]: ")
	var response string
	fmt.Scan(&response)

	return strings.ToLower(response) == "y"
}

func LoadAccountsFromConfig() ([]Account, error) {
	v := viper.New()
	v.SetConfigName(".config") // 設定ファイルの名前（拡張子なし）
	v.AddConfigPath(".")       // 設定ファイルのディレクトリ
	v.SetConfigType("yaml")    // 使用するファイル形式

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var accounts []Account
	if err := v.UnmarshalKey("accounts", &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func FilterAccounts(roleArns []Account, target string) []Account {
	if target == "all" {
		return roleArns
	}

	var filtered []Account
	for _, roleArn := range roleArns {
		if matchesTarget(roleArn.Name, target) {
			filtered = append(filtered, roleArn)
		}
	}
	return filtered
}

func matchesTarget(accountName, target string) bool {
	targets := strings.Split(target, ",")

	for _, t := range targets {
		if matched, _ := filepath.Match(t, accountName); matched {
			return true
		}
	}
	return false
}
