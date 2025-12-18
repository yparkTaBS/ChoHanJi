package Logging

import (
	"ChoHanJi/infrastructure/ContextKeys"
	"context"
	"errors"
)

func GetJobName(ctx context.Context) (string, error) {
	jobName, ok := ctx.Value(ContextKeys.JobName).(string)
	if !ok {
		return "", errors.New("either the jobName is not found or is not of string")
	}
	return jobName, nil
}
