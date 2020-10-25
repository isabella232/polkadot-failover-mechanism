package gcp

import (
	"context"
	"fmt"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
)

type Validator struct {
	InstanceName string
	Metric       int
}

func GetValidator(
	ctx context.Context,
	client *monitoring.MetricClient,
	project,
	prefix string,
	checkValue int,
	instanceNames ...string,
) (Validator, error) {
	points, err := GetValidatorMetrics(ctx, client, project, prefix, instanceNames...)

	if err != nil {
		return Validator{}, err
	}

	var validators []Validator

	for instanceName, points := range points {
		if len(points) == 0 {
			continue
		}
		lastPoint := points[0]

		value := lastPoint.Value.GetInt64Value()

		if int(value) == checkValue {
			validators = append(validators, Validator{
				InstanceName: instanceName,
				Metric:       checkValue,
			})
			continue
		}

	}

	switch len(validators) {
	case 0:
		return Validator{}, fmt.Errorf("cannot find validators")
	case 1:
		return validators[0], nil
	default:
		return Validator{}, fmt.Errorf("found %d validators: %#v", len(validators), validators)
	}

}
