package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatServerSentEvent(t *testing.T) {
	type TestData struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	testCases := []struct {
		event string
		data  interface{}
		want  string
	}{
		{"event1", "Test data 1", "event: event1\ndata: {\"data\":\"Test data 1\"}\n\n\n"},
		{"event2", 12345, "event: event2\ndata: {\"data\":12345}\n\n\n"},
		{"event3", TestData{"hello", 42}, "event: event3\ndata: {\"data\":{\"field1\":\"hello\",\"field2\":42}}\n\n\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.event, func(t *testing.T) {
			got, err := FormatServerSentEvent(tc.event, tc.data)
			if err != nil {
				t.Fatalf("FormatServerSentEvent(%s, %v) returned an error: %v", tc.event, tc.data, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
