package jcal

import (
	"testing"
)

func TestUnmarshalEvent(t *testing.T) {
	testCases := []struct {
		name        string
		properties  []JCalProperty
		isAllDay    bool
		expectError bool
	}{
		{
			name: "valid properties",
			properties: []JCalProperty{
				{
					Name: "created",
					Values: []interface{}{
						"2021-01-01T00:00:00Z",
					},
					TypeName: "date-time",
				},
				{
					Name: "dtstamp",
					Values: []interface{}{
						"2021-01-01T00:00:00Z",
					},
					TypeName: "date-time",
				},
				{
					Name: "last-modified",
					Values: []interface{}{
						"2021-01-01T00:00:00Z",
					},
					TypeName: "date-time",
				},
				{
					Name: "sequence",
					Values: []interface{}{
						1,
					},
				},
				{
					Name: "uid",
					Values: []interface{}{
						"test-uid",
					},
				},
				{
					Name: "dtstart",
					Values: []interface{}{
						"2022-01-01T00:00:00Z",
					},
					TypeName: "date-time",
				},
				{
					Name: "dtend",
					Values: []interface{}{
						"2022-01-01T11:11:11Z",
					},
					TypeName: "date-time",
				},
				{
					Name: "status",
					Values: []interface{}{
						"confirmed",
					},
				},
				{
					Name: "summary",
					Values: []interface{}{
						"Test Event",
					},
				},
				{
					Name: "location",
					Values: []interface{}{
						"Berlin",
					},
				},
				{
					Name: "description",
					Values: []interface{}{
						"Berlin",
					},
				},
			},
			expectError: false,
		},
		{
			name:     "whole day event",
			isAllDay: true,
			properties: []JCalProperty{
				{
					Name: "dtstart",
					Values: []interface{}{
						"2024-05-04",
					},
					TypeName: "date",
				},
				{
					Name: "dtend",
					Values: []interface{}{
						"2024-05-04",
					},
					TypeName: "date",
				},
				{
					Name: "status",
					Values: []interface{}{
						"confirmed",
					},
				},
				{
					Name: "summary",
					Values: []interface{}{
						"Whole day Event",
					},
				},
				{
					Name: "location",
					Values: []interface{}{
						"Berlin",
					},
				},
				{
					Name: "description",
					Values: []interface{}{
						"Berlin",
					},
				},
			},
			expectError: false,
		},
		{
			name: "wrong typename for time format",
			properties: []JCalProperty{
				{
					Name: "dtstart",
					Values: []interface{}{
						"2024-05-04",
					},
					TypeName: "date-time",
				},
			},
			expectError: true,
		},
		{
			name: "unrecognized typename",
			properties: []JCalProperty{
				{
					Name: "dtstart",
					Values: []interface{}{
						"2024-05-04",
					},
					TypeName: "foo",
				},
			},
			expectError: true,
		},
		{
			name:        "empty properties",
			properties:  []JCalProperty{},
			expectError: true,
		},
		{
			name: "invalid date",
			properties: []JCalProperty{
				{
					Name: "dtstart",
					Values: []interface{}{
						"invalid date",
					},
				},
			},
			expectError: true,
		},
		{
			name: "multiple dates",
			properties: []JCalProperty{
				{
					Name: "dtstart",
					Values: []interface{}{
						"2022-01-01T00:00:00Z",
						"2022-01-02T00:00:00Z",
					},
				},
			},
			expectError: true,
		},
		{
			name: "no date",
			properties: []JCalProperty{
				{
					Name: "summary",
					Values: []interface{}{
						"Test Event",
					},
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := unmarshalEvent(tc.properties)
			if (err != nil) != tc.expectError {
				t.Errorf("unmarshalEvent() error = %v, expectError %v", err, tc.expectError)
			}
			if event.IsAllDay != tc.isAllDay {
				t.Errorf("unmarshalEvent() IsAllDay property does not match: expected: %v, actual: %v", tc.isAllDay, event.IsAllDay)
			}
		})
	}
}
