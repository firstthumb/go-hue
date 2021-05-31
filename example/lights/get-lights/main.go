package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/firstthumb/go-hue"
)

var (
	clientID = flag.String("clientId", "", "ClientId for Hue for API access.")
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
	result, _, _ := client.Lights.GetAll(context.Background())
	lights, _ := json.Marshal(result)
	fmt.Println(string(lights))
}
