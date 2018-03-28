package response

import (
	"net/http"
	"testing"
)

func TestConflict_DroppedConnection(t *testing.T) {
	resp := http.Response{}
	if WasConflict(&resp) {
		t.Fatalf("wasConflict should return `false` for a dropped connection")
	}
}

func TestConflcit_StatusCodes(t *testing.T) {
	testCases := []struct {
		statusCode     int
		expectedResult bool
	}{
		{http.StatusOK, false},
		{http.StatusInternalServerError, false},
		{http.StatusNotFound, false},
		{http.StatusConflict, true},
	}

	for _, test := range testCases {
		resp := http.Response{
			StatusCode: test.statusCode,
		}
		result := WasConflict(&resp)
		if test.expectedResult != result {
			t.Fatalf("Expected '%+v' for status code '%d' - got '%+v'",
				test.expectedResult, test.statusCode, result)
		}
	}
}

func TestNotFound_DroppedConnection(t *testing.T) {
	resp := http.Response{}
	if WasNotFound(&resp) {
		t.Fatalf("wasNotFound should return `false` for a dropped connection")
	}
}

func TestNotFound_StatusCodes(t *testing.T) {
	testCases := []struct {
		statusCode     int
		expectedResult bool
	}{
		{http.StatusOK, false},
		{http.StatusInternalServerError, false},
		{http.StatusNotFound, true},
	}

	for _, test := range testCases {
		resp := http.Response{
			StatusCode: test.statusCode,
		}
		result := WasNotFound(&resp)
		if test.expectedResult != result {
			t.Fatalf("Expected '%+v' for status code '%d' - got '%+v'",
				test.expectedResult, test.statusCode, result)
		}
	}
}
