package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateVPC(scope constructs.Construct, config *configs.Config) awsec2.Vpc {
	subnetConfigurations := []*awsec2.SubnetConfiguration{
		{
			Name:       jsii.String(fmt.Sprintf("public-subnet-%s", config.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		{
			Name:       jsii.String(fmt.Sprintf("private1-subnet-%s", config.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
		{
			Name:       jsii.String(fmt.Sprintf("private2-subnet-%s", config.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
		},
	}
	props := awsec2.VpcProps{
		VpcName:             jsii.String(config.Vpc.Name),
		IpAddresses:         awsec2.IpAddresses_Cidr(jsii.String(config.Vpc.CidrBlock)),
		MaxAzs:              jsii.Number(2),
		SubnetConfiguration: &subnetConfigurations,
	}

	return awsec2.NewVpc(scope, jsii.String("vpc-id"), &props)
}
