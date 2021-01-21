package network_test

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
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
			Input: "    ",
			Valid: false,
		},
		{
			Name:  "Sql Server",
			Input: "sqlServer",
			Valid: true,
		},
		{
			Name:  "Sql Full Stop Server",
			Input: "sql.Server",
			Valid: true,
		},
		{
			Name:  "Blob Underscore Secondary",
			Input: "blob_secondary",
			Valid: true,
		},
		{
			Name:  "Blob Dash Secondary",
			Input: "blob-secondary",
			Valid: true,
		},
		{
			Name:  "Blob Full Stop Secondary",
			Input: "blob.secondary",
			Valid: true,
		},
		{
			Name:  "Blob Space Secondary",
			Input: "blob secondary",
			Valid: false,
		},
		{
			Name:  "Minimum Value Valid",
			Input: "dfs",
			Valid: true,
		},
		{
			Name:  "Minimum Value Invalid",
			Input: "~~~",
			Valid: false,
		},
		{
			Name:  "Too Short",
			Input: "AB",
			Valid: false,
		},
		{
			Name:  "Minimum Full Stop",
			Input: "S.S",
			Valid: true,
		},
		{
			Name:  "Minimum Underscore",
			Input: "S_S",
			Valid: true,
		},
		{
			Name:  "Minimum Dash",
			Input: "S-S",
			Valid: true,
		},
		{
			Name:  "Max Length",
			Input: "123456789012345678901234567890123456789012345678901234567890123",
			Valid: true,
		},
		{
			Name:  "Too Long",
			Input: "the-name-of-this-subresource-is-way-too-looong-for-this-resource",
			Valid: false,
		},
		{
			Name:  "Start With Dash",
			Input: "-SqlServer",
			Valid: false,
		},
		{
			Name:  "Start With Underscore",
			Input: "_SqlServer",
			Valid: false,
		},
		{
			Name:  "Start Full Stop",
			Input: ".SqlServer",
			Valid: false,
		},
		{
			Name:  "End With Dash",
			Input: "SqlServer-",
			Valid: false,
		},
		{
			Name:  "End With Underscore",
			Input: "SqlServer_",
			Valid: false,
		},
		{
			Name:  "End Full Stop",
			Input: "SqlServer.",
			Valid: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		_, errors := network.ValidatePrivateLinkSubResourceName(v.Input, "private_link_endpoint_subresource")
		isValid := len(errors) == 0
		if v.Valid != isValid {
			t.Fatalf("Expected %t but got %t", v.Valid, isValid)
		}
	}
}
