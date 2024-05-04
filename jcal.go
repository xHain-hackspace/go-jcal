package jcal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Event struct {
	Created      time.Time `jcal:"created"`
	DtStamp      time.Time `jcal:"dtstamp"`
	LastModified time.Time `jcal:"last-modified"`
	Sequence     int       `jcal:"sequence"`
	UID          string    `jcal:"uid"`
	DtStart      time.Time `jcal:"dtstart"`
	DtEnd        time.Time `jcal:"dtend"`
	Status       string    `jcal:"status"`
	Summary      string    `jcal:"summary"`
	Location     string    `jcal:"location"`
	Description  string    `jcal:"description"`
}

type JCalObject struct {
	ComponentName string         `json:"-"`
	Properties    []JCalProperty `json:"-"`
	SubComponents []JCalObject   `json:"-"`
	Events        []Event        `json:"-"` // Store "vevent" components here
}

type JCalProperty struct {
	Name       string                 `json:"-"`
	Parameters map[string]interface{} `json:"-"`
	TypeName   string                 `json:"-"`
	Values     []interface{}          `json:"-"`
}

// Custom UnmarshalJSON to handle the jCal format
func (obj *JCalObject) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) != 3 {
		return fmt.Errorf("jCal object must have exactly 3 elements")
	}

	// Unmarshal component name
	if err := json.Unmarshal(raw[0], &obj.ComponentName); err != nil {
		return err
	}

	// Unmarshal properties
	var props []json.RawMessage
	if err := json.Unmarshal(raw[1], &props); err != nil {
		return err
	}
	for _, p := range props {
		var prop JCalProperty
		if err := prop.UnmarshalJSON(p); err != nil {
			return err
		}
		obj.Properties = append(obj.Properties, prop)
	}

	// Unmarshal sub-components
	var comps []json.RawMessage
	if err := json.Unmarshal(raw[2], &comps); err != nil {
		return err
	}
	for _, c := range comps {
		var comp JCalObject
		if err := comp.UnmarshalJSON(c); err != nil {
			return err
		}
		obj.SubComponents = append(obj.SubComponents, comp)
	}

	// Example logic to convert a JCalProperty to an Event field
	for _, comp := range obj.SubComponents {
		if comp.ComponentName == "vevent" {
			event, err := unmarshalEvent(comp.Properties)
			if err != nil {
				return err
			}
			obj.Events = append(obj.Events, event)
		}
	}
	return nil
}

func unmarshalEvent(properties []JCalProperty) (Event, error) {
	// Check if nil
	if properties == nil {
		return Event{}, fmt.Errorf("nil properties")
	}
	if len(properties) == 0 {
		return Event{}, fmt.Errorf("empty properties")
	}
	var event Event
	t := reflect.TypeOf(event)
	v := reflect.ValueOf(&event).Elem()
	for _, prop := range properties {
		if len(prop.Values) != 1 {
			return event, fmt.Errorf("property %s must have exactly 1 value", prop.Name)
		}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("jcal")

			if tag == prop.Name {
				switch field.Type {
				case reflect.TypeOf(time.Time{}):
					// Parse date-time strings into time.Time
					strVal, ok := prop.Values[0].(string)
					if !ok {
						return event, fmt.Errorf("%s can not be interpreted as string", prop.Values[0])
					}
					// try RFC3339 first but also support simple dates without times
					parseOk := false
					var err error
					var timeVal time.Time
					for _, timeformat := range []string{time.RFC3339, "2006-01-02"} {
						timeVal, err = time.Parse(timeformat, strVal)
						if err == nil {
							parseOk = true
							break
						}
					}
					if !parseOk {
						return event, err
					}
					v.Field(i).Set(reflect.ValueOf(timeVal))
				case reflect.TypeOf(""):
					if strVal, ok := prop.Values[0].(string); ok {
						v.Field(i).SetString(strVal)
					}
				case reflect.TypeOf(0), reflect.TypeOf(1.2):
					floatVal, floatOk := prop.Values[0].(float64)
					intVal, intOk := prop.Values[0].(int)
					var val int64
					if floatOk {
						val = int64(floatVal)
					} else if intOk {
						val = int64(intVal)
					} else {
						return event, fmt.Errorf("property %s must be a number", prop.Name)
					}
					v.Field(i).SetInt(val)
				}
				break // Move to the next property once a match is found
			}
		}
	}

	// check if dtstart dtend and summary are set
	if event.DtStart.IsZero() {
		return event, fmt.Errorf("dtstart is not set")
	}
	if event.DtEnd.IsZero() {
		return event, fmt.Errorf("dtend is not set")
	}
	if event.Summary == "" {
		return event, fmt.Errorf("summary is not set")
	}

	return event, nil
}

func (prop *JCalProperty) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 4 {
		return fmt.Errorf("jCal property must have at least 4 elements")
	}

	// Unmarshal property name
	if err := json.Unmarshal(raw[0], &prop.Name); err != nil {
		return err
	}

	// Unmarshal parameters
	if err := json.Unmarshal(raw[1], &prop.Parameters); err != nil {
		return err
	}

	// Unmarshal type name
	if err := json.Unmarshal(raw[2], &prop.TypeName); err != nil {
		return err
	}

	// Unmarshal values (remaining elements)
	for _, v := range raw[3:] {
		var value interface{}
		if err := json.Unmarshal(v, &value); err != nil {
			return err
		}
		prop.Values = append(prop.Values, value)
	}

	return nil
}
