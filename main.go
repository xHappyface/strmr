package main

import (
    "fmt"
	"log"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
)

func main() {
	fmt.Println("Hello World!")
	client, err := goobs.New("localhost:4455", goobs.WithPassword("test123"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	version, err := client.General.GetVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	fmt.Printf("Websocket server version: %s\n", version.ObsWebSocketVersion)

	resp, _ := client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: "status",
		InputSettings: map[string]interface{}{
			"text": "NEW!!",
		},
	})
	fmt.Println(resp)
}