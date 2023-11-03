# ape (awscli parallel executer)

## 概要

- 複数のAWSアカウントに対してawsコマンドを実行するツール
- 複数アカウントに対してコマンドを並列で実行するため、高速に処理を行うことができる
- jq等のコマンドと組み合わせた処理が可能
- `-y` オプションを付与することでプレビューをスキップすることができる
  
## 前提

- コマンド実行サーバにawscliがインストールされていること
- コマンド実行対象のアカウントにSwitchRole権限が付与されていること

## インストール方法

```sh
arch="darwin-arm64"

curl -L https://github.com/yasuyuki0321/ape/releases/latest/download/ape-${arch}.tar.gz | tar zxvf -
chmod 755 ape-${arch}

※ 必要に応じてリンクを作成したり、/bin等、PATHの通っているディレクトリに配置する
ln -s ./ape-${arch} ./ape
mv ./ape-${arch} /bin/
```

## 設定ファイル

- ファイル名: `.config.yaml`
- apeバイナリと同じディレクトリに配置する
- ape config add/removeで `.config.yaml` ファイルの操作が可能

```yaml
accounts:
  - name: "account001"
    roleArn: "arn:aws:iam::account001:role/assume-role-test"
  - name: "account002"
    roleArn: "arn:aws:iam::account002:role/assume-role-test"
```

## 実行方法

### config

```sh
manage .config.yaml file

Usage:
  ape config [command]

Available Commands:
  add         Add an account
  list        List all accounts
  remove      Remove an account
  test        Test connection to AWS accounts
```

### exec

```sh
Execute AWS command across multiple accounts

Usage:
  ape exec [flags]

Flags:
  -c, --command string   Command to execute
  -h, --help             help for exec
  -y, --skip-preview     Skip the preview and execute the command directly
  -t, --target string    Command target accounts (default "all")
```

## コマンドの実行例

### config

.config.yamlの登録されているアカウントの確認

```sh
./ape config list
Accounts:
Name: test001, RoleArn: arn:aws:iam::123456789123:role/assume-role-test
```

アカウントの登録

```sh
./ape config add -n test002 -r arn:aws:iam::987654321987:role/assume-role-test
Account: test002(987654321987) - OK
```

アカウントの削除

```sh
./ape config remove -n test002
```

接続テスト

```sh
./ape config test
Account: test001(123456789123) - OK
Account: test002(987654321987) - OK
```

### exec

登録されているすべてのアカウントに対してコマンドを実行

```sh
./ape exec -c "aws s3 ls"
Accounts:
Account: test001(123456789123)
Account: test002(987654321987)

Command: aws s3 ls

Do you want to continue? [y/N]: y
----------
Time: 2023-11-03 12:51:52
Account: test001(123456789123)
Command: aws s3 ls
----------
2022-04-27 20:32:11 sts-test-001

----------
Time: 2023-11-03 12:51:52
Account: test001(123456789123)
Command: aws s3 ls
----------
2022-04-27 20:32:11 sts-test-002
```

特定のアカウントに対してコマンドを実行

```sh
./ape exec -t test001 -c "aws s3 ls" -y
Accounts:
Account: test001(123456789123)

Command: aws s3 ls

Do you want to continue? [y/N]: y
----------
Time: 2023-11-03 12:56:11
Account: test001(123456789123)
Command: aws s3 ls
----------
2022-04-27 20:32:11 sts-test-001
```

アカウントの指定にはワイルドカードの使用が可能

```sh
./ape exec -t 'test*' -c "aws s3 ls" -y
----------
Time: 2023-11-03 12:59:00
Account: test001(123456789123)
Command: aws s3 ls
----------
2022-04-27 20:32:11 sts-test-001

----------
Time: 2023-11-03 12:59:00
Account: test002(987654321987)
Command: aws s3 ls
----------
2022-04-27 20:32:11 sts-test-002
```

jq等と組み合わせて使用することも可能

```sh
./ape -t test001 -c "aws ec2 describe-instances | jq .Reservations[].Instances[].InstanceId" -y
----------
Time: 2023-11-03 13:02:16
Account: test001(123456789123)
Command: aws ec2 describe-instances | jq .Reservations[].Instances[].InstanceId
----------
"i-09ef99b3c57018602"
"i-0a9ad44aa54f06a79"
```
