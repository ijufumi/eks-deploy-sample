# EKS へ自動でデプロイするサンプルコード

このコードは、 `GitHub Action` で `Release` を作成したら
自動で `EKS(Elastic Kubernetes Service)` にデプロイします。

## 全体図

![全体図](./assets/overall.png)

## ディレクトリ構成

```bash
├── LICENSE
├── README.md
├── app
├── assets
└── deployments
```

### `app`

サンプルアプリケーションのコード。  
必要最小限のコードとデプロイするための最低限のファイルなどを格納している。

### `assets`

`REAME.md` 用の画像ファイルなどを格納している。

### `deployments`

環境を作成するための `AWS CDK` のコードを格納している。


