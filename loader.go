package ssmloader

import (
	"context"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

const (
	tagKey = "ssm"
)

type Loader struct {
	awsConfig aws.Config
}

func NewLoader(awsConf aws.Config) *Loader {
	return &Loader{
		awsConfig: awsConf,
	}
}

func (l *Loader) Load(ctx context.Context, v any) error {
	ssm := aws_ssm.NewFromConfig(l.awsConfig)

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("v must be a pointer of a struct")
	}

	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Struct {
		return errors.New("v must be a pointer of a struct")
	}

	rt := rv.Type()

	keys := make([]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		keys[i] = rt.Field(i).Tag.Get(tagKey)
	}

	params := map[string]string{}
	{
		var names []string
		for _, key := range keys {
			if len(key) == 0 {
				continue
			}

			names = append(names, key)
		}

		out, err := ssm.GetParameters(ctx, &aws_ssm.GetParametersInput{
			Names:          names,
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return errors.Wrap(err, "failed to get SSM parameters")
		}
		if len(out.InvalidParameters) > 0 {
			return errors.Errorf("invalid params: %s", strings.Join(out.InvalidParameters, ","))
		}

		for _, param := range out.Parameters {
			params[*param.Name] = *param.Value
		}
	}

	for i, key := range keys {
		if len(key) == 0 {
			continue
		}

		rv.Field(i).SetString(params[key])
	}

	return nil
}
