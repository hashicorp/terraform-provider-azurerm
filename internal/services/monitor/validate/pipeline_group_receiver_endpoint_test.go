// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestPipelineGroupReceiverEndpoint(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "", Valid: false},
		{Input: "0.0.0.0:514", Valid: true},
		{Input: "0.0.0.0:4317", Valid: true},
		{Input: ":4317", Valid: true},
		{Input: "[::]:4317", Valid: true},
		{Input: "example.com:443", Valid: true},
		{Input: "udp://example.com:5000", Valid: true},
		{Input: "203.0.113.10:5514", Valid: true},
		{Input: "localhost:8080", Valid: true},
		{Input: "127.0.0.1:4317", Valid: true},
		{Input: "10.0.0.5:4317", Valid: true},
		{Input: "192.168.1.10:514", Valid: true},
		{Input: "0.0.0.0", Valid: false},
		{Input: "noport", Valid: false},
		{Input: "0.0.0.0:", Valid: false},
		{Input: ":0", Valid: false},
		{Input: ":65536", Valid: false},
		{Input: ":abc", Valid: false},
		{Input: ":-1", Valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := PipelineGroupReceiverEndpoint(tc.Input, "endpoint")
			valid := len(errors) == 0
			if valid != tc.Valid {
				t.Fatalf("expected %q to be valid=%t, got valid=%t", tc.Input, tc.Valid, valid)
			}
		})
	}
}
