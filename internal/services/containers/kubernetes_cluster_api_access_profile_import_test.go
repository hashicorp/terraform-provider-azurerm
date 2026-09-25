// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesClusterAPIAccessProfileRetainedEmptyImport(t *testing.T) {
	resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
		"api_server_access_profile": resourceKubernetesCluster().Schema["api_server_access_profile"],
	}}
	profiles := map[string]*managedclusters.ManagedClusterAPIServerAccessProfile{
		"nil profile":  nil,
		"nil fields":   {},
		"empty ranges": {AuthorizedIPRanges: pointer.To([]string{})},
		"API defaults": {AuthorizedIPRanges: pointer.To([]string{}), EnableVnetIntegration: pointer.To(false), SubnetId: pointer.To("")},
	}
	for name, profile := range profiles {
		t.Run(name, func(t *testing.T) {
			for _, ranges := range []string{"empty", "null", "omitted"} {
				t.Run(ranges, func(t *testing.T) {
					block := map[string]interface{}{}
					switch ranges {
					case "empty":
						block["authorized_ip_ranges"] = []interface{}{}
					case "null":
						block["authorized_ip_ranges"] = nil
					}
					config := map[string]interface{}{"api_server_access_profile": []interface{}{block}}
					imported := schema.TestResourceDataRaw(t, resource.Schema, nil)
					imported.SetId("existing-cluster")
					if err := imported.Set("api_server_access_profile", flattenKubernetesClusterAPIAccessProfile(profile, false)); err != nil {
						t.Fatal(err)
					}
					if got := imported.Get("api_server_access_profile").([]interface{}); len(got) != 0 {
						t.Fatalf("import must not infer an empty block from API defaults: %#v", got)
					}

					diff, err := resource.SimpleDiff(context.Background(), imported.State(), terraform.NewResourceConfigRaw(config), nil)
					if err != nil {
						t.Fatal(err)
					}
					if diff.Empty() || diff.RequiresNew() {
						t.Fatalf("expected an in-place reconciliation of the configured empty block after import, got %#v", diff)
					}
					resource.UpdateContext = func(_ context.Context, d *pluginsdk.ResourceData, _ interface{}) diag.Diagnostics {
						if err := d.Set("api_server_access_profile", flattenKubernetesClusterAPIAccessProfile(profile, len(d.Get("api_server_access_profile").([]interface{})) > 0)); err != nil {
							return diag.FromErr(err)
						}
						return nil
					}
					state, diagnostics := resource.Apply(context.Background(), imported.State(), diff, nil)
					if diagnostics.HasError() {
						t.Fatal(diagnostics)
					}
					diff, err = resource.SimpleDiff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
					if err != nil {
						t.Fatal(err)
					}
					if !diff.Empty() {
						t.Fatalf("reconciled empty block produced a follow-up diff: %#v", diff.Attributes)
					}
				})
			}
		})
	}
}
