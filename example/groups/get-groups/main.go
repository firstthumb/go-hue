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
	client := hue.NewClient(host, token, nil)
	result, _, _ := client.Group.GetAll(context.Background())
	groups, _ := json.Marshal(result)
	fmt.Println(string(groups))
}
