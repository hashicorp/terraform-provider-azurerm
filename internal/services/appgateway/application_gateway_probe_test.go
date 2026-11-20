// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appgateway

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
)

func TestFlattenApplicationGatewayProbes_MatchByProtocol(t *testing.T) {
	defaultMatch := &applicationgateways.ApplicationGatewayProbeHealthResponseMatch{
		Body:        pointer.To(""),
		StatusCodes: &[]string{"200-399"},
	}

	testCases := []struct {
		name              string
		protocol          applicationgateways.ApplicationGatewayProtocol
		expectedMatchSize int
	}{
		{
			name:              "Http retains match block",
			protocol:          applicationgateways.ApplicationGatewayProtocolHTTP,
			expectedMatchSize: 1,
		},
		{
			name:              "Https retains match block",
			protocol:          applicationgateways.ApplicationGatewayProtocolHTTPS,
			expectedMatchSize: 1,
		},
		{
			name:              "Tcp omits API default match block",
			protocol:          applicationgateways.ApplicationGatewayProtocolTcp,
			expectedMatchSize: 0,
		},
		{
			name:              "Tls omits API default match block",
			protocol:          applicationgateways.ApplicationGatewayProtocolTls,
			expectedMatchSize: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &[]applicationgateways.ApplicationGatewayProbe{
				{
					Name: pointer.To("test-probe"),
					Properties: &applicationgateways.ApplicationGatewayProbePropertiesFormat{
						Protocol:           pointer.To(tc.protocol),
						Interval:           pointer.To(int64(30)),
						Timeout:            pointer.To(int64(30)),
						UnhealthyThreshold: pointer.To(int64(3)),
						Match:              defaultMatch,
					},
				},
			}

			flattened := flattenApplicationGatewayProbes(input)
			if len(flattened) != 1 {
				t.Fatalf("expected 1 flattened probe, got %d", len(flattened))
			}

			probeMap := flattened[0].(map[string]any)
			matches, ok := probeMap["match"].([]any)
			if !ok {
				t.Fatalf("expected match to be []any, got %T", probeMap["match"])
			}

			if len(matches) != tc.expectedMatchSize {
				t.Fatalf("expected match length %d for protocol %q, got %d", tc.expectedMatchSize, tc.protocol, len(matches))
			}
		})
	}
}
