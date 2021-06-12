package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/firstthumb/go-hue"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
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
	result, _, _ := client.Lights.Get(context.Background(), "1")
	fmt.Printf("%v\n", result)

	// var v uint16 = 45535
	// color := gamut.Hex("#33B8FF")
	// r, g, b, _ := color.RGBA()
	// v := uint16(r + g + b)

	// a, _ := strconv.ParseInt("#33B8FF", 16, 64)
	// v := uint16(a)

	cl := gamut.Tones(gamut.Hex("#1f11ff"), 12)

	for _, _ = range cl {
		time.Sleep(2 * time.Second)
		// clr, ok := colorful.MakeColor(clll)
		clr := colorful.WarmColor()
		// if !ok {
		// 	fmt.Println("Cannot create.....")
		// }
		// clr, err := colorful.Hex("#1030F3")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		h, s, v := clr.Hsv()

		a := uint16((h * 65535) / 360)
		b := uint8(s * 255)
		c := uint8(v * 255)
		fmt.Println("****")
		fmt.Println(a)
		fmt.Println(b)
		fmt.Println(c)
		fmt.Println("****")

		responses, _, err := client.Lights.SetState(context.Background(), "1", hue.SetStateParams{
			On:  hue.Bool(true),
			Hue: &a,
			Sat: &b,
			Bri: &c,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range responses {
			fmt.Printf("\n\n%+v\n", r)
		}
	}

}
