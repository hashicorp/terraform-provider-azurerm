// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateNetAppVolumeProtocolConversion(t *testing.T) {
	cases := []struct {
		Name                string
		OldProtocols        []string
		NewProtocols        []string
		KerberosEnabled     bool
		DataReplication     []interface{}
		ExportPolicyRules   []interface{}
		ExpectedErrors      int
		ExpectedErrorString string
	}{
		{
			Name:              "ValidConversionNFSv3ToNFSv41",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv4.1"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionNFSv41ToNFSv3",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv3"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "NoChangeNoValidation",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv3"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "InitialCreationNoValidation",
			OldProtocols:      []string{},
			NewProtocols:      []string{"NFSv3"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:                "InvalidKerberosNFSv41ToNFSv3Conversion",
			OldProtocols:        []string{"NFSv4.1"},
			NewProtocols:        []string{"NFSv3"},
			KerberosEnabled:     true,
			DataReplication:     []interface{}{},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      1,
			ExpectedErrorString: "cannot convert an NFSv4.1 volume with Kerberos enabled to NFSv3",
		},
		{
			Name:                "InvalidDualProtocolConversionOld",
			OldProtocols:        []string{"NFSv3", "CIFS"},
			NewProtocols:        []string{"NFSv4.1"},
			KerberosEnabled:     false,
			DataReplication:     []interface{}{},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      2, // Both dual-protocol and CIFS conversion errors
			ExpectedErrorString: "cannot change the NFS version of a dual-protocol volume",
		},
		{
			Name:                "InvalidDualProtocolConversionNew",
			OldProtocols:        []string{"NFSv3"},
			NewProtocols:        []string{"NFSv4.1", "CIFS"},
			KerberosEnabled:     false,
			DataReplication:     []interface{}{},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      2, // Both dual-protocol and CIFS conversion errors
			ExpectedErrorString: "cannot change the NFS version of a dual-protocol volume",
		},
		{
			Name:              "InvalidCIFSConversionOld",
			OldProtocols:      []string{"CIFS"},
			NewProtocols:      []string{"NFSv3"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0, // No NFS protocol change detected, no validation
		},
		{
			Name:              "InvalidCIFSConversionNew",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"CIFS"},
			KerberosEnabled:   false,
			DataReplication:   []interface{}{},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0, // No NFS protocol change detected, no validation
		},
		{
			Name:                "InvalidCIFSToDualProtocol",
			OldProtocols:        []string{"NFSv3"},
			NewProtocols:        []string{"NFSv4.1", "CIFS"},
			KerberosEnabled:     false,
			DataReplication:     []interface{}{},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      2, // Both dual-protocol and CIFS conversion errors
			ExpectedErrorString: "cannot change the NFS version of a dual-protocol volume",
		},
		{
			Name:            "InvalidDestinationVolumeConversion",
			OldProtocols:    []string{"NFSv3"},
			NewProtocols:    []string{"NFSv4.1"},
			KerberosEnabled: false,
			DataReplication: []interface{}{
				map[string]interface{}{
					"endpoint_type": "dst",
				},
			},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      1,
			ExpectedErrorString: "cannot convert a destination volume in a cross-region replication relationship",
		},
		{
			Name:            "ValidSourceVolumeConversion",
			OldProtocols:    []string{"NFSv3"},
			NewProtocols:    []string{"NFSv4.1"},
			KerberosEnabled: false,
			DataReplication: []interface{}{
				map[string]interface{}{
					"endpoint_type": "src",
				},
			},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:            "ValidExportPolicyDuringProtocolConversion",
			OldProtocols:    []string{"NFSv3"},
			NewProtocols:    []string{"NFSv4.1"},
			KerberosEnabled: false,
			DataReplication: []interface{}{},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocol": []interface{}{"NFSv4.1"}, // Will be updated to match new volume protocol
				},
			},
			ExpectedErrors: 0, // Export policy validation is skipped during protocol conversion
		},
		{
			Name:            "ValidExportPolicyWithProtocolsEnabledDuringConversion",
			OldProtocols:    []string{"NFSv4.1"},
			NewProtocols:    []string{"NFSv3"},
			KerberosEnabled: false,
			DataReplication: []interface{}{},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocols_enabled": []interface{}{"NFSv3"}, // Will be updated to match new volume protocol
				},
			},
			ExpectedErrors: 0, // Export policy validation is skipped during protocol conversion
		},
		{
			Name:            "ValidExportPolicyProtocolMatchWithProtocolField",
			OldProtocols:    []string{"NFSv3"},
			NewProtocols:    []string{"NFSv4.1"},
			KerberosEnabled: false,
			DataReplication: []interface{}{},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocol": []interface{}{"NFSv4.1"}, // Matches volume protocol NFSv4.1
				},
			},
			ExpectedErrors: 0,
		},
		{
			Name:            "ValidExportPolicyProtocolMatchWithProtocolsEnabledField",
			OldProtocols:    []string{"NFSv4.1"},
			NewProtocols:    []string{"NFSv3"},
			KerberosEnabled: false,
			DataReplication: []interface{}{},
			ExportPolicyRules: []interface{}{
				map[string]interface{}{
					"protocols_enabled": []interface{}{"NFSv3"}, // Matches volume protocol NFSv3
				},
			},
			ExpectedErrors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeProtocolConversion(tc.OldProtocols, tc.NewProtocols, tc.KerberosEnabled, tc.DataReplication, tc.ExportPolicyRules)

			if len(errors) != tc.ExpectedErrors {
				t.Fatalf("expected ValidateNetAppVolumeProtocolConversion to return %d error(s), got %d: %v", tc.ExpectedErrors, len(errors), errors)
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
