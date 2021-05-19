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
	groups, _, err := client.Group.GetAll(ctx)
	if err != nil {
		t.Errorf("Group.GetAll returned error: %v", err)
	}
	var result map[string]*Group
	json.Unmarshal(bytes, &result)

	for i, g := range result {
		id, _ := strconv.Atoi(i)
		g.ID = &id
	}

	want := funk.Values(result).([]*Group)
	sort.Slice(want, func(i, j int) bool {
		return *want[i].ID < *want[j].ID
	})

	if !reflect.DeepEqual(groups, want) {
		t.Errorf("Group.GetAll returned %+v, want %+v", groups, want)
	}
}
