package stacks

import (
	"fmt"

	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateCodepipeline(scope constructs.Construct, config *configs.Config, bucket awss3.IBucket) pipeline.Pipeline {
	sourceRole := awsiam.NewRole(scope, jsii.String("codepipeline-source-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codepipeline.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	sourceRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:Get*", "s3:List*"),
		Resources: jsii.Strings(*bucket.BucketArn(), *jsii.String(fmt.Sprintf("%s/*", *bucket.BucketArn()))),
	}))

	sourceOutput := pipeline.NewArtifact(jsii.String("source"))
	sourceAction := actions.NewS3SourceAction(
		&actions.S3SourceActionProps{
			ActionName: jsii.String("Source"),
			Bucket:     bucket,
			BucketKey:  jsii.String("sample.zip"),
			Output:     sourceOutput,
			Role:       sourceRole,
		},
	)

	buildRole := awsiam.NewRole(scope, jsii.String("codepipeline-source-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codepipeline.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	buildRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ecr:*"),
		Resources: jsii.Strings("*"),
	}))
	buildProject := build.NewPipelineProject(scope, jsii.String("id-codebuild"), &build.PipelineProjectProps{
		Role: buildRole,
	})

	buildAction := actions.NewCodeBuildAction(
		&actions.CodeBuildActionProps{
			ActionName: jsii.String("build"),
			Input:      sourceOutput,
			Project:    buildProject,
			EnvironmentVariables: &map[string]*build.BuildEnvironmentVariable{
				"IMAGE_REPO_NAME": {
					Value: jsii.String(config.Repository.Name),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
				"WEB_HOOK_URL": {
					Value: jsii.String(config.Slack.WebHookURL),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
				"AWS_ACCOUNT_ID": {
					Value: jsii.String(config.AwsAccessKeyID),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
			},
		},
	)

	stages := []*pipeline.StageProps{
		{
			StageName: jsii.String("Source"),
			Actions: &[]pipeline.IAction{
				sourceAction,
			},
		},
		{
			StageName: jsii.String("Build"),
			Actions: &[]pipeline.IAction{
				buildAction,
			},
		},
	}

	role := awsiam.NewRole(scope, jsii.String("codepipeline-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codepipeline.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:Get*", "s3:List*"),
		Resources: jsii.Strings(*bucket.BucketArn(), *jsii.String(fmt.Sprintf("%s/*", *bucket.BucketArn()))),
	}))

	props := &pipeline.PipelineProps{
		PipelineName: jsii.String(config.Pipeline.Name),
		Stages:       &stages,
		Role:         role,
	}

	return pipeline.NewPipeline(scope, jsii.String("id-codepipeline"), props)
}
