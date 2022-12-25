# 環境作成用のコード

## 注意

2022/12/23現在、一部ちゃんと動かない可能性があります。

## ファイル・ディレクトリ構成

```bash
.
├── deploy.go
├── deploy_test.go
├── lambda
└── pkg
    ├── config
    └── stacks
```

### `deploy.go`

`CDK` の本体

### `deploy_test.go`

`CDK` のテストファイル

### `labmda`

`Lambda` 用のコード

### `pkg/config`

`CDK` 用の設定

### `pkg/stacks`

各リソースの `Stack` ファイル

## セットアップ

### 必要なコマンド

以下のコマンドがインストールされていること

#### `awscli`

Mac で `brew` を使っている場合、以下のコマンドでインストール可能。

```bash
brew install awscli
```

#### `awscdk`

Mac で `brew` を使っている場合、以下のコマンドでインストール可能。

```bash
brew install aws-cdk
```

`CDK` の使い方については、 [CDK README](./README.cdk.md) を参照

### `.env` ファイル

`.env.example` をコピーして `.env` を作成し、以下の内容を記述する

| 項目 | 説明 |
| ---- | ---- |
| `AWS_ACCOUNT_ID` | `ECR` の名前解決に使う `AWS_ACCOUNT_ID` |
| `VPC_CIDR_BLOCK` | `VPC` に設定する `CIDR` |
| `VPC_NAME` | `VPC` に設定する名前 |
| `REPOSITORY_NAME` | `ECR` の名前 |
| `LAMBDA_IMAGE_FILE` | `Lambda` 用の `docker` イメージ |
| `LAMBDA_TIMEOUT_SEC` | `Lambda` に設定するタイムアウト（秒）最大は `900` |
| `CLUSTER_NAME` | `EKS` の名前 |
| `CLUSTER_K8S_VERSION` | `EKS` に設定する `Kubernetes` のバージョン |
| `CLUSTER_APP_NAME` | `EKS` にデプロイする `app` の名前 |
| `CLUSTER_APP_IMAGE` | `EkS` にデプロイする最初の `docker` イメージ |
| `CLUSTER_ADMIN_USERS` | `EKS` の `ConfigMap` に設定するユーザの `ARN` リスト。 `,` で区切って指定する |
| `PIPELINE_NAME` | `Codepipeline` の名前 |
| `S3_BUCKET_NAME` | `S3` バケットの名前 |
| `GITHUB_ACCESS_TOKEN` | `Lambda` から `Github` にアクセスするためのトークン |
| `WEB_HOOK_URL` | `Slack` に通知するための `URL`. 未設定の場合は通知しない |
| `DOCKER_USER` | `Code build` で `docker login` する際のユーザ |
| `DOCKER_TOKEN` | `Code build` で `docker login` する際のユーザのパスワード |

### `AWS`アカウント

もし、`aws configure` で `profile` が `default` じゃない場合、以下のコマンドを実行する。

```bash
aws configure --profile [YOUR PROFILE NAME]
```

### public aws ecr へのログイン

```bash
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
```

## このコードで作成する `Stack` たち

### `VPC`

`EKS` を作成する `VPC` を作成するための `Stack`

### `Lambda`

`GitHub Action` から実行される `Lambda` を作成するための `Stack`

### `Code pipeline` + `Code build`

`docker` イメージの `build` 〜 `EKS` へのデプロイを実行するための `Code pipeline` を作成するための `Stack` 

### `EKS`

`EKS` を作成するための `Stack`

### `ECR`

`EKS` にデプロイする `docker` イメージを格納する `ECR` を作成するための `Stack`

### `S3`

`Codepipline` で `Source` のファイルを格納する `bucket` を作成する `Stack`
