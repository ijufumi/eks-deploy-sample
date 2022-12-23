package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	kubectl "github.com/aws/aws-cdk-go/awscdk/v2/lambdalayerkubectl"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateEKS(scope constructs.Construct, config *configs.Config, vpc awsec2.Vpc) awseks.Cluster {
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
		Version:      awseks.KubernetesVersion_Of(jsii.String(config.Cluster.K8SVersion)),
		KubectlLayer: kubectl.NewKubectlLayer(scope, jsii.String("id-kubectl-layer")),
		ClusterName:  jsii.String(config.Cluster.Name),
		// Role:        eksTaskRole,
		Vpc:        vpc,
		VpcSubnets: &subnets,
		AlbController: &awseks.AlbControllerOptions{
			Version: awseks.AlbControllerVersion_V2_4_1(),
		},
	}

	cluster := awseks.NewFargateCluster(scope, jsii.String(id), &props)

	appName := jsii.String(config.Cluster.App.Name)

	cluster.AddManifest(jsii.String("id-app-manifest"), &map[string]interface{}{
		"apiVersion": jsii.String("apps/v1"),
		"kind":       jsii.String("Deployment"),
		"metadata": map[string]*string{
			"name": appName,
		},
		"spec": map[string]interface{}{
			"replicas": jsii.Number(1),
			"selector": map[string]interface{}{
				"matchLabels": map[string]*string{
					"app": appName,
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]*string{
						"app": appName,
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":            appName,
							"image":           jsii.String(config.Cluster.App.Image),
							"imagePullPolicy": jsii.String("IfNotPresent"),
							"resources": map[string]interface{}{
								"limits": map[string]*string{
									"cpu":    jsii.String("128m"),
									"memory": jsii.String("128Mi"),
								},
								"requests": map[string]*string{
									"cpu":    jsii.String("128m"),
									"memory": jsii.String("128Mi"),
								},
							},
							"ports": []map[string]interface{}{
								{
									"name":          jsii.String("task"),
									"containerPort": jsii.Number(80),
								},
							},
							"readinessProbe": map[string]interface{}{
								"tcpSocket": map[string]interface{}{
									"port": jsii.Number(80),
								},
								"initialDelaySeconds": jsii.Number(15),
								"timeoutSeconds":      jsii.Number(2),
							},
							"livenessProbe": map[string]interface{}{
								"tcpSocket": map[string]interface{}{
									"port": jsii.Number(80),
								},
								"initialDelaySeconds": jsii.Number(45),
								"timeoutSeconds":      jsii.Number(2),
							},
						},
					},
				},
			},
		},
	})

	cluster.AddManifest(jsii.String("id-service-manifest"), &map[string]interface{}{
		"apiVersion": jsii.String("v1"),
		"kind":       jsii.String("Service"),
		"metadata": map[string]interface{}{
			"name": jsii.String("service"),
			"annotations": map[string]*string{
				"service.beta.kubernetes.io/aws-load-balancer-type":            jsii.String("external"),
				"service.beta.kubernetes.io/aws-load-balancer-nlb-target-type": jsii.String("ip"),
			},
		},
		"spec": map[string]interface{}{
			"type": jsii.String("LoadBalancer"),
			"selector": map[string]*string{
				"app": appName,
			},
			"ports": []map[string]interface{}{
				{
					"protocol":   jsii.String("TCP"),
					"port":       jsii.Number(80),
					"targetPort": jsii.Number(80),
				},
			},
		},
	})

	return cluster
}
