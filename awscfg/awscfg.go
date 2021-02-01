package awscfg

import (
	"github.com/AlekSi/pointer"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/NYTimes/gizmo/config/aws"
)

// Service is a dependency that will provide a gizmo aws.Config struct
var Service = dependency.Service{
	Name: "aws-cfg",
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("aws-access-key-id", "", "The key to use to communicate")
		set.String("aws-mfa-serial-number", "", "The mfa serial number to communicate to AWS with")
		set.String("aws-default-region", "eu-west-1", "The AWS region that we are accessing resources in")
		set.String("aws-role-arn", "", "The AWS Role ARN to use to communicate to AWS with")
		set.String("aws-secret-access-key", "", "The secret key associated with the access key, to use to communicate to AWS")
		set.String("aws-session-token", "", "The session token to use to communicate to AWS with")
		set.String("aws-endpoint-url", "", "The endpoint URL that is used as AWSs API endpoint")
	},
	Constructor: NewAWSCfg,
}

// NewAWSCfg configures an aws.Config with the given dependency.ConfigGetter
func NewAWSCfg(config dependency.ConfigGetter) aws.Config {
	return aws.Config{
		AccessKey:       config.GetString("aws-access-key-id"),
		MFASerialNumber: config.GetString("aws-mfa-serial-number"),
		Region:          config.GetString("aws-default-region"),
		RoleARN:         config.GetString("aws-role-arn"),
		SecretKey:       config.GetString("aws-secret-access-key"),
		SessionToken:    config.GetString("aws-session-token"),
		EndpointURL:     pointer.ToStringOrNil(config.GetString("aws-endpoint-url")),
	}
}
