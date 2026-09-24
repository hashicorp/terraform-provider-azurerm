package helper

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDnsRecordImporter_DifferentRecordTypes(t *testing.T) {
	tests := []struct {
		name       string
		idPath     string
		recordType recordsets.RecordType
		expectErr  bool
	}{
		{
			name:       "A record with A type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/A/test",
			recordType: recordsets.RecordTypeA,
			expectErr:  false,
		},
		{
			name:       "AAAA record with AAAA type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/AAAA/test",
			recordType: recordsets.RecordTypeAAAA,
			expectErr:  false,
		},
		{
			name:       "CNAME record with CNAME type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/CNAME/test",
			recordType: recordsets.RecordTypeCNAME,
			expectErr:  false,
		},
		{
			name:       "MX record with MX type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/MX/test",
			recordType: recordsets.RecordTypeMX,
			expectErr:  false,
		},
		{
			name:       "NS record with NS type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/NS/test",
			recordType: recordsets.RecordTypeNS,
			expectErr:  false,
		},
		{
			name:       "PTR record with PTR type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/PTR/test",
			recordType: recordsets.RecordTypePTR,
			expectErr:  false,
		},
		{
			name:       "SRV record with SRV type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/SRV/test",
			recordType: recordsets.RecordTypeSRV,
			expectErr:  false,
		},
		{
			name:       "TXT record with TXT type",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/TXT/test",
			recordType: recordsets.RecordTypeTXT,
			expectErr:  false,
		},
		{
			name:       "A record with TXT type (mismatch)",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/A/test",
			recordType: recordsets.RecordTypeTXT,
			expectErr:  true,
		},
		{
			name:       "TXT record with A type (mismatch)",
			idPath:     "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myGroup/providers/Microsoft.Network/dnsZones/example.com/TXT/test",
			recordType: recordsets.RecordTypeA,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			d.SetId(tt.idPath)

			result, err := ResourceDnsRecordImporter(d, tt.recordType)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(result) != 1 {
				t.Errorf("Expected 1 resource data, got: %d", len(result))
			}

			if result[0].Id() != d.Id() {
				t.Errorf("Expected ID %s, got: %s", d.Id(), result[0].Id())
			}
		})
	}
}
