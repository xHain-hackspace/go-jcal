package jcal

import (
	"testing"
)

func TestUnmarshalEvent(t *testing.T) {
	testCases := []struct {
		name        string
		properties  []JCalProperty
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
				},
				{
					Name: "dtstamp",
					Values: []interface{}{
						"2021-01-01T00:00:00Z",
					},
				},
				{
					Name: "last-modified",
					Values: []interface{}{
						"2021-01-01T00:00:00Z",
					},
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
				},
				{
					Name: "dtend",
					Values: []interface{}{
						"2022-01-01T11:11:11Z",
					},
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
			_, err := unmarshalEvent(tc.properties)
			if (err != nil) != tc.expectError {
				t.Errorf("unmarshalEvent() error = %v, expectError %v", err, tc.expectError)
			}
		})
	}
}
