package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateZoneBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPayload := `{"name":"example.org.", "kind": "Native", "masters": [], "nameservers": ["ns1.example.org.", "ns2.example.org."]}`
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("Unable to read body")
		}

		if AreEqualJSON(string(data), expectedPayload) {
			t.Errorf("'%v' should be equal to '%v'\n", expectedPayload, string(data))
		}
		w.WriteHeader(201)
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "changeme", "localhost")
	code, err := c.AddZone("example.org.", []string{"ns1.example.org", "ns2.example.org"})
	if err != nil {
		t.Error(err)
	}
	if code != 201 {
		t.Error()
	}

}
func TestCreateZoneHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "changeme" {
			t.Error("The API key was not configured correctly...")
		}
		w.WriteHeader(201)
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "changeme", "localhost")
	code, err := c.AddZone("example.org.", []string{"ns1.example.org", "ns2.example.org"})
	if err != nil {
		t.Error(err)
	}
	if code != 201 {
		t.Error()
	}

}
func TestEnsureCorrectCreateZoneURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.String()
		if path != "/api/v1/servers/localhost/zones" {
			t.Errorf("'%v' should have been '%v'", path, "/api/v1/servers/localhost/zones")
		}
		w.WriteHeader(201)
	}))
	defer ts.Close()
	c := NewClient(ts.URL, "changeme", "localhost")
	code, err := c.AddZone("example.org.", []string{"ns1.example.org", "ns2.example.org"})
	if err != nil {
		t.Error(err)
	}
	if code != 201 {
		t.Error("Code should be 201, was", code)
	}

}

func AreEqualJSON(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
