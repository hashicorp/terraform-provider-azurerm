// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
)

func TestValidateRecordTypeID(t *testing.T) {
	cases := []struct {
		RecordType recordsets.RecordType
		Value      string
		Errors     int
	}{
		{
			RecordType: recordsets.RecordTypeA,
			Value:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Network/dnsZones/domain.com/TXT/testrecord",
			Errors:     1,
		},
		{
			RecordType: recordsets.RecordTypeTXT,
			Value:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Network/dnsZones/domain.com/TXT/testrecord",
			Errors:     0,
		},
		{
			RecordType: recordsets.RecordTypeTXT,
			Value:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group",
			Errors:     1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Value, func(t *testing.T) {
			_, errors := ValidateRecordTypeID(tc.RecordType)(tc.Value, "id")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected ValidateRecordTypeID to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
