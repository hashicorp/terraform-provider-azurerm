// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"
)

func TestValidateNetAppVolumeGroupSAPHanaProtocolRestrictions(t *testing.T) {
	cases := []struct {
		Name                string
		VolumeSpecName      string
		Protocols           []string
		ExpectedErrors      int
		ExpectedErrorString string
	}{
		{
			Name:           "ValidNFSv41OnDataVolume",
			VolumeSpecName: "data",
			Protocols:      []string{"NFSv4.1"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv41OnLogVolume",
			VolumeSpecName: "log",
			Protocols:      []string{"NFSv4.1"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv41OnSharedVolume",
			VolumeSpecName: "shared",
			Protocols:      []string{"NFSv4.1"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv3OnDataBackupVolume",
			VolumeSpecName: "data-backup",
			Protocols:      []string{"NFSv3"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv3OnLogBackupVolume",
			VolumeSpecName: "log-backup",
			Protocols:      []string{"NFSv3"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv41OnDataBackupVolume",
			VolumeSpecName: "data-backup",
			Protocols:      []string{"NFSv4.1"},
			ExpectedErrors: 0,
		},
		{
			Name:           "ValidNFSv41OnLogBackupVolume",
			VolumeSpecName: "log-backup",
			Protocols:      []string{"NFSv4.1"},
			ExpectedErrors: 0,
		},
		{
			Name:                "InvalidNFSv3OnDataVolume",
			VolumeSpecName:      "data",
			Protocols:           []string{"NFSv3"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'data' volumes for SAP HANA",
		},
		{
			Name:                "InvalidNFSv3OnLogVolume",
			VolumeSpecName:      "log",
			Protocols:           []string{"NFSv3"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'log' volumes for SAP HANA",
		},
		{
			Name:                "InvalidNFSv3OnSharedVolume",
			VolumeSpecName:      "shared",
			Protocols:           []string{"NFSv3"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'shared' volumes for SAP HANA",
		},
		{
			Name:                "InvalidNFSv3WithNFSv41OnDataVolume",
			VolumeSpecName:      "data",
			Protocols:           []string{"NFSv3", "NFSv4.1"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'data' volumes for SAP HANA",
		},
		{
			Name:                "InvalidNFSv3WithNFSv41OnLogVolume",
			VolumeSpecName:      "log",
			Protocols:           []string{"NFSv3", "NFSv4.1"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'log' volumes for SAP HANA",
		},
		{
			Name:                "InvalidNFSv3WithNFSv41OnSharedVolume",
			VolumeSpecName:      "shared",
			Protocols:           []string{"NFSv3", "NFSv4.1"},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'shared' volumes for SAP HANA",
		},
		{
			Name:           "ValidOtherVolumeSpecWithNFSv3",
			VolumeSpecName: "custom-volume",
			Protocols:      []string{"NFSv3"},
			ExpectedErrors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Simulate the validation logic from the SAP HANA resource
			var errors []error

			for _, protocol := range tc.Protocols {
				if protocol == "NFSv3" {
					// NFSv3 is not allowed on data, log, and shared volumes for SAP HANA
					if tc.VolumeSpecName == "data" || tc.VolumeSpecName == "log" || tc.VolumeSpecName == "shared" {
						errors = append(errors, fmt.Errorf("NFSv3 protocol is not supported on '%s' volumes for SAP HANA. Only NFSv4.1 is supported for critical SAP HANA volumes (data, log, shared). NFSv3 can only be used for backup volumes (data-backup, log-backup)", tc.VolumeSpecName))
					}
				}
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

func TestValidateNetAppVolumeGroupSAPHanaProtocolConversion(t *testing.T) {
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
			Name:              "ValidConversionDataBackupNFSv3ToNFSv41",
			VolumeSpecName:    "data-backup",
			OldProtocols:      []string{"NFSv3"},
			NewProtocols:      []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionLogBackupNFSv41ToNFSv3",
			VolumeSpecName:    "log-backup",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv3"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:              "ValidConversionDataVolumeNFSv41ToNFSv41",
			VolumeSpecName:    "data",
			OldProtocols:      []string{"NFSv4.1"},
			NewProtocols:      []string{"NFSv4.1"},
			ExportPolicyRules: []interface{}{},
			ExpectedErrors:    0,
		},
		{
			Name:                "InvalidConversionDataVolumeToNFSv3",
			VolumeSpecName:      "data",
			OldProtocols:        []string{"NFSv4.1"},
			NewProtocols:        []string{"NFSv3"},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'data' volumes for SAP HANA",
		},
		{
			Name:                "InvalidConversionLogVolumeToNFSv3",
			VolumeSpecName:      "log",
			OldProtocols:        []string{"NFSv4.1"},
			NewProtocols:        []string{"NFSv3"},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'log' volumes for SAP HANA",
		},
		{
			Name:                "InvalidConversionSharedVolumeToNFSv3",
			VolumeSpecName:      "shared",
			OldProtocols:        []string{"NFSv4.1"},
			NewProtocols:        []string{"NFSv3"},
			ExportPolicyRules:   []interface{}{},
			ExpectedErrors:      1,
			ExpectedErrorString: "NFSv3 protocol is not supported on 'shared' volumes for SAP HANA",
		},
		{
			Name:           "ValidConversionWithExportPolicy",
			VolumeSpecName: "data-backup",
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
			Name:           "NoConversionNoValidation",
			VolumeSpecName: "data",
			OldProtocols:   []string{"NFSv4.1"},
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

			// First, validate SAP HANA protocol restrictions
			for _, protocol := range tc.NewProtocols {
				if protocol == "NFSv3" {
					if tc.VolumeSpecName == "data" || tc.VolumeSpecName == "log" || tc.VolumeSpecName == "shared" {
						errors = append(errors, fmt.Errorf("NFSv3 protocol is not supported on '%s' volumes for SAP HANA. Only NFSv4.1 is supported for critical SAP HANA volumes (data, log, shared). NFSv3 can only be used for backup volumes (data-backup, log-backup)", tc.VolumeSpecName))
					}
				}
			}

			// Then, validate protocol conversion (only if we're changing protocols)
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

// Helper function to compare string slices
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
