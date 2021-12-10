package validate

import "testing"

func TestClientID(t *testing.T) {
	testData := []struct {
		input string
		error bool
	}{
		{
			input: "",
			error: true,
		},
		{
			input: " ",
			error: true,
		},
		{
			input: "hello",
			error: false,
		},
		{
			input: "msi",
			error: true,
		},
		{
			input: "Msi",
			error: true,
		},
		{
			input: "MSI",
			error: true,
		},
		{
			input: "abc123",
			error: false,
		},
	}

	for _, v := range testData {
		t.Logf("Testing %q..", v.input)
		_, errors := ClientID(v.input, "client_id")
		hasErrors := len(errors) > 0
		if v.error != hasErrors {
			t.Fatalf("Expected %t but got %t", v.error, hasErrors)
		}
	}
}
