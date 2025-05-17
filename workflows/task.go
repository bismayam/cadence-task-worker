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
	Username    string
	Password    string
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

	var token string
	err := workflow.ExecuteActivity(ctx, LoginActivity, input.Username, input.Password).Get(ctx, &token)
	if err != nil {
		return fmt.Errorf("login activity failed: %w", err)
	}

	taskAPIRequest := TaskAPIRequest{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
		Token:       token,
	}

	return workflow.ExecuteActivity(ctx, CallTaskAPI, taskAPIRequest).Get(ctx, nil)
}

type TaskAPIRequest struct {
	Title       string
	Description string
	Completed   bool
	Token       string
}

func CallTaskAPI(ctx context.Context, input TaskAPIRequest) error {
	fmt.Println("===> Inside CallTaskAPI")
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+input.Token).
		SetBody(map[string]interface{}{
			"title":       input.Title,
			"description": input.Description,
			"completed":   input.Completed,
		}).
		Post("http://localhost:8080/api/tasks")
	fmt.Println("===> Response:", resp)
	fmt.Println("===> Response Body:", resp.Body())
	fmt.Println("===> Response Status Code:", resp.StatusCode())
	fmt.Println("===> Response Status:", resp.Status())
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("API error: %s", resp.Status())
	}
	return nil
}

func LoginActivity(ctx context.Context, username, password string) (string, error) {
	fmt.Println("===> Inside LoginActivity")
	client := resty.New()

	loginResp := struct {
		Token string `json:"token"`
	}{}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(&loginResp).
		Post("http://localhost:8080/api/auth/login")

	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("login HTTP error: %s", resp.Status())
	}
	fmt.Println("===> Login successful, token:", loginResp.Token)
	return loginResp.Token, nil
}
