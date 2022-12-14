package stacks

import (
	"os"
	"path"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateLambda(scope constructs.Construct, config *configs.Config) awslambda.DockerImageFunction {
	current, _ := os.Getwd()
	imageProps := &awslambda.AssetImageCodeProps{
		Platform: awsecrassets.Platform_LINUX_AMD64(),
	}

	props := &awslambda.DockerImageFunctionProps{
		Code: awslambda.DockerImageCode_FromImageAsset(
			jsii.String(path.Join(current, config.Lambda.Image.File)),
			imageProps,
		),
		Environment: &map[string]*string{
			"CODEPIPELINE_NAME": jsii.String(config.Pipeline.Name),
			"BUCKET_NAME":       jsii.String(config.S3.BucketName),
			"ACCESS_TOKEN":      jsii.String(config.Github.AccessToken),
		},
	}

	id := jsii.String("id-lambda")
	lambdaFunction := awslambda.NewDockerImageFunction(scope, id, props)
	lambdaFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	return lambdaFunction
}
