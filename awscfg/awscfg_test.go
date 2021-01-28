package awscfg_test

import (
	"reflect"
	"testing"

	"github.com/BlackBX/service-framework/awscfg"
	"github.com/BlackBX/service-framework/config"
	"github.com/NYTimes/gizmo/config/aws"
	"github.com/spf13/cobra"
)

func TestNewAWSCfg(t *testing.T) {
	cmd := &cobra.Command{}
	awscfg.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewFactory().Configure(cmd)
	if err != nil {
		t.Fatal(err)
	}
	gotCfg := awscfg.NewAWSCfg(cfg)
	expectedCfg := aws.Config{
		Region: "eu-west-1",
	}
	if !reflect.DeepEqual(expectedCfg, gotCfg) {
		t.Fatalf("expected (%v), got (%v)", expectedCfg, gotCfg)
	}
}
