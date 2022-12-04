package stacks

import (
	"fmt"
	"os"
	"path"

	assets "github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateImage(scope constructs.Construct, configuration config.Config) {
	current, _ := os.Getwd()
	imageID := fmt.Sprintf("id-image")

	props := assets.DockerImageAssetProps{
		File: jsii.String(path.Join(current, configuration.Lambda.Image.File)),
	}

	imageAssets := assets.NewDockerImageAsset(scope, jsii.String(imageID), &props)
}
