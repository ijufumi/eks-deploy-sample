package stacks

import (
	"fmt"

	build "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipeline"
	actions "github.com/aws/aws-cdk-go/awscdk/v2/awscodepipelineactions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateCodepipeline(scope constructs.Construct, config *configs.Config, bucket awss3.IBucket, repository awsecr.IRepository, cluster awseks.ICluster, eksMasterRole awsiam.IRole) pipeline.Pipeline {
	buildRole := awsiam.NewRole(scope, jsii.String("codebuild-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codebuild.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	role := awsiam.NewRole(scope, jsii.String("codepipeline-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("codepipeline.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:Get*", "s3:List*"),
		Resources: jsii.Strings(*bucket.BucketArn(), *jsii.String(fmt.Sprintf("%s/*", *bucket.BucketArn()))),
	}))
	role.AddToPolicy(awsiam.NewPolicyStatement(
		&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings("sts:AssumeRole"),
			Effect:    awsiam.Effect_ALLOW,
			Resources: jsii.Strings("*"),
		},
	))

	sourceOutput := pipeline.NewArtifact(jsii.String("source"))
	sourceAction := actions.NewS3SourceAction(
		&actions.S3SourceActionProps{
			ActionName: jsii.String("Source"),
			Bucket:     bucket,
			BucketKey:  jsii.String("sample.zip"),
			Output:     sourceOutput,
			Role:       role,
			Trigger:    actions.S3Trigger_NONE,
		},
	)

	buildProject := build.NewPipelineProject(scope, jsii.String("id-codebuild"), &build.PipelineProjectProps{
		Role: buildRole,
		Environment: &build.BuildEnvironment{
			Privileged: jsii.Bool(true),
			BuildImage: build.LinuxBuildImage_AMAZON_LINUX_2_4(),
		},
	})

	buildRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ecr:PutImage"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings(*repository.RepositoryArn(), fmt.Sprintf("%s/*", *repository.RepositoryArn())),
	}))
	buildRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ecr:GetAuthorizationToken"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings("*"),
	}))
	buildRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("eks:DescribeCluster"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings(*cluster.ClusterArn()),
	}))
	buildRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("sts:AssumeRole"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings(*eksMasterRole.RoleArn()),
	}))

	dockerUser := awsssm.NewStringParameter(scope, jsii.String("id-docker-user"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("dodker-user"),
		StringValue:   jsii.String(config.Docker.User),
	})
	dockerToken := awsssm.NewStringParameter(scope, jsii.String("id-docker-token"), &awsssm.StringParameterProps{
		ParameterName: jsii.String("dodker-token"),
		StringValue:   jsii.String(config.Docker.Token),
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
					Value: config.GetAwsAccountID(),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
				"EKS_CLUSTER_NAME": {
					Value: jsii.String(config.Cluster.Name),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
				"EKS_CLUSTER_ROLE": {
					Value: eksMasterRole.RoleArn(),
					Type:  build.BuildEnvironmentVariableType_PLAINTEXT,
				},
				"DOCKER_USER": {
					Value: dockerUser.ParameterName(),
					Type:  build.BuildEnvironmentVariableType_PARAMETER_STORE,
				},
				"DOCKER_TOKEN": {
					Value: dockerToken.ParameterName(),
					Type:  build.BuildEnvironmentVariableType_PARAMETER_STORE,
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

	props := &pipeline.PipelineProps{
		PipelineName: jsii.String(config.Pipeline.Name),
		Stages:       &stages,
		Role:         role,
	}

	return pipeline.NewPipeline(scope, jsii.String("id-codepipeline"), props)
}
