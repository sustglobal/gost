package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrFieldMissing         = errors.New("event field not found")
	ErrFieldIncorrectType   = errors.New("event field is of incorrect type")
	ErrFieldUnexpectedValue = errors.New("event field contains unexpected value")
)

type EventType string
type EventFieldKey string

func NewEvent(typ EventType, fields ...EventField) *Event {
	return &Event{
		Type:   typ,
		Fields: fields,
	}
}

type Event struct {
	Type   EventType    `json:"type"`
	Fields []EventField `json:"fields"`
}

func (ev *Event) Field(key EventFieldKey) (interface{}, error) {
	for _, ef := range ev.Fields {
		if ef.Key == key {
			return ef.Value, nil
		}
	}
	return nil, ErrFieldMissing
}

func (ev *Event) IntField(key EventFieldKey) (int, error) {
	val, err := ev.Field(key)
	if err != nil {
		return 0, err
	}
	fv, ok := val.(float64)
	if !ok {
		return 0, ErrFieldIncorrectType
	}
	iv := int(fv)
	return iv, nil
}

func (ev *Event) StringField(key EventFieldKey) (string, error) {
	val, err := ev.Field(key)
	if err != nil {
		return "", err
	}
	sv, ok := val.(string)
	if !ok {
		return "", ErrFieldIncorrectType
	}
	return sv, nil
}

func (ev *Event) TimeField(key EventFieldKey) (time.Time, error) {
	val, err := ev.Field(key)
	if err != nil {
		return time.Time{}, err
	}
	sv, ok := val.(string)
	if !ok {
		return time.Time{}, ErrFieldIncorrectType
	}
	tv, err := time.Parse(time.RFC3339, sv)
	if err != nil {
		return time.Time{}, ErrFieldUnexpectedValue
	}
	if tv.IsZero() {
		return time.Time{}, ErrFieldUnexpectedValue
	}
	return tv, nil
}

func (ev *Event) JSONField(key EventFieldKey, dst interface{}) error {
	val, err := ev.Field(key)
	if err != nil {
		return err
	}

	//NOTE(bcwaldon): must marshal here to get from interface{} to something
	// we can try unmarshalling again into the provided dst
	bval, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed marshaling partially-decoded JSON field: %v", err)
	}

	if err := json.Unmarshal(bval, dst); err != nil {
		return ErrFieldUnexpectedValue
	}

	return nil
}

type EventField struct {
	Key   EventFieldKey `json:"key"`
	Value interface{}   `json:"value"`
}

func Field(k EventFieldKey, v interface{}) EventField {
	return EventField{Key: k, Value: v}
}
