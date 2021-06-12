package hue

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	funk "github.com/thoas/go-funk"
)

func TestGroupService_GetAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Group_GetAll.json")
	mux.HandleFunc("/username/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.GetAll(ctx)
	if err != nil {
		t.Errorf("Group.GetAll returned error: %+v", err)
	}
	var result map[string]Group
	json.Unmarshal(bytes, &result)

	for i, g := range result {
		id, _ := strconv.Atoi(i)
		g.ID = id
	}

	want := funk.Values(result).([]Group)

	cmp.Equal(got, want, cmpopts.SortSlices(func(x, y Light) bool {
		return x.ID > y.ID
	}))
}

func TestGroupService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Group_Get.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s", testGroupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.Get(ctx, testGroupId)
	if err != nil {
		t.Errorf("Group.Get returned error: %+v", err)
	}
	want := &Group{}
	json.Unmarshal(bytes, want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Group.Get returned %+v, want %+v", got, want)
	}
}

func TestGroupService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	groupId := "1"
	bytes, _ := ioutil.ReadFile("testdata/Group_Update.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s", groupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.Update(ctx, groupId, String("Bedroom"), []string{testLightId}, nil)
	if err != nil {
		t.Errorf("Group.Update returned error: %+v", err)
	}
	want := true

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Group.Update returned %+v, want %+v", got, want)
	}
}

func TestGroupService_SetState(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Group_SetState.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s/action", testGroupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	got, _, err := client.Groups.SetState(ctx, testGroupId, SetStateParams{On: Bool(true)})
	if err != nil {
		t.Errorf("Group.SetState returned error: %+v", err)
	}

	want := []ApiResponse{
		{
			Success: map[string]interface{}{
				"/groups/1/action/on": true,
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Group.SetState returned %+v, want %+v", got, want)
	}
}

func TestGroupService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	bytes, _ := ioutil.ReadFile("testdata/Group_Delete.json")
	mux.HandleFunc(fmt.Sprintf("/username/groups/%s", testGroupId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	})

	ctx := context.Background()
	_, err := client.Groups.Delete(ctx, testGroupId)
	if err != nil {
		t.Errorf("Group.Delete returned error: %+v", err)
	}
}
