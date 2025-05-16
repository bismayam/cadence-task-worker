package main

import (
	"cadence-task-worker/workflows"
	"context"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "TASK_QUEUE", worker.Options{})
	w.RegisterWorkflow(workflows.CallTaskAPIWorkflow)
	w.RegisterActivity(workflows.CallTaskAPI)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
	defer w.Stop()
	log.Println("✅ Worker started")
	workflowOptions := client.StartWorkflowOptions{
		ID:        "task-workflow-001",
		TaskQueue: "TASK_QUEUE",
	}

	input := workflows.TaskRequest{
		Title:       "Test Task",
		Description: "This is a test task",
		Completed:   false,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CallTaskAPIWorkflow, input)
	if err != nil {
		log.Fatalf("Unable to execute workflow: %v", err)
	}

	log.Println("✅ Workflow started:", "WorkflowID:", we.GetID(), "RunID:", we.GetRunID())

}
