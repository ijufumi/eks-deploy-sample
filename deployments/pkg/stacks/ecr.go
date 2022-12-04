package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateECR(scope constructs.Construct, configuration config.Config) awsecr.Repository {
	repositoryID := fmt.Sprintf("id-%s", configuration.Repository.Name)
	reporitoyProps := awsecr.RepositoryProps{
		RepositoryName:  jsii.String(configuration.Repository.Name),
		ImageScanOnPush: jsii.Bool(true),
	}

	repository := awsecr.NewRepository(scope, jsii.String(repositoryID), &reporitoyProps)

	return repository
}
