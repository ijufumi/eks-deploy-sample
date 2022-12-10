package stacks

import (
	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateCodepipeline(scope constructs.Construct, config *configs.Config) pipeline.Pipeline {
	sourceOutput := pipeline.NewArtifact(jsii.String("source"))
	sourceAction := actions.NewS3SourceAction(
		&actions.S3SourceActionProps{
			ActionName: jsii.String("Source"),
			Bucket:     awss3.Bucket_FromBucketName(scope, jsii.String("id-s3-bucket"), jsii.String(config.S3.BucketName)),
			BucketKey:  jsii.String("sample.zip"),
			Output:     sourceOutput,
		},
	)

	buildProject := build.NewPipelineProject(scope, jsii.String("id-codebuild"), &build.PipelineProjectProps{})

	buildAction := actions.NewCodeBuildAction(
		&actions.CodeBuildActionProps{
			ActionName: jsii.String("build"),
			Input:      sourceOutput,
			Project:    buildProject,
		},
	)

	stages := []*pipeline.StageProps{
		{
			StageName: jsii.String("Source"),
			Actions: &[]pipeline.IAction{
				sourceAction,
				buildAction,
			},
		},
	}

	props := &pipeline.PipelineProps{
		PipelineName: jsii.String(config.Pipeline.Name),
		Stages:       &stages,
	}

	return pipeline.NewPipeline(scope, jsii.String("id-codepipeline"), props)
}
