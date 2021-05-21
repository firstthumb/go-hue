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

Philips Hue uses local authorization. First you need to create user.

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

[More Examples](https://github.com/firstthumb/go-hue/tree/main/example)

## Coverage

Currently the following services are supported:

- [x] [Lights API](https://developers.meethue.com/develop/hue-api/lights-api/)
  - [x] Get all lights
  - [x] Get new lights
  - [x] Search for new lights
  - [x] Get light attributes and state
  - [x] Set light attributes (rename)
  - [x] Set light state
  - [x] Delete lights
- [ ] [Groups API](https://developers.meethue.com/develop/hue-api/groupds-api/)
  - [x] Get all groups
  - [ ] Create group
  - [ ] Get group attributes
  - [ ] Set group attributes
  - [ ] Set group state
  - [ ] Delete group

## Show your support

Give a ⭐️ if this project helped you!
