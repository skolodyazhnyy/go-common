package configsystem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	url := "http://url.com"
	env := "environment"
	httpc := &http.Client{}

	cli := NewClient(url, env, httpc)
	expected := &Client{
		url:        url,
		env:        env,
		httpClient: httpc,
	}

	if !reflect.DeepEqual(cli, expected) {
		t.Error("Unexpected struct received")
	}
}

func TestClients(t *testing.T) {
	type ExpectedType struct {
		Value string
	}

	expectedValue := &ExpectedType{Value: "expected"}

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" || req.RequestURI != "/client/snoopy/environment/env/scope/loc/merged/key" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		d, _ := json.Marshal(expectedValue)
		_, _ = rw.Write(d)
	}))

	defer srv.Close()

	cli := NewClient(srv.URL, "env", &http.Client{})

	var retrieved ExpectedType

	err := cli.Value("snoopy", "loc", "key", &retrieved)
	if err != nil {
		t.Fatal("Request not processed", err)
	}

	if retrieved.Value != expectedValue.Value {
		t.Errorf("Expecting %s got %s instead", expectedValue.Value, retrieved.Value)
	}

	err = cli.Value("snoopy", "unknown-environment", "key", &retrieved)
	if err == nil {
		t.Fatal("Bad requests should return error")
	}
}

func TestValue(t *testing.T) {
	type ExpectedType struct {
		Value string
	}

	expectedValue := &ExpectedType{Value: "expected"}

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.RequestURI != "/client/snoopy/environment/env/scope/loc/merged/key" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		d, _ := json.Marshal(expectedValue)
		_, _ = rw.Write(d)
	}))

	defer srv.Close()

	cli := NewClient(srv.URL, "env", &http.Client{})

	var retrieved ExpectedType

	err := cli.Value("snoopy", "loc", "key", &retrieved)
	if err != nil {
		t.Fatal("Request not processed", err)
	}

	if retrieved.Value != expectedValue.Value {
		t.Errorf("Expecting %s got %s instead", expectedValue.Value, retrieved.Value)
	}

	err = cli.Value("snoopy", "undefined", "key", &retrieved)
	if err == nil {
		t.Fatal("Bad requests should return error")
	}
}

func TestValueWhenReturnedValueIsAString(t *testing.T) {
	expectedValue := "expected value"

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.RequestURI != "/client/snoopy/environment/env/scope/loc/merged/key" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = rw.Write([]byte(expectedValue))
	}))

	defer srv.Close()

	cli := NewClient(srv.URL, "env", &http.Client{})

	retrieved, err := cli.ValueAsString("snoopy", "loc", "key")
	if err != nil {
		t.Fatal("Request not processed", err)
	}

	if retrieved != expectedValue {
		t.Errorf("Expecting %s got %s instead", expectedValue, retrieved)
	}
}