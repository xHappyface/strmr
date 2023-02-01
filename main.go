package main

import (
	"log"

	"github.com/andreykaipov/goobs"
	"github.com/jnrprgmr/dog/pkg/obs"
)

func main() {
	client, err := goobs.New("localhost:4455", goobs.WithPassword("test123"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	obs := obs.New(client)

	obs.SetTask("UPDATE##")
}
