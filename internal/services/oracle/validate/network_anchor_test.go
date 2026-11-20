// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDomainName(t *testing.T) {
	for _, test := range []struct {
		input any
		valid bool
	}{
		{"", false},
		{"test-domain", false},
		{"good-domain.com,example.net", false},
		{"bad-do#main.com", false},
		{"-bad123domain.com", false},
		{"good-domain.com", true},
		{"abc.def.test.example.com", true},
		{42, false},
	} {
		_, errors := DomainName(test.input, "domain_names")
		if (len(errors) == 0) != test.valid {
			t.Errorf("DomainName(%q): expected valid=%t, got errors=%v", test.input, test.valid, errors)
		}
	}
}
