package compute

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestParseUsernameFromAuthorizedKeysPath(t *testing.T) {
	testData := []struct {
		Input    string
		Expected *string
	}{
		{
			Input:    "",
			Expected: nil,
		},
		{
			Input:    "/home/",
			Expected: nil,
		},
		{
			Input:    "/home/username",
			Expected: nil,
		},
		{
			Input:    "/home/username/",
			Expected: nil,
		},
		{
			Input:    "/home/username/.ssh",
			Expected: nil,
		},
		{
			Input:    "/home/username/.ssh/",
			Expected: nil,
		},
		{
			Input:    "/home/username/.ssh/authorized_keys",
			Expected: utils.String("username"),
		},
		{
			Input:    "/home/abc123/.ssh/authorized_keys",
			Expected: utils.String("abc123"),
		},
		{
			Input:    "/home/!!&abc-123!/.ssh/authorized_keys",
			Expected: utils.String("!!&abc-123!"),
		},
	}
	for _, v := range testData {
		actual := parseUsernameFromAuthorizedKeysPath(v.Input)
		if v.Expected != nil && actual == nil {
			t.Fatalf("Expected %q but got nil", *v.Expected)
		}
		if v.Expected == nil && actual != nil {
			t.Fatalf("Expected nil but got %q", *actual)
		}
		if v.Expected != nil && actual != nil && *v.Expected != *actual {
			t.Fatalf("Expected %q but got %q", *v.Expected, *actual)
		}
	}
}
