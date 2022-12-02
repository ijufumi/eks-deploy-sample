# 環境作成用のコード

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
