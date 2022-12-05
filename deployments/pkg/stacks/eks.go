package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/config"
)

func CreateEKS(scope constructs.Construct, configuration config.Config) {
	eksTaskRole := awsiam.NewRole(scope, jsii.String("eks-task-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	id := fmt.Sprintf("id-%s", configuration.Cluster.Name)
	props := awseks.ClusterProps{
		Version:     awseks.KubernetesVersion_Of(jsii.String(configuration.Cluster.K8SVersion)),
		ClusterName: jsii.String(configuration.Cluster.Name),
		Role:        eksTaskRole,
	}

	cluster := awseks.NewCluster(scope, jsii.String(id), &props)
}
