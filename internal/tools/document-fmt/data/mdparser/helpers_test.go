// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdparser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestExtractFieldFromLine(t *testing.T) {
	tests := []struct {
		name         string
		args         string
		wantName     string
		wantRequired bool
		wantOptional bool
		wantContent  string
		wantEnums    []string
	}{
		{
			name:         "empty",
			args:         "",
			wantName:     "",
			wantRequired: false,
			wantOptional: false,
			wantContent:  "",
			wantEnums:    nil,
		},
		{
			name:         "required field with enums",
			args:         "* `store_name` - (Required) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			wantName:     "store_name",
			wantRequired: true,
			wantOptional: false,
			wantContent:  "* `store_name` - (Required) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			wantEnums:    []string{"CertificateAuthority", "Root"},
		},
		{
			name:         "optional field with enums",
			args:         "* `store_name` - (Optional) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			wantName:     "store_name",
			wantRequired: false,
			wantOptional: true,
			wantContent:  "* `store_name` - (Optional) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			wantEnums:    []string{"CertificateAuthority", "Root"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := extractFieldFromLine(tt.args)
			if field.Name != tt.wantName {
				t.Errorf("extractFieldFromLine() gotName = %v, want %v", field.Name, tt.wantName)
			}
			if field.Required != tt.wantRequired {
				t.Errorf("extractFieldFromLine() gotRequired = %v, want %v", field.Required, tt.wantRequired)
			}
			if field.Optional != tt.wantOptional {
				t.Errorf("extractFieldFromLine() gotOptional = %v, want %v", field.Optional, tt.wantOptional)
			}
			if !reflect.DeepEqual(field.PossibleValues(), tt.wantEnums) {
				t.Errorf("extractFieldFromLine() gotEnums = %v, want %v", field.PossibleValues(), tt.wantEnums)
			}
		})
	}
}

func TestExtractBlockNames(t *testing.T) {
	tests := []struct {
		line  string
		names []string
	}{
		{
			"A `management`, `portal`, `developer_portal` and `scm` block supports the following:",
			[]string{"management", "portal", "developer_portal", "scm"},
		},
		{
			"An `identity` block supports the following:",
			[]string{"identity"},
		},
		{
			"A `policy` block supports the following:",
			[]string{"policy"},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			names := extractBlockNames(test.line)
			if !reflect.DeepEqual(names, test.names) {
				t.Fatalf("test %d want: %v, got: %v", idx, test.names, names)
			}
		})
	}
}

func TestGetDefaultValue(t *testing.T) {
	lines := []string{
		"* `load_balancing_mode` - (Optional) The Site load balancing. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.",
		"* `local_mysql_enabled` - (Optional) Use Local MySQL. Defaults to `false`.",
		"* `local_mysql_enabled` - (Optional) Use Local MySQL. Defaults to `\"\"`.",
		"* `minimum_tls_version` - (Optional) The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
		"* `export_policy_enabled` - (Optional) Boolean value that indicates whether export policy is enabled. Defaults to `true`. In order to set it to `false`, make sure the `public_network_access_enabled` is also set to `false`.\n",
		"* `probe_threshold` - (Optional) The number of consecutive successful or failed probes that allow or deny traffic to this endpoint. Possible values range from `1` to `100`. The default value is `1`.\n",
	}
	values := []string{
		"LeastRequests",
		"false",
		`""`,
		"1.2",
		"true",
		"1",
	}
	for idx, line := range lines {
		val := getDefaultValue(line)
		if values[idx] != val {
			t.Fatalf("idx %d want: %s, got: %v", idx, values[idx], val)
		}
	}
}

