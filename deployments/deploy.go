package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/stacks"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DeployStackProps struct {
	awscdk.StackProps
}

func NewDeployStack(scope constructs.Construct, id string, props *DeployStackProps, config *configs.Config) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create stacks
	vpc := stacks.CreateVPC(stack, config)
	s3 := stacks.CreateS3(stack, config)
	_ = stacks.CreateECR(stack, config)
	cluster := stacks.CreateEKS(stack, config, vpc)
	_ = stacks.CreateCodepipeline(stack, config, s3)
	lambda := stacks.CreateLambda(stack, config)

	// Output results
	awscdk.NewCfnOutput(stack, jsii.String("id-output-labmda"), &awscdk.CfnOutputProps{
		Value:      lambda.FunctionName(),
		ExportName: jsii.String("labmda-function-url"),
	})

	fmt.Printf("labmda-function-url is %s", *lambda.FunctionName())
	fmt.Printf("alb-endpoint-url is %s", *cluster.AlbController().ToString())

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	config := configs.Load()
	NewDeployStack(app, "DeployStack", &DeployStackProps{
		awscdk.StackProps{
			Env: env(config),
		},
	}, config)

	app.Synth(nil)
}

func env(config *configs.Config) *awscdk.Environment {
	var account = config.AwsAccessKeyID
	if len(account) == 0 {
		account = config.CdkDefaultAccount
	}
	var region = config.AwsRegion
	if len(region) == 0 {
		region = config.CdkDefaultRegion
	}
	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}

}
