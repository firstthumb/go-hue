package hue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/muesli/gamut"
	"github.com/sirupsen/logrus"
)

var (
	baseURLPath = "/api"

	testLightId  = "1"
	testGroupId  = "1"
	testColor    = gamut.Hex("#FF0000")
	testColorHex = "#FF0000"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	url, _ := url.Parse(server.URL)
	client = NewClient(url.Host, "username", &ClientOptions{LogLevel: logrus.DebugLevel})

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func getPayload(t *testing.T, r *http.Request, payload interface{}) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Request payload failed to read")
	}
	err = json.Unmarshal(bytes, payload)
	if err != nil {
		t.Errorf("Request payload failed to unmarshal")
	}
}
