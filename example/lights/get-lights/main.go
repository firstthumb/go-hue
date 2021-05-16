package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/firstthumb/go-hue"
)

func main() {
	host := "<BRIDGE_HOST>"
	token := "<YOUR_USER_TOKEN>"
	client := hue.NewClient(nil, host, token)
	result, _, _ := client.Light.GetAll(context.Background())
	lights, _ := json.Marshal(result)
	fmt.Println(string(lights))
}
