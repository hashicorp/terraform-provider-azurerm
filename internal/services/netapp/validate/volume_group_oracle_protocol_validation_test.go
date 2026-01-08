// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateNetAppVolumeGroupOracleProtocolConversion(t *testing.T) {
	cases := []struct {
		Name                string
		VolumeSpecName      string
		OldProtocols        []string
		NewProtocols        []string
		ExportPolicyRules   []interface{}
		ExpectedErrors      int
		ExpectedErrorString string
	}{
		{
			Name:              "ValidConversionDataNFSv3ToNFSv41",
			VolumeSpecName:    "data",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionDataNFSv41ToNFSv3",
			VolumeSpecName:    "data",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionLogsNFSv3ToNFSv41",
			VolumeSpecName:    "logs",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionLogsNFSv41ToNFSv3",
			VolumeSpecName:    "logs",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionBinaryNFSv3ToNFSv41",
			VolumeSpecName:    "binary",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionBinaryNFSv41ToNFSv3",
			VolumeSpecName:    "binary",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "NoChangeNoValidation",
			VolumeSpecName:    "data",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "InitialCreationNoValidation",
			VolumeSpecName:    "data",
			OldProtocols:      []string{},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:           "ValidConversionWithMatchingExportPolicy",
			VolumeSpecName: "data",
			OldProtocols:   []string{"NFSv3"},
			NewProtocols:   []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocols_enabled": []interface{}{"NFSv4.1"},
				},
			},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidConversionWithExportPolicyUpdate",
			VolumeSpecName: "logs",
			OldProtocols:   []string{"NFSv4.1"},
			NewProtocols:   []string{"NFSv3"},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocol": []interface{}{"NFSv3"}, // Will be updated to match new volume protocol
				},
			},
			ExpectedErrors: 0, // Export policy validation is skipped during protocol conversion
		},
		{
			Name:           "ValidConversionMultipleVolumes",
			VolumeSpecName: "binary",
			OldProtocols:   []string{"NFSv3"},
			NewProtocols:   []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocols_enabled": []interface{}{"NFSv4.1"},
				},
				map[string]interface{}{
					"protocol": []interface{}{"NFSv4.1"},
				},
			},
			ExpectedErrors: 0,
		},
		{
			Name:                "InvalidDualProtocolConversion",
			VolumeSpecName:      "data",
			OldProtocols:        []string{"NFSv3", "CIFS"},
			NewProtocols:        []string{"NFSv4.1"},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      2, // Both dual-protocol and CIFS conversion errors
			ExpectedErrorString: "cannot change the NFS version of a dual-protocol volume",
		},
		{
			Name:                "InvalidCIFSConversion",
			VolumeSpecName:      "logs",
			OldProtocols:        []string{"NFSv3"},
			NewProtocols:        []string{"NFSv4.1", "CIFS"},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      2, // Both dual-protocol and CIFS conversion errors
			ExpectedErrorString: "cannot change the NFS version of a dual-protocol volume",
		},
		{
			Name:           "ValidCustomVolumeSpecName",
			VolumeSpecName: "custom-oracle-volume",
			OldProtocols:   []string{"NFSv3"},
			NewProtocols:   []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocols_enabled": []interface{}{"NFSv4.1"},
				},
			},
			ExpectedErrors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var errors []error

			// Oracle volume groups don't have specific protocol restrictions like SAP HANA
			// They only need to validate protocol conversion rules
			if len(tc.OldProtocols) > 0 && !slicesEqual(tc.OldProtocols, tc.NewProtocols) {
				// For volume groups, kerberos and data replication are not directly supported
				var kerberosEnabled bool
				var dataReplication []interface{}

				conversionErrors := ValidateNetAppVolumeProtocolConversion(tc.OldProtocols, tc.NewProtocols, kerberosEnabled, dataReplication, tc.ExportPolicyRules)
				errors = append(errors, conversionErrors...)
			}

			if len(errors) != tc.ExpectedErrors {
				t.Fatalf("expected %d error(s), got %d: %v", tc.ExpectedErrors, len(errors), errors)
			}

			if tc.ExpectedErrorString != "" && len(errors) > 0 {
				found := false
				for _, err := range errors {
					if err != nil && len(err.Error()) > 0 {
						if len(tc.ExpectedErrorString) <= len(err.Error()) &&
							err.Error()[:len(tc.ExpectedErrorString)] == tc.ExpectedErrorString {
							found = true
							break
						}
					}
				}
				if !found {
					t.Fatalf("expected error containing '%s', got errors: %v", tc.ExpectedErrorString, errors)
				}
			}
		})
	}
}
