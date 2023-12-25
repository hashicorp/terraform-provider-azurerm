// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
)

func TestExtractListItem(t *testing.T) {
	tests := []struct {
		name         string
		args         string
		wantName     string
		wantOptional int
		wantContent  string
		wantEnums    []string
	}{
		{
			"", "", "", model.Default, "", nil,
		},
		{
			"store_name", "* `store_name` - (Required) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			"store_name", model.Required, "The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			[]string{"CertificateAuthority", "Root"},
		},
		{
			"store_name", "* `store_name` - (Optional) The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			"store_name", model.Optional, "The name of the Certificate. Possible values are `CertificateAuthority` and `Root`.",
			[]string{"CertificateAuthority", "Root"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := extractFieldFromLine(tt.args)
			if field.Name != tt.wantName {
				t.Errorf("extractFieldFromLine() gotName = %v, want %v", field.Name, tt.wantName)
			}
			if field.Required != tt.wantOptional {
				t.Errorf("extractFieldFromLine() gotOptional = %v, want %v", field.Required, tt.wantOptional)
			}
			if !reflect.DeepEqual(field.PossibleValues(), tt.wantEnums) {
				t.Errorf("extractFieldFromLine() gotEnums = %v, want %v", field.PossibleValues(), tt.wantEnums)
			}
		})
	}
}

func TestExtractBlockNames(t *testing.T) {
	var tests = []struct {
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
				t.Fatalf("test %d want: %v, got: A%v", idx, test.names, names)
			}
		})
	}
}

func TestScanOrSplit(t *testing.T) {
	content, _ := os.ReadFile(path.Join(ResourceDir(), "windows_function_app_slot.html.markdown"))
	content = bytes.TrimSuffix(content, []byte{'\n'})
	lines := strings.Split(string(content), "\n")
	sc := bufio.NewScanner(bytes.NewBuffer(content))
	var lines2 []string
	for sc.Scan() {
		lines2 = append(lines2, sc.Text())
	}
	if len(lines2) != len(lines) {
		t.Fatalf("%d: %d", len(lines), len(lines2))
	}
	for idx, line := range lines {
		if line != lines2[idx] {
			t.Fatalf("%d: %s: %s", idx, line, lines2[idx])
		}
	}
}

func TestDefaultValueReg(t *testing.T) {
	var lines = []string{
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
	for idx, line := range lines {
		val := getDefaultValue(line)
		t.Logf("%d idxs: %v", idx, val)
	}
}

func TestForceNewReg(t *testing.T) {
	// str := "* `address` - (Required) The list of upto 3 lines for address information. Changing this forces a new Databox Edge Order to be created.\n"
	str := "* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group to which this Virtual Machine should be assigned. Changing this forces a new resource to be created"
	res := ForceNewReg.MatchString(str)
	t.Log(res)
	str = ForceNewReg.ReplaceAllString(str, "")
	t.Log(str)
}
