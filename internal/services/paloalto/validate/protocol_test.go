// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestProtocolWithPort(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErrors int
	}{
		{
			name:       "single port TCP",
			input:      "TCP:80",
			wantErrors: 0,
		},
		{
			name:       "single port UDP",
			input:      "UDP:443",
			wantErrors: 0,
		},
		{
			name:       "port range TCP",
			input:      "TCP:1024-1206",
			wantErrors: 0,
		},
		{
			name:       "port range UDP",
			input:      "UDP:5000-5100",
			wantErrors: 0,
		},
		{
			name:       "any keyword",
			input:      "any",
			wantErrors: 0,
		},
		{
			name:       "application-default keyword",
			input:      "application-default",
			wantErrors: 0,
		},
		{
			name:       "invalid single port zero",
			input:      "TCP:0",
			wantErrors: 1,
		},
		{
			name:       "invalid port too high",
			input:      "TCP:70000",
			wantErrors: 1,
		},
		{
			name:       "invalid start greater than end",
			input:      "TCP:2000-1000",
			wantErrors: 1,
		},
		{
			name:       "invalid start port zero in range",
			input:      "TCP:0-100",
			wantErrors: 1,
		},
		{
			name:       "invalid end port zero in range",
			input:      "TCP:100-0",
			wantErrors: 1,
		},
		{
			name:       "invalid protocol",
			input:      "ICMP:80",
			wantErrors: 1,
		},
		{
			name:       "invalid missing port",
			input:      "TCP:",
			wantErrors: 1,
		},
		{
			name:       "invalid missing colon",
			input:      "TCP80",
			wantErrors: 1,
		},
		{
			name:       "invalid malformed range",
			input:      "TCP:100-200-300",
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := ProtocolWithPort(tt.input, "test_field")
			if len(errors) != tt.wantErrors {
				t.Errorf("ProtocolWithPort(%q) got %d errors, want %d", tt.input, len(errors), tt.wantErrors)
			}
		})
	}
}
