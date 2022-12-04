package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateVPC(scope constructs.Construct, configuration config.Config) awsec2.Vpc {
	vpcID := fmt.Sprintf("id-%s", configuration.Vpc.Name)
	subnetConfigurations := []*awsec2.SubnetConfiguration{
		{
			Name:       jsii.String(fmt.Sprintf("public-subnet-%s", configuration.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		{
			Name:       jsii.String(fmt.Sprintf("private1-subnet-%s", configuration.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_NAT,
		},
		{
			Name:       jsii.String(fmt.Sprintf("private2-subnet-%s", configuration.Vpc.Name)),
			SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
		},
	}
	props := awsec2.VpcProps{
		VpcName:             jsii.String(configuration.Vpc.Name),
		Cidr:                jsii.String(configuration.Vpc.CidrBlock),
		MaxAzs:              jsii.Number(2),
		SubnetConfiguration: &subnetConfigurations,
	}

	return awsec2.NewVpc(scope, jsii.String(vpcID), &props)
}
