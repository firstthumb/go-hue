<h1 align="center">go-hue</h1>
<h3 align="center">*** Work In Progress ***</h3>

<p align="center">
  <a href="https://github.com/firstthumb/go-hue/commits/main">
    <img src="https://img.shields.io/github/last-commit/firstthumb/go-hue.svg" target="_blank" />
  </a>
  <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/firstthumb/go-hue">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/firstthumb/go-hue">
  <a href="http://godoc.org/github.com/firstthumb/go-hue">
    <img src="https://godoc.org/github.com/firstthumb/go-hue?status.svg" target="_blank" />
  </a>
  <a href="https://github.com/firstthumb/go-hue/issues?q=is%3Apr+is%3Aclosed">
    <img alt="GitHub closed pull requests" src="https://img.shields.io/github/issues-pr-closed-raw/firstthumb/go-hue"> 
  </a>
  <a href="https://github.com/firstthumb/go-hue/pulls">
    <img alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr/firstthumb/go-hue">
  </a>
  <a href="https://github.com/firstthumb/go-hue/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/firstthumb/go-hue">
  </a>
  <a href="https://github.com/firstthumb/go-hue/graphs/contributors">
    <img alt="GitHub contributors" src="https://img.shields.io/github/contributors/firstthumb/go-hue">
  </a>
  <a href="https://github.com/firstthumb/go-hue/blob/main/LICENSE.md">
    <img alt="License: BSD" src="https://img.shields.io/badge/license-MIT-green.svg" target="_blank" />
  </a>
</p>

> go-hue is a Go client library for accessing the [Philips Hue API](https://developers.meethue.com/develop/hue-api/)

## Install

```sh
go get github.com/firstthumb/go-hue
```

## Authentication

Philips Hue uses local and remote authorization. First you need to create user.

## Usage

Import the package into your project.

```Go
import "github.com/firstthumb/go-hue"
```

Use existing user and access Hue services. For example:

```Go
// Discover your network and finds the first bridge
host, _ := hue.Discover()
client := hue.NewClient(host, "<YOUR USER TOKEN>", nil)
lights, resp, err := client.Light.GetAll(context.Background())
```

Or create user. Don't forget to save the clientId 

```Go
// You must press Philips Hue bridge button before
host, _ := hue.Discover()
client, _ := hue.CreateUser(host, "<CLIENT_NAME>", nil)
client.GetClientID() // Save clientID for next time
lights, resp, err := client.Light.GetAll(context.Background())
```

Supports remote API

```Go
// Create your clientId and clientSecret at https://developers.meethue.com/my-apps/
// set your environment variables HUE_CLIENT_ID, HUE_CLIENT_SECRET and HUE_APP_ID
// use the same callback url defined in your app
auth := hue.NewAuthenticator("http://localhost:8181/callback")
client, err := auth.Authenticate()
if err != nil {
  panic(err)
}
	
username, err := client.CreateRemoteUser()
if err != nil {
  panic(err)
}
	
client.Login(username)
result, _, _ := client.Light.GetAll(context.Background())
lights, _ := json.Marshal(result)
fmt.Println(string(lights))
```

[More Examples](https://github.com/firstthumb/go-hue/tree/main/example)

## Coverage

Currently the following services are supported:

- [x] [Remote API](https://developers.meethue.com/develop/hue-api/remote-api-quick-start-guide/)
  - [x] Remote Login
- [x] [Lights API](https://developers.meethue.com/develop/hue-api/lights-api/)
  - [x] Get all lights
  - [x] Get new lights
  - [x] Search for new lights
  - [x] Get light attributes and state
  - [x] Set light attributes (rename)
  - [x] Set light state
  - [x] Delete lights
- [x] [Groups API](https://developers.meethue.com/develop/hue-api/groupds-api/)
  - [x] Get all groups
  - [x] Create group
  - [x] Get group attributes
  - [x] Set group attributes
  - [x] Set group state
  - [x] Delete group
- [ ] [Schedules API](https://developers.meethue.com/develop/hue-api/3-schedules-api/)
- [ ] [Scenes API](https://developers.meethue.com/develop/hue-api/4-scenes/)

## Show your support

Give a ⭐️ if this project helped you!
