package event

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// used as a fixture for JSON field parsing below
type jsonMessage struct {
	Foo string
	Bar float64
}

func TestFieldCoreTypes(t *testing.T) {
	raw := `{
  "type": "example_type",
  "fields": [
    {"key": "example_str", "value": "XYZ"},
    {"key": "example_int", "value": 3},
	{"key": "example_time", "value": "2022-03-04T01:13:44.1429Z"}
  ]
}`
	var ev Event
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	gotStr, err := ev.StringField(EventFieldKey("example_str"))
	if err != nil {
		t.Errorf("failed unmarshaling JSON: %v", err)
	}

	wantStr := "XYZ"
	if wantStr != gotStr {
		t.Errorf("received incorrect str value: want=%v got=%v", wantStr, gotStr)
	}

	gotInt, err := ev.IntField(EventFieldKey("example_int"))
	if err != nil {
		t.Errorf("failed unmarshaling JSON: %v", err)
	}

	wantInt := 3
	if wantInt != gotInt {
		t.Errorf("received incorrect int value: want=%v got=%v", wantInt, gotInt)
	}

	gotTime, err := ev.TimeField(EventFieldKey("example_time"))
	if err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	wantTime := time.Date(2022, time.March, 4, 1, 13, 44, 142900000, time.UTC)
	if wantTime != gotTime {
		t.Fatalf("received incorrect time value: want=%v got=%v", wantTime, gotTime)
	}
}

func TestFieldCoreTypesIncorrect(t *testing.T) {
	raw := `{
  "type": "example_type",
  "fields": [
    {"key": "example_str", "value": "XYZ"},
    {"key": "example_int", "value": 3},
	{"key": "example_time", "value": "2022-03-04T01:13:44.1429Z"}
  ]
}`
	var ev Event
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	if _, err := ev.StringField(EventFieldKey("example_int")); err != ErrFieldIncorrectType {
		t.Errorf("received incorrect error when parsing int as str: %v", err)
	}

	if _, err := ev.IntField(EventFieldKey("example_time")); err != ErrFieldIncorrectType {
		t.Errorf("received incorrect error when parsing time as int: %v", err)
	}

	if _, err := ev.TimeField(EventFieldKey("example_int")); err != ErrFieldIncorrectType {
		t.Errorf("received incorrect error when parsing int as time: %v", err)
	}
}

func TestFieldTimeUnexpectedValue(t *testing.T) {
	raw := `{
  "type": "example_type",
  "fields": [
    {"key": "example_time1", "value": "XYZ"},
	{"key": "example_time2", "value": "0001-01-01T00:00:00Z"}
  ]
}`
	var ev Event
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	for _, f := range ev.Fields {
		if _, err := ev.TimeField(f.Key); err != ErrFieldUnexpectedValue {
			t.Errorf("received incorrect error when parsing time field: %v", err)
		}
	}
}

func TestFieldJSONType(t *testing.T) {
	raw := `{
  "type": "example_type",
  "fields": [
	{"key": "example_json", "value": {"foo":"XYZ", "bar": 12.43}}
  ]
}`
	var ev Event
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	want := jsonMessage{Foo: "XYZ", Bar: 12.43}
	got := jsonMessage{}

	if err := ev.JSONField(EventFieldKey("example_json"), &got); err != nil {
		t.Errorf("received incorrect error when parsing JSON field: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("parsed value unexpected: want=%+v got=%+v", want, got)
	}
}

func TestFieldJSONUnexpectedValue(t *testing.T) {
	raw := `{
  "type": "example_type",
  "fields": [
	{"key": "example_json_empty_string", "value": ""},
	{"key": "example_json_stringified_json", "value": "{}"}
  ]
}`
	var ev Event
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("failed unmarshaling JSON: %v", err)
	}

	for _, f := range ev.Fields {
		if err := ev.JSONField(f.Key, &jsonMessage{}); err != ErrFieldUnexpectedValue {
			t.Errorf("received incorrect error when parsing JSON field: %v", err)
		}
	}
}
