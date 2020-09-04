package main

import (
	"log"
	"sync"

	"github.com/csthompson/riptrace/agent/internal/handlers"
	"github.com/csthompson/riptrace/agent/pkg/types"
	"github.com/csthompson/riptrace/agent/service/natsvc"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)

	natsClient := natsvc.GetClient("localhost:4222")
	defer natsClient.Client.Close()

	//Handle requests to the hosts profile
	profileHandler := handlers.ProfileHandler{NatsClient: &natsClient}
	natsClient.RegisterTopicHandler("profile.get", profileHandler.Get)

	//TODO: Add handlers to attach a delve process. Can we embed delve as a package?

	result := types.Profile{}
	natsClient.Request("profile.get", nil, &result)

	for _, p := range result.GoProcs {
		log.Println(p)
	}

	wg.Wait()
}
