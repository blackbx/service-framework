package sqs

import (
	"context"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/NYTimes/gizmo/config/aws"
	"github.com/NYTimes/gizmo/pubsub"
	aws2 "github.com/NYTimes/gizmo/pubsub/aws"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	DefaultSleepInterval = 2 * time.Second
)

// Service is a dependency that provides an SQS Subscriber for a service
var Service = dependency.Service{
	Name: "gizmo-sqs",
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("aws-sqs-queue-name", "", "The name of the SQS Queue you want to read from")
		set.String("aws-sqs-queue-owner-account-id", "", "The account ID of the owner of the SQS Queue")
		set.String("aws-sqs-queue-url", "", "The URL of the SQS Queue")
		set.Int64("aws-sqs-max-messages", 10, "The number of bulk messages the SQSSubscriber will attempt to fetch on each receive.")
		set.Int64("aws-sqs-timeout-seconds", 2, "The number of seconds the SQS client will wait before timing out.")
		set.Duration("aws-sqs-sleep-interval", DefaultSleepInterval, "The time the SQSSubscriber will wait if it sees no messages on the queue.")
		set.Int("aws-sqs-delete-buffer-size", 0, "The limit of messages allowed in the delete buffer before executing a 'delete batch' request.")
		set.Bool("aws-sqs-consume-base64", false, "A flag to signal the subscriber to base64 decode the payload before returning it.")
	},
	Dependencies: fx.Provide(
		NewSQSConfig,
	),
	Constructor: aws2.NewSubscriber,
	InvokeFunc:  Invoke,
}

// StopParams are the parameters required in order to register an invoke
// function, that is used to stop the queue reader.
type StopParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Queue     pubsub.Subscriber
}

// Invoke is the function that is registered with fx in order to gracefully
// shut down the server.
func Invoke(params StopParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Stopping queue")
			err := params.Queue.Stop()
			params.Logger.Info("Queue stopped")
			return err
		},
	})
}

// NewSQSConfig gives you a new instance of aws.SQSConfig for a given dependency.ConfigGetter
func NewSQSConfig(awscfg aws.Config, config dependency.ConfigGetter) aws2.SQSConfig {
	return aws2.SQSConfig{
		Config:              awscfg,
		QueueName:           config.GetString("aws-sqs-queue-name"),
		QueueOwnerAccountID: config.GetString("aws-sqs-queue-owner-account-id"),
		QueueURL:            config.GetString("aws-sqs-queue-url"),
		MaxMessages:         pointer.ToInt64(config.GetInt64("aws-sqs-max-messages")),
		TimeoutSeconds:      pointer.ToInt64(config.GetInt64("aws-sqs-timeout-seconds")),
		SleepInterval:       pointer.ToDuration(config.GetDuration("aws-sqs-sleep-interval")),
		DeleteBufferSize:    pointer.ToInt(config.GetInt("aws-sqs-delete-buffer-size")),
		ConsumeBase64:       pointer.ToBool(config.GetBool("aws-sqs-consume-base64")),
	}
}
