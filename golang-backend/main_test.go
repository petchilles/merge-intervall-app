package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  []Interval
	}{
		{
			name: "overlapping intervals",
			intervals: []Interval{
				{Start: 25, End: 30},
				{Start: 2, End: 19},
				{Start: 14, End: 23},
				{Start: 4, End: 8},
			},
			expected: []Interval{
				{Start: 2, End: 23},
				{Start: 25, End: 30},
			},
		},
		{
			name: "adjacent intervals",
			intervals: []Interval{
				{Start: 25, End: 30},
				{Start: 35, End: 40},
				{Start: 20, End: 25},
				{Start: 4, End: 8},
			},
			expected: []Interval{
				{Start: 4, End: 8},
				{Start: 20, End: 30},
				{Start: 35, End: 40},
			},
		},
		{
			name: "intervals having same starting point",
			intervals: []Interval{
				{Start: 25, End: 30},
				{Start: 25, End: 40},
				{Start: 25, End: 35},
				{Start: 25, End: 27},
			},
			expected: []Interval{
				{Start: 25, End: 40},
			},
		},
		{
			name: "no overlapping intervals",
			intervals: []Interval{
				{Start: 1, End: 5},
				{Start: 10, End: 15},
			},
			expected: []Interval{
				{Start: 1, End: 5},
				{Start: 10, End: 15},
			},
		},
		{
			name: "negative intervals",
			intervals: []Interval{
				{Start: -1, End: 5},
				{Start: -10, End: -5},
				{Start: -1, End: 2},
				{Start: -15, End: -5},
				{Start: -10, End: -3},
			},
			expected: []Interval{
				{Start: -15, End: -3},
				{Start: -1, End: 5},
			},
		},
		{
			name:      "empty input",
			intervals: []Interval{},
			expected:  []Interval{},
		},
		{
			name: "single interval",
			intervals: []Interval{
				{Start: 5, End: 10},
			},
			expected: []Interval{
				{Start: 5, End: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeIntervals(tt.intervals)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to create a request and response recorder
func executeRequest(req *http.Request, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)
	return rr
}

func TestMergeHandler_Success(t *testing.T) {
	// Valid input intervals
	intervals := MergeRequest{
		Intervals: []Interval{
			{Start: 25, End: 30},
			{Start: 2, End: 19},
			{Start: 14, End: 23},
			{Start: 4, End: 8},
		},
	}

	// Encode intervals to JSON
	jsonIntervals, err := json.Marshal(intervals)
	if err != nil {
		t.Fatalf("Failed to marshal intervals: %v", err)
	}

	req, err := http.NewRequest("POST", "/merge", bytes.NewBuffer(jsonIntervals))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute the request
	response := executeRequest(req, mergeHandler)

	// Check status code
	assert.Equal(t, http.StatusOK, response.Code)

	// Decode response
	var mergeResponse MergeResponse
	err = json.NewDecoder(response.Body).Decode(&mergeResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check response data
	expected := []Interval{{Start: 2, End: 23}, {Start: 25, End: 30}}
	assert.Equal(t, expected, mergeResponse.Result)
}

func TestMergeHandler_BadRequest(t *testing.T) {
	// Invalid JSON input
	invalidJSON := []byte(`{"intervals": [ { "start": 25, "end": 30 }, { "start": 2, "end": 19 }`)

	req, err := http.NewRequest("POST", "/merge", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute the request
	response := executeRequest(req, mergeHandler)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, response.Code)

	// Check response body
	expectedError := "unexpected EOF\n"
	assert.Equal(t, expectedError, response.Body.String())
}

func TestMergeHandler_InternalServerError(t *testing.T) {
	// Valid input intervals
	intervals := MergeRequest{
		Intervals: []Interval{
			{Start: 25, End: 30},
		},
	}

	// Encode intervals to JSON
	jsonIntervals, err := json.Marshal(intervals)
	if err != nil {
		t.Fatalf("Failed to marshal intervals: %v", err)
	}

	req, err := http.NewRequest("POST", "/merge", bytes.NewBuffer(jsonIntervals))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Mock jsonEncode to produce an error
	oldJsonEncode := jsonEncode
	defer func() { jsonEncode = oldJsonEncode }()
	jsonEncode = func(w http.ResponseWriter, v interface{}) error {
		return errors.New("encoding error")
	}

	// Execute the request
	response := executeRequest(req, mergeHandler)

	// Check status code
	assert.Equal(t, http.StatusInternalServerError, response.Code)

	// Check response body
	expectedError := "encoding error\n"
	assert.Equal(t, expectedError, response.Body.String())
}