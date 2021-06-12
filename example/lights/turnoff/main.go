package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/firstthumb/go-hue"
)

var (
	clientID = flag.String("clientId", "", "ClientId for Hue for API access.")
	lightID  = flag.String("lightID", "1", "LightId to turn off.")
)

func main() {
	flag.Parse()

	if len(*clientID) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a Client ID")
		fmt.Println("Flags: ")
		flag.PrintDefaults()
		os.Exit(2)
	}

	host, err := hue.Discover()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not find the bridge")
		os.Exit(1)
	}

	client := hue.NewClient(host, *clientID, nil)
	if err := client.Lights.TurnOff(context.Background(), *lightID); err != nil {
		fmt.Println("Cannot turn off the light", err)
	} else {
		fmt.Println("Turned off the light")
	}
}
