# ape (aws parallel executer)

## 概要

- 複数のAWSアカウントに対してawsコマンドを実行するツール
- コマンドを実行するサーバにそれぞれのアカウントに対してSwitchRoleする権限が付与されていること
- `-y` オプションを付与することでプレビューをスキップすることができる

## 前提

- 実行サーバにawscliがインストールされていること
- コマンド実行対象のアカウントにSwitchRole権限が付与されていること

## 実行方法

```sh
# 設定関連のコマンド
aps config list
aps config add -n <名前> -i <アカウントID> -r <role名>
aps config del -i <アカウントID> -r <role名>
※ 変更は削除追加で対応

# SwitchRoleで切るか確認する
aps config check -a
aps config check -t <アカウント名>

# 全アカウントに対して実行する
aps exec -a -c <コマンド>

# アカウント指定で実行する
aps exec -t <アカウント名> -c <コマンド>
```

## その他
