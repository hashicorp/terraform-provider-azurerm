package pluginsdk

import (
	"fmt"
	"strings"
	"testing"
)

func TestSnakeCase(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			Input:  "ResourceGroupName",
			Output: "resource_group_name",
		},
		{
			Input:  "ServerGroupsv2Name",
			Output: "server_groups_v2_name",
		},
		{
			Input:  "ServerGroupsV2Name",
			Output: "server_groups_v2_name",
		},
		{
			Input:  "ServerGroupsV42Name",
			Output: "server_groups_v42_name",
		},
		{
			Input:  "ServerGroupsV2",
			Output: "server_groups_v2",
		},
		{
			Input:  "Terraform404",
			Output: "terraform_404",
		},
		{
			Input:  "1TerraformProvider",
			Output: "1_terraform_provider",
		},
		{
			Input:  "TFProvider",
			Output: "tf_provider",
		},
		{
			Input:  "Peer2Peer",
			Output: "peer_2_peer",
		},
		{
			Input:  "V2_Terraform",
			Output: "v2_terraform",
		},
		{
			Input:  "vV2_Terraform",
			Output: "v_v2_terraform",
		},
		{
			Input:  "VV2_Terraform",
			Output: "v_v2_terraform",
		},
	}

	failures := make([]string, 0)
	for _, tc := range cases {
		if v := toSnakeCase(tc.Input); v != tc.Output {
			failures = append(failures, fmt.Sprintf("expected %s, got %s", tc.Output, v))
		}
	}
	if len(failures) > 0 {
		t.Fatal(strings.Join(failures, "\n"))
	}
}
