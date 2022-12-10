# 環境作成用のコード

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

## 動かすための前提条件

以下のコマンドがインストールされていること

### `awscli`

Mac で `brew` を使っている場合、以下のコマンドでインストール可能。

```bash
brew install awscli
```

### `awscdk`

Mac で `brew` を使っている場合、以下のコマンドでインストール可能。

```bash
brew install aws-cdk
```

`CDK` の使い方については、 [CDK README](./README.cdk.md) を参照

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
