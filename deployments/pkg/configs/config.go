package configs

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config is application configuration
type Config struct {
	AwsAccountID      string `env:"AWS_ACCOUNT_ID"`
	AwsRegion         string `env:"AWS_REGION"`
	CdkDefaultAccount string `env:"CDK_DEFAULT_ACCOUNT"`
	CdkDefaultRegion  string `env:"CDK_DEFAULT_REGION"`

	Vpc struct {
		Name      string `env:"VPC_NAME"`
		CidrBlock string `env:"VPC_CIDR_BLOCK"`
	}

	Repository struct {
		Name string `env:"REPOSITORY_NAME"`
	}

	Lambda struct {
		Image struct {
			File string `env:"LAMBDA_IMAGE_FILE" envDefault:"../lambda"`
		}
		TimeoutSec float64 `env:"LAMBDA_TIMEOUT_SEC" envDefault:"180"`
	}

	Cluster struct {
		Name       string `env:"CLUSTER_NAME"`
		K8SVersion string `env:"CLUSTER_K8S_VERSION" envDefault:"1.24"`
		App        struct {
			Name  string `env:"CLUSTER_APP_NAME"`
			Image string `env:"CLUSTER_APP_IMAGE"`
		}
		AdminUsers []string `env:"CLUSTER_ADMIN_USERS" envSeparator:","`
		AdminRoles []string `env:"CLUSTER_ADMIN_ROLES" envSeparator:","`
	}

	Pipeline struct {
		Name string `env:"PIPELINE_NAME" envDefault:"pipeline"`
	}

	S3 struct {
		BucketName string `env:"S3_BUCKET_NAME" envDefault:"pipeline_bucket"`
	}

	Github struct {
		AccessToken string `env:"GITHUB_ACCESS_TOKEN"`
	}

	Slack struct {
		WebHookURL string `env:"WEB_HOOK_URL"`
	}

	Docker struct {
		User  string `env:"DOCKER_USER"`
		Token string `env:"DOCKER_TOKEN"`
	}
}

func (c *Config) GetAwsAccountID() *string {
	if len(c.AwsAccountID) != 0 {
		return jsii.String(c.AwsAccountID)
	}
	return jsii.String(c.CdkDefaultAccount)
}

func (c *Config) GetAwsRegion() *string {
	if len(c.AwsRegion) != 0 {
		return jsii.String(c.AwsRegion)
	}
	return jsii.String(c.CdkDefaultRegion)
}

// Load returns configuration made from environment variables
func Load() *Config {
	_ = godotenv.Load()
	c := Config{}
	_ = env.Parse(&c)

	return &c
}
