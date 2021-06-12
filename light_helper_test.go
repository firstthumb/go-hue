package hue

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLightService_TurnOn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_TurnOn.json")
	mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", testLightId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var payload SetStateParams
		getPayload(t, r, &payload)

		assert.Equal(t, true, *payload.On)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	err := client.Lights.TurnOn(ctx, testLightId)
	if err != nil {
		t.Errorf("Light.TurnOn returned error: %+v", err)
	}

	assert.Nil(t, err)
}

func TestLightService_TurnOff(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_TurnOff.json")
	mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", testLightId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var payload SetStateParams
		getPayload(t, r, &payload)

		assert.Equal(t, false, *payload.On)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	err := client.Lights.TurnOff(ctx, testLightId)
	if err != nil {
		t.Errorf("Light.TurnOff returned error: %+v", err)
	}

	assert.Nil(t, err)
}

func TestLightService_TurnOnAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var wg sync.WaitGroup

	lightIds := []string{"1", "2", "3"}
	bytes, _ := ioutil.ReadFile("testdata/Light_TurnOnAll.json")
	for _, lightId := range lightIds {
		wg.Add(1)
		mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", lightId), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")

			var payload SetStateParams
			getPayload(t, r, &payload)

			assert.Equal(t, true, *payload.On)

			w.Header().Set("Content-Type", "application/json")
			tmpl, _ := template.New("test").Parse(string(bytes))
			tmpl.Execute(w, struct{ LightId string }{LightId: lightId})
			wg.Done()
		})
	}

	ctx := context.Background()
	client.Lights.TurnOnAll(ctx, lightIds...)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Errorf("Failed to turn on all lights on time")
	}
}

func TestLightService_TurnOffAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var wg sync.WaitGroup

	lightIds := []string{"1", "2", "3"}
	bytes, _ := ioutil.ReadFile("testdata/Light_TurnOffAll.json")
	for _, lightId := range lightIds {
		wg.Add(1)
		mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", lightId), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")

			var payload SetStateParams
			getPayload(t, r, &payload)

			assert.Equal(t, false, *payload.On)

			w.Header().Set("Content-Type", "application/json")
			tmpl, _ := template.New("test").Parse(string(bytes))
			tmpl.Execute(w, struct{ LightId string }{LightId: lightId})
			wg.Done()
		})
	}

	ctx := context.Background()
	client.Lights.TurnOffAll(ctx, lightIds...)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Errorf("Failed to turn off all lights on time")
	}
}

func TestLightService_SetColor(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_SetColor.json")
	mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", testLightId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var payload SetStateParams
		getPayload(t, r, &payload)

		h, s, b, _ := colorToHSV(testColor)

		assert.Equal(t, true, *payload.On)
		assert.Equal(t, h, *payload.Hue)
		assert.Equal(t, s, *payload.Sat)
		assert.Equal(t, b, *payload.Bri)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	err := client.Lights.SetColor(ctx, testLightId, testColor)
	if err != nil {
		t.Errorf("Light.SetColor returned error: %+v", err)
	}

	assert.Nil(t, err)
}

func TestLightService_SetColorHex(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_SetColor.json")
	mux.HandleFunc(fmt.Sprintf("/username/lights/%s/state", testLightId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var payload SetStateParams
		getPayload(t, r, &payload)

		h, s, b, _ := hexColorToHSV(testColorHex)

		assert.Equal(t, true, *payload.On)
		assert.Equal(t, h, *payload.Hue)
		assert.Equal(t, s, *payload.Sat)
		assert.Equal(t, b, *payload.Bri)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	err := client.Lights.SetColorHex(ctx, testLightId, testColorHex)
	if err != nil {
		t.Errorf("Light.SetColorHex returned error: %+v", err)
	}

	assert.Nil(t, err)
}