func TestIsForceNew(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{
			name: "force new resource",
			line: "* `address` - (Required) The list of upto 3 lines for address information. Changing this forces a new Databox Edge Order to be created.",
			want: true,
		},
		{
			name: "force new without period",
			line: "* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group to which this Virtual Machine should be assigned. Changing this forces a new resource to be created",
			want: true,
		},
		{
			name: "partial force new (should be false)",
			line: "* `field` - (Optional) Description. Changing this forces a new resource created when something happens.",
			want: false,
		},
		{
			name: "no force new",
			line: "* `field` - (Optional) Description of the field.",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isForceNew(tt.line); got != tt.want {
				t.Errorf("isForceNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPossibleValues(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantEnums []string
	}{
		{
			name:      "possible values are",
			line:      "* `store_name` - (Required) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			wantEnums: []string{"CertificateAuthority", "Root"},
		},
		{
			name:      "possible values include",
			line:      "* `load_balancing_mode` - (Optional) The Site load balancing. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.",
			wantEnums: []string{"WeightedRoundRobin", "LeastRequests", "LeastResponseTime", "WeightedTotalTraffic", "RequestHash", "PerSiteRoundRobin"},
		},
		{
			name:      "possible values with version numbers",
			line:      "* `minimum_tls_version` - (Optional) The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
			wantEnums: []string{"1.0", "1.1", "1.2"},
		},
		{
			name:      "possible values with quotes",
			line:      "* `sku` - (Required) The SKU of the resource. Possible values are `Basic`, `Standard`, and `Premium`.",
			wantEnums: []string{"Basic", "Standard", "Premium"},
		},
		{
			name:      "must be one of",
			line:      "* `protocol` - (Required) The protocol to use. Must be one of `TCP`, `UDP`, or `ICMP`.",
			wantEnums: []string{"TCP", "UDP", "ICMP"},
		},
		{
			name:      "valid values",
			line:      "* `type` - (Required) The type of the resource. Valid values are `TypeA`, `TypeB`, and `TypeC`.",
			wantEnums: []string{"TypeA", "TypeB", "TypeC"},
		},
		{
			name:      "allowed values",
			line:      "* `tier` - (Optional) The tier of the service. Allowed values include `Free`, `Shared`, `Basic`, `Standard`, and `Premium`.",
			wantEnums: []string{"Free", "Shared", "Basic", "Standard", "Premium"},
		},
		{
			name:      "no possible values",
			line:      "* `name` - (Required) The name of the resource.",
			wantEnums: nil,
		},
		{
			name:      "possible values range (extracts range bounds)",
			line:      "* `probe_threshold` - (Optional) The number of consecutive successful or failed probes. Possible values range from `1` to `100`. The default value is `1`.",
			wantEnums: []string{"1", "100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := extractFieldFromLine(tt.line)
			if !reflect.DeepEqual(field.PossibleValues(), tt.wantEnums) {
				t.Errorf("extractFieldFromLine() PossibleValues = %v, want %v", field.PossibleValues(), tt.wantEnums)
			}
		})
	}
}

func TestParseErrorsInFieldExtraction(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		wantError     bool
		expectedError string
	}{
		{
			name:          "missing field name",
			line:          "* - (Required) This is a description without a field name.",
			wantError:     true,
			expectedError: NoFieldNameFound,
		},
		{
			name:          "missing field name",
			line:          "*`fieldname`- (Required) This is a description without a field name.",
			wantError:     true,
			expectedError: NoFieldNameFound,
		},
		{
			name:      "valid field",
			line:      "* `name` - (Required) The name of the resource.",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := extractFieldFromLine(tt.line)

			if tt.wantError {
				if len(field.ParseErrors) == 0 {
					t.Errorf("expected parse error but got none")
				}
				found := false
				for _, err := range field.ParseErrors {
					if strings.Contains(err, tt.expectedError) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error %q but got %v", tt.expectedError, field.ParseErrors)
				}
			} else if len(field.ParseErrors) > 0 {
				t.Errorf("unexpected parse errors: %v", field.ParseErrors)
			}
		})
	}
}

func TestBlockPropertyDetection(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantBlock bool
	}{
		{
			name:      "block with 'as defined below'",
			line:      "* `config` - (Required) A `config` block as defined below.",
			wantBlock: true,
		},
		{
			name:      "block with 'as detailed below'",
			line:      "* `settings` - (Optional) One or more `settings` blocks as detailed below.",
			wantBlock: true,
		},
		{
			name:      "not a block",
			line:      "* `name` - (Required) The name of the resource.",
			wantBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isBlock := guessBlockProperty(tt.line)
			if isBlock != tt.wantBlock {
				t.Errorf("guessBlockProperty() = %v, want %v", isBlock, tt.wantBlock)
			}
		})
	}
}
