package ssmloader

import (
	"context"
	"testing"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/go-cmp/cmp"
)

type TestConfig struct {
	Field1 string `ssm:"TestString"`
	Field2 string
	Field3 string `ssm:"TestSecureString"`
}

func TestLoad(t *testing.T) {
	ctx := context.Background()

	awsConf, err := aws_config.LoadDefaultConfig(ctx, aws_config.WithRegion("ap-northeast-1"))
	if err != nil {
		t.Fatal(err)
	}

	l := NewLoader(awsConf)

	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			Name         string
			Input        any
			ErrorMessage string
		}{
			{
				"string",
				"string",
				"v must be a pointer of a struct",
			},
			{
				"nil",
				nil,
				"v must be a pointer of a struct",
			},
			{
				"struct",
				TestConfig{},
				"v must be a pointer of a struct",
			},
			{
				"struct pointer with invalid tag",
				&struct {
					Field string `ssm:"InvalidString"`
				}{},
				"invalid params: InvalidString",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				err := l.Load(ctx, tc.Input)
				if err == nil {
					t.Errorf("err must not be nil")
					return
				}
				if diff := cmp.Diff(err.Error(), tc.ErrorMessage); len(diff) > 0 {
					t.Errorf("mismatch:\n%s", diff)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		out := TestConfig{
			Field1: "TS",
			Field3: "TSS",
		}

		var conf TestConfig
		if err := l.Load(ctx, &conf); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(conf, out); len(diff) > 0 {
			t.Errorf("mismatch:\n%s", diff)
		}
	})
}
