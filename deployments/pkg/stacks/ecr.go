package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateECR(scope constructs.Construct, config *configs.Config) awsecr.Repository {
	repositoryID := "ecr-repository-id"
	reporitoyProps := awsecr.RepositoryProps{
		RepositoryName:  jsii.String(config.Repository.Name),
		ImageScanOnPush: jsii.Bool(true),
	}

	return awsecr.NewRepository(scope, jsii.String(repositoryID), &reporitoyProps)
}
