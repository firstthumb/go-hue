package hue

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestGroupService_TurnOn(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	groupId := "1"
	bytes, _ := ioutil.ReadFile("testdata/Group_TurnOn.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s/action", groupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.TurnOn(ctx, groupId)
	if err != nil {
		t.Errorf("Group.TurnOn returned error: %+v", err)
	}

	want := []*ApiResponse{
		{
			Success: map[string]interface{}{
				"/groups/1/action/on": true,
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Group.TurnOn returned %+v, want %+v", got, want)
	}
}

func TestGroupService_TurnOff(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	groupId := "1"
	bytes, _ := ioutil.ReadFile("testdata/Group_TurnOff.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s/action", groupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.TurnOff(ctx, groupId)
	if err != nil {
		t.Errorf("Group.TurnOff returned error: %+v", err)
	}

	want := []*ApiResponse{
		{
			Success: map[string]interface{}{
				"/groups/1/action/on": false,
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Group.TurnOff returned %+v, want %+v", got, want)
	}
}

func TestGroupService_TurnOnAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var wg sync.WaitGroup

	groupIds := []string{"1", "2", "3"}
	bytes, _ := ioutil.ReadFile("testdata/Group_TurnOn.json")
	for _, groupId := range groupIds {
		wg.Add(1)
		mux.HandleFunc(fmt.Sprintf("/username/groups/%s/action", groupId), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(bytes))
			wg.Done()
		})
	}

	ctx := context.Background()
	client.Groups.TurnOnAll(ctx, groupIds...)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Errorf("Failed to turn on all groups on time")
	}
}

func TestGroupService_TurnOffAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var wg sync.WaitGroup

	groupIds := []string{"1", "2", "3"}
	bytes, _ := ioutil.ReadFile("testdata/Group_TurnOff.json")
	for _, groupId := range groupIds {
		wg.Add(1)
		mux.HandleFunc(fmt.Sprintf("/username/groups/%s/action", groupId), func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(bytes))
			wg.Done()
		})
	}

	ctx := context.Background()
	client.Groups.TurnOffAll(ctx, groupIds...)

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Errorf("Failed to turn off all groups on time")
	}
}
