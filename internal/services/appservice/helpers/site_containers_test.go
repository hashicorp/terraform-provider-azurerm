// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import "testing"

func TestValidateSiteContainerEnvironmentVariables(t *testing.T) {
	cases := []struct {
		name    string
		input   []SiteContainerEnvironmentVariable
		wantErr bool
	}{
		{
			name:    "empty",
			input:   nil,
			wantErr: false,
		},
		{
			name: "unique names",
			input: []SiteContainerEnvironmentVariable{
				{Name: "FOO", AppSettingName: "FOO_SETTING"},
				{Name: "BAR", AppSettingName: "BAR_SETTING"},
			},
			wantErr: false,
		},
		{
			name: "duplicate name different app setting",
			input: []SiteContainerEnvironmentVariable{
				{Name: "FOO", AppSettingName: "FOO_SETTING"},
				{Name: "FOO", AppSettingName: "FOO_SETTING_ALT"},
			},
			wantErr: true,
		},
		{
			name: "duplicate name identical app setting",
			input: []SiteContainerEnvironmentVariable{
				{Name: "FOO", AppSettingName: "FOO_SETTING"},
				{Name: "FOO", AppSettingName: "FOO_SETTING"},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSiteContainerEnvironmentVariables(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
		})
	}
}
