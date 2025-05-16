package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TaskRequest struct {
	Description string
	Title       string
	Completed   bool
}

func CallTaskAPIWorkflow(ctx workflow.Context, input TaskRequest) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 1,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 2,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	return workflow.ExecuteActivity(ctx, CallTaskAPI, input).Get(ctx, nil)
}

func CallTaskAPI(ctx context.Context, input TaskRequest) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"title":       input.Title,
			"description": input.Description,
			"completed":   input.Completed,
		}).
		Post("http://localhost:8080/api/tasks")

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.Status())
	}
	return nil
}
