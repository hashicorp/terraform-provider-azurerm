// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

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
			Expected: false,
		},
		{
			Input:    "application-default",
			Expected: false,
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

func TestProtocol(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{Input: "any", Expected: true},
		{Input: "application-default", Expected: true},
		{Input: "TCP:80", Expected: true},
		{Input: "UDP:53", Expected: true},
		{Input: "TCP:1024-1206", Expected: true},
		{Input: "UDP:5000-5100", Expected: true},
		{Input: "TCP:80-80", Expected: true},
		{Input: "Any", Expected: false},
		{Input: "Application-Default", Expected: false},
		{Input: "TCP:0", Expected: false},
		{Input: "TCP:65536", Expected: false},
		{Input: "TCP:1206-1024", Expected: false},
		{Input: "TCP:1024-", Expected: false},
		{Input: "ICMP:80", Expected: false},
		{Input: "TCP", Expected: false},
		{Input: "TCP:80:90", Expected: false},
		{Input: "", Expected: false},
	}
	for _, v := range testCases {
		_, errors := Protocol(v.Input, "protocol")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result for %q to be %t but got %t (and %d errors)", v.Input, v.Expected, result, len(errors))
		}
	}
}

func TestProtocol_nonString(t *testing.T) {
	if _, errors := Protocol(80, "protocol"); len(errors) == 0 {
		t.Fatal("Expected an error for a non-string input to Protocol")
	}
	if _, errors := ProtocolWithPort(80, "protocol_ports"); len(errors) == 0 {
		t.Fatal("Expected an error for a non-string input to ProtocolWithPort")
	}
}

func TestProtocol_errorMessageListsStandaloneValues(t *testing.T) {
	for _, input := range []string{"", "TCP", "anything"} {
		_, errors := Protocol(input, "protocol")
		if len(errors) != 1 {
			t.Fatalf("Expected 1 error for %q but got %d", input, len(errors))
		}
		if msg := errors[0].Error(); !strings.Contains(msg, "`any`") || !strings.Contains(msg, "`application-default`") {
			t.Fatalf("Expected the error for %q to mention `any` and `application-default`, got: %s", input, msg)
		}
	}

	// `protocol_ports` entries don't accept these values, so the message shouldn't suggest them
	_, errors := ProtocolWithPort("TCP", "protocol_ports")
	if len(errors) != 1 {
		t.Fatalf("Expected 1 error but got %d", len(errors))
	}
	if msg := errors[0].Error(); strings.Contains(msg, "`any`") {
		t.Fatalf("Expected the `ProtocolWithPort` error not to mention `any`, got: %s", msg)
	}
}
