// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestProtocolWithPort(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "TCP:80",
			Expected: true,
		},
		{
			Input:    "UDP:53",
			Expected: true,
		},
		{
			Input:    "TCP:1",
			Expected: true,
		},
		{
			Input:    "TCP:65535",
			Expected: true,
		},
		{
			Input:    "any",
			Expected: true,
		},
		{
			Input:    "application-default",
			Expected: true,
		},
		{
			Input:    "TCP:1024-1206",
			Expected: true,
		},
		{
			Input:    "UDP:5000-5100",
			Expected: true,
		},
		{
			Input:    "TCP:80-80",
			Expected: true,
		},
		{
			Input:    "TCP:0",
			Expected: false,
		},
		{
			Input:    "TCP:65536",
			Expected: false,
		},
		{
			Input:    "TCP:abc",
			Expected: false,
		},
		{
			Input:    "TCP:1206-1024",
			Expected: false,
		},
		{
			Input:    "TCP:1024-70000",
			Expected: false,
		},
		{
			Input:    "TCP:1024-",
			Expected: false,
		},
		{
			Input:    "ICMP:80",
			Expected: false,
		},
		{
			Input:    "TCP",
			Expected: false,
		},
		{
			Input:    "TCP:80:90",
			Expected: false,
		},
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "Any",
			Expected: false,
		},
	}
	for _, v := range testCases {
		_, errors := ProtocolWithPort(v.Input, "protocol_ports")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result for %q to be %t but got %t (and %d errors)", v.Input, v.Expected, result, len(errors))
		}
	}
}
