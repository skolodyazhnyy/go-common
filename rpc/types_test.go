package rpc

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeMarshal(t *testing.T) {
	tm := Time(time.Unix(1503659291, 0).UTC())

	data, err := json.Marshal(tm)
	if err != nil {
		t.Errorf("An error occured while marshalling time: %v", err)
	}

	got := string(data)
	want := `"2017-08-25T11:08:11Z"`

	if got != want {
		t.Fatalf("Serialized time does not match format, got: %v, want: %v", got, want)
	}
}

func TestTimeUnmarshal(t *testing.T) {
	tm := Time{}

	if err := json.Unmarshal([]byte(`"2017-08-25T11:08:11Z"`), &tm); err != nil {
		t.Errorf("An error occured while unmarshalling time: %v", err)
	}

	got := time.Time(tm)
	want := time.Unix(1503659291, 0)

	if got.Unix() != want.Unix() {
		t.Fatalf("Serialized time does not match format, got: %v, want: %v", got, want)
	}
}

func TestDateMarshal(t *testing.T) {
	tm := Date(time.Unix(1503619200, 0).UTC())

	data, err := json.Marshal(tm)
	if err != nil {
		t.Errorf("An error occured while marshalling time: %v", err)
	}

	got := string(data)
	want := `"2017-08-25"`

	if got != want {
		t.Fatalf("Serialized time does not match format, got: %v, want: %v", got, want)
	}
}

func TestDateUnmarshal(t *testing.T) {
	tm := Date{}

	if err := json.Unmarshal([]byte(`"2017-08-25"`), &tm); err != nil {
		t.Errorf("An error occured while unmarshalling time: %v", err)
	}

	got := time.Time(tm)
	want := time.Unix(1503619200, 0)

	if got.Unix() != want.Unix() {
		t.Fatalf("Serialized time does not match format, got: %v, want: %v", got.Unix(), want)
	}
}
