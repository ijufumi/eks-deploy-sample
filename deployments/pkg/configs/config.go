package configs

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config is application configuration
type Config struct {
	AwsAccessKeyID    string `env:"AWS_ACCESS_KEY_ID"`
	AwsRegion         string `env:"AWS_REGION"`
	CdkDefaultAccount string `env:"CDK_DEFAULT_ACCOUNT"`
	CdkDefaultRegion  string `env:"CDK_DEFAULT_REGION"`

	Vpc struct {
		Name      string `env:"VPC_NAME"`
		CidrBlock string `env:"CIDR_BLOCK"`
	}

	Repository struct {
		Name string `env:"REPOSITORY_NAME"`
	}

	Lambda struct {
		Image struct {
			File string `env:"LAMBDA_IMAGE_FILE" default:"../lambda"`
			Tag  string `env:"LAMBDA_IMAGE_TAG" default:"latest"`
		}
		Repository struct {
			Name string `env:"LAMBDA_REPOSITORY_NAME" default:"lambda"`
		}
	}

	Cluster struct {
		Name       string `env:"CLUSTER_NAME"`
		K8SVersion string `env:"CLUSTER_K8S_VERSION" default:"1.21"`
		App        struct {
			Name  string `env:"CLUSTER_APP_NAME"`
			Image string `env:"CLUSTER_APP_IMAGE"`
		}
	}

	Pipeline struct {
		Name string `env:"PIPELINE_NAME" default:"pipeline"`
	}

	S3 struct {
		BucketName string `env:"S3_BUCKET_NAME" default:"pipeline_bucket"`
	}

	Github struct {
		OrganizationName string `env:"GITHUB_ORGANIZATION_NAME"`
		RepositoryName   string `env:"GITHUB_REPOSITORY_NAME"`
		AccessToken      string `env:"GITHUB_ACCESS_TOKEN"`
	}

	Slack struct {
		WebHookURL string `env:"WEB_HOOK_URL"`
	}
}

// Load returns configuration made from environment variables
func Load() *Config {
	_ = godotenv.Load()
	c := Config{}
	_ = env.Parse(&c)

	return &c
}
