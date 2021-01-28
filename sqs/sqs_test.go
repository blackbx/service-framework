package sqs_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/sqs"
	"github.com/NYTimes/gizmo/config/aws"
	aws2 "github.com/NYTimes/gizmo/pubsub/aws"
	"github.com/spf13/cobra"
)

func TestNewSQSConfig(t *testing.T) {
	cmd := &cobra.Command{}
	sqs.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewFactory().Configure(cmd)
	if err != nil {
		t.Fatal(err)
	}

	awscfg := aws.Config{
		AccessKey: "foo",
		SecretKey: "bar",
		Region:    "eu-west-1",
	}

	expectedCfg := aws2.SQSConfig{
		Config:              awscfg,
		QueueName:           "",
		QueueOwnerAccountID: "",
		QueueURL:            "",
		MaxMessages:         pointer.ToInt64(10),
		TimeoutSeconds:      pointer.ToInt64(2),
		SleepInterval:       pointer.ToDuration(2 * time.Second),
		DeleteBufferSize:    pointer.ToInt(0),
		ConsumeBase64:       pointer.ToBool(false),
	}
	gotCfg := sqs.NewSQSConfig(awscfg, cfg)
	if !reflect.DeepEqual(expectedCfg, gotCfg) {
		t.Fatalf("expected (%v), got (%v)", expectedCfg, gotCfg)
	}
}
