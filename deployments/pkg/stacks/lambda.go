package stacks

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	codepipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateLambda(scope constructs.Construct, config *configs.Config, s3 awss3.IBucket, pipeline codepipeline.Pipeline) awslambda.DockerImageFunction {
	current, _ := os.Getwd()
	imageProps := &awslambda.AssetImageCodeProps{
		Platform: awsecrassets.Platform_LINUX_AMD64(),
	}

	role := awsiam.NewRole(scope, jsii.String("lambda-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObject"),
		Resources: jsii.Strings(*s3.BucketArn(), *jsii.String(fmt.Sprintf("%s/*", *s3.BucketArn()))),
	}))
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("codepipeline:UpdatePipeline", "codepipeline:ListPipelines", "codepipeline:GetPipeline", "codepipeline:StartPipelineExecution"),
		Resources: jsii.Strings(*pipeline.PipelineArn(), *jsii.String(fmt.Sprintf("%s/*", *pipeline.PipelineArn()))),
	}))
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("iam:PassRole"),
		Resources: jsii.Strings("*"),
		Conditions: &map[string]interface{}{
			"StringEquals": map[string]interface{}{
				"iam:PassedToService": []string{
					"codepipeline.amazonaws.com",
				},
			},
		},
	}))

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
		Timeout: awscdk.Duration_Hours(jsii.Number(config.Lambda.TimeoutHours)),
		Role:    role,
	}

	id := jsii.String("id-lambda")
	lambdaFunction := awslambda.NewDockerImageFunction(scope, id, props)
	lambdaFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	return lambdaFunction
}
