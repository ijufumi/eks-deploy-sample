package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/ijufumi/eks-deploy-sample/deployments/pkg/configs"
)

func CreateS3(scope constructs.Construct, config *configs.Config) awss3.IBucket {
	return awss3.NewBucket(scope, jsii.String("id-s3-bucket"), &awss3.BucketProps{
		BucketName:      jsii.String(config.S3.BucketName),
		Versioned:       jsii.Bool(true),
		ObjectOwnership: awss3.ObjectOwnership_BUCKET_OWNER_ENFORCED,
	})
}
