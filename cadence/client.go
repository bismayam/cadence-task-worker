package cadence

import (
	"log"

	"go.temporal.io/sdk/client"
)

var CadenceClient client.Client

func InitCadenceClient() {
	var err error
	CadenceClient, err = client.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
}
