package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"

	//	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateEKS(scope constructs.Construct, configuration config.Config, vpc awsec2.Vpc) awseks.Cluster {
	// eksTaskRole := awsiam.NewRole(scope, jsii.String("eks-task-role"), &awsiam.RoleProps{
	// 	AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	// })

	subnets := []*awsec2.SubnetSelection{
		{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_NAT,
		},
	}
	id := "eks-cluster-id"
	props := awseks.FargateClusterProps{
		Version:     awseks.KubernetesVersion_Of(jsii.String(configuration.Cluster.K8SVersion)),
		ClusterName: jsii.String(configuration.Cluster.Name),
		// Role:        eksTaskRole,
		Vpc:        vpc,
		VpcSubnets: &subnets,
	}

	cluster := awseks.NewFargateCluster(scope, jsii.String(id), &props)

	return cluster
}
