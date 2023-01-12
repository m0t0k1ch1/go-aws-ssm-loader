package ssmloader

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestConfig struct {
	Field1 string `ssm:"TestString"`
	Field2 string `ssm:"TestSecureString"`
}

func TestLoad(t *testing.T) {
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
		}

		for _, tc := range tcs {
			t.Run(tc.Name, func(t *testing.T) {
				err := Load(tc.Input)
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
			Field2: "TSS",
		}

		var conf TestConfig
		if err := Load(&conf); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(conf, out); len(diff) > 0 {
			t.Errorf("mismatch:\n%s", diff)
		}
	})
}
