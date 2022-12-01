# サンプルアプリケーション

## ファイル・ディレクトリ構成

```bash
.
├── Dockerfile
├── README.md
├── buildspec.yml
├── deploy
│   ├── apply-image.sh
│   └── deployment.yaml
└── scripts
    └── notification.sh
```

### `Dockerfile`

`ECR` に `push` する `docker` イメージを作成するときに使用する `Dockerfile`

### `buildspeck.yml`

`Code build` 用の設定ファイル

### `deploy/apply-image.sh`

新しく作成した `docker` イメージを `Kubernetes` にデプロイするためのスクリプト

### `deploy/deployment.yml`

`Kubernetes` にデプロイするための定義ファイル

### `scripts/notification.sh`

`Slack` に通知するためのスクリプト
