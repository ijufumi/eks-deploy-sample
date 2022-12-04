package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateEKS(scope constructs.Construct, configuration config.Config) {
	id := fmt.Sprintf("id-%s", configuration.Cluster.Name)
	props := awseks.ClusterProps{
		Version:     awseks.KubernetesVersion_Of(jsii.String(configuration.Cluster.K8SVersion)),
		ClusterName: jsii.String(configuration.Cluster.Name),
	}

	cluster := awseks.NewCluster(scope, jsii.String(id), &props)
}
