package hue

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	funk "github.com/thoas/go-funk"
)

func TestLightService_GetAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_GetAll.json")
	mux.HandleFunc("/username/lights", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	lights, _, err := client.Light.GetAll(ctx)
	if err != nil {
		t.Errorf("Lights.GetAll returned error: %v", err)
	}
	var result map[string]*Light
	json.Unmarshal(bytes, &result)

	for i, l := range result {
		id, _ := strconv.Atoi(i)
		l.ID = &id
	}

	want := funk.Values(result).([]*Light)
	sort.Slice(want, func(i, j int) bool {
		return *want[i].ID < *want[j].ID
	})

	if !reflect.DeepEqual(lights, want) {
		t.Errorf("Lights.GetAll returned %+v, want %+v", lights, want)
	}
}

func TestLightService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_Get.json")
	mux.HandleFunc("/username/lights/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	light, _, err := client.Light.Get(ctx, "1")
	if err != nil {
		t.Errorf("Light.Get returned error: %v", err)
	}
	want := &Light{}
	json.Unmarshal(bytes, want)

	if !reflect.DeepEqual(light, want) {
		t.Errorf("Light.Get returned %+v, want %+v", light, want)
	}
}

func TestLightService_GetNew(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_GetNew.json")
	mux.HandleFunc("/username/lights/new", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	lights, _, err := client.Light.GetNew(ctx)
	if err != nil {
		t.Errorf("Light.GetNew returned error: %v", err)
	}

	want := []*Light{
		{
			ID:   Int(8),
			Name: String("new lamb"),
		},
	}

	if !reflect.DeepEqual(lights, want) {
		t.Errorf("Light.GetNew returned %+v, want %+v", lights, want)
	}
}

func TestLightService_Search(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_Search.json")
	mux.HandleFunc("/username/lights", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Light.Search(ctx)
	if err != nil {
		t.Errorf("Light.Search returned error: %v", err)
	}

	want := true

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Light.Search returned %+v, want %+v", got, want)
	}
}

func TestLightService_Rename(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_Rename.json")
	mux.HandleFunc("/username/lights/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Light.Rename(ctx, "1", "new_name")
	if err != nil {
		t.Errorf("Light.Rename returned error: %v", err)
	}

	want := true

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Light.Rename returned %+v, want %+v", got, want)
	}
}

func TestLightService_SetState(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Light_SetState.json")
	mux.HandleFunc("/username/lights/1/state", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Light.SetState(ctx, "1", &SetStateRequest{On: Bool(false), Bri: UInt8(200)})
	if err != nil {
		t.Errorf("Light.SetState returned error: %v", err)
	}

	want := []*ApiResponse{
		{
			Success: map[string]interface{}{
				"/lights/1/state/on": false,
			},
		},
		{
			Success: map[string]interface{}{
				"/lights/1/state/bri": float64(200),
			},
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("Light.SetState returned %+v, want %+v", got, want)
	}
}

func TestLightService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/username/lights/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(""))
	})

	ctx := context.Background()
	got, resp, err := client.Light.Delete(ctx, "1")
	if err != nil {
		t.Errorf("Light.Delete returned error: %v", err)
	}

	want := true

	if !cmp.Equal(got, want) {
		t.Errorf("Light.Delete returned %+v, want %+v", got, want)
	}

	if !cmp.Equal(resp.StatusCode, http.StatusOK) {
		t.Errorf("Light.Delete returned status code %+v, want %+v", resp.StatusCode, http.StatusOK)
	}
}
