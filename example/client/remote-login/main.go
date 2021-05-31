package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/firstthumb/go-hue"
)

func main() {
	auth := hue.NewAuthenticator("http://localhost:8181/callback")
	client, err := auth.Authenticate()
	if err != nil {
		panic(err)
	}
	// Client is created but not authenticated by bridge
	username, err := client.CreateRemoteUser()
	if err != nil {
		panic(err)
	}
	// Save "username" for next usage
	// username := "S8uATAguQtowBJnTPpMr8q8nDkskQ6hHbdfAUn1C"
	client.Login(username)
	result, _, _ := client.Lights.GetAll(context.Background())
	lights, _ := json.Marshal(result)
	fmt.Println(string(lights))
}
