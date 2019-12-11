package network

import (
	"testing"
)

func TestValidatePrivateLinkSubResourceName(t *testing.T) {
	testData := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Empty Value",
			Input: "",
			Valid: false,
		},
		{
			Name:  "Whitespace",
			Input: " ",
			Valid: false,
		},
		{
			Name:  "Sql Server",
			Input: "sqlServer",
			Valid: true,
		},
		{
			Name:  "Sql Server All Stop",
			Input: "sql.Server",
			Valid: true,
		},
		{
			Name:  "Blob Secondary",
			Input: "blob_secondary",
			Valid: true,
		},
		{
			Name:  "Blob Secondary Invalid",
			Input: "blob-secondary",
			Valid: false,
		},
		{
			Name:  "Blob Secondary Space",
			Input: "blob secondary",
			Valid: false,
		},
		{
			Name:  "Minimum Value Valid",
			Input: "A",
			Valid: true,
		},
		{
			Name:  "Minimum Value Invalid",
			Input: "~",
			Valid: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		_, errors := ValidatePrivateLinkSubResourceName(v.Input, "private_link_endpoint_subresource")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
