package validate

import (
	"strings"
	"testing"
)

func TestFirewallName(t *testing.T) {
	// The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.
	var validNames = []string{
		"a",
		"abc123",
		"a_b_c",
		"hy-ph-en",
		"valid_",
		"v-a_l1.d_",
		strings.Repeat("w", 65),
	}
	for _, v := range validNames {
		_, errors := FirewallName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Firewall Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"_invalid",
		"-invalid",
		".invalid",
		"!invalid",
		"hel!!o",
		"invalid.",
		"invalid-",
		"invalid!",
	}
	for _, v := range invalidNames {
		_, errors := FirewallName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Firewall Name", v)
		}
	}
}
