package utils

import (
	"net/http"
	"testing"

	"github.com/Azure/go-autorest/autorest"
)

func TestResponseNotFound_DroppedConnection(t *testing.T) {
	resp := autorest.Response{}
	if ResponseWasNotFound(resp) {
		t.Fatalf("responseWasNotFound should return `false` for a dropped connection")
	}
}

func TestResponseNotFound_StatusCodes(t *testing.T) {
	testCases := []struct {
		statusCode     int
		expectedResult bool
	}{
		{http.StatusOK, false},
		{http.StatusInternalServerError, false},
		{http.StatusNotFound, true},
	}

	for _, test := range testCases {
		resp := autorest.Response{
			Response: &http.Response{
				StatusCode: test.statusCode,
			},
		}
		result := ResponseWasNotFound(resp)
		if test.expectedResult != result {
			t.Fatalf("Expected '%+v' for status code '%d' - got '%+v'",
				test.expectedResult, test.statusCode, result)
		}
	}
}
