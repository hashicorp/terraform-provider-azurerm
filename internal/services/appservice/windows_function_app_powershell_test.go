// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	appserviceclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app_slot"
)

func TestWindowsFunctionAppPowerShellWorkerCustomizeDiff(t *testing.T) {
	cases := []struct {
		name    string
		version cty.Value
		worker  cty.Value
		wantErr bool
	}{
		{name: "7.6 explicit true", version: cty.StringVal("7.6"), worker: cty.True, wantErr: true},
		{name: "7.6 unset", version: cty.StringVal("7.6"), worker: cty.NullVal(cty.Bool)},
		{name: "7.6 explicit false", version: cty.StringVal("7.6"), worker: cty.False},
		{name: "7.6 unknown worker", version: cty.StringVal("7.6"), worker: cty.UnknownVal(cty.Bool)},
		{name: "unknown version", version: cty.UnknownVal(cty.String), worker: cty.True},
		{name: "7.4 explicit true", version: cty.StringVal("7.4"), worker: cty.True},
		{name: "dotnet explicit true", version: cty.NullVal(cty.String), worker: cty.True},
	}
	resources := map[string]struct {
		siteConfig *schema.Schema
		customize  sdk.ResourceFunc
	}{
		"app":  {siteConfig: helpers.SiteConfigSchemaWindowsFunctionApp(), customize: windows_function_app.Resource{}.CustomizeDiff()},
		"slot": {siteConfig: windows_function_app_slot.SiteConfigSchemaWindowsFunctionAppSlot(), customize: windows_function_app_slot.Resource{}.CustomizeDiff()},
	}
	for kind, resource := range resources {
		t.Run(kind, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					raw := cty.ObjectVal(map[string]cty.Value{
						"site_config": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
							"use_32_bit_worker": tc.worker,
							"application_stack": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
								"powershell_core_version": tc.version,
							})}),
						})}),
					})
					r := &schema.Resource{
						Schema: map[string]*schema.Schema{"site_config": resource.siteConfig},
						CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, _ interface{}) error {
							return resource.customize.Func(ctx, sdk.ResourceMetaData{
								ResourceDiff: d,
								Client:       &clients.Client{AppService: &appserviceclient.Client{}},
							})
						},
					}
					// Resource.Diff obtains the protocol raw config from InstanceState.RawConfig.
					state := &terraform.InstanceState{
						ID:        "test",
						RawConfig: raw,
						Attributes: map[string]string{
							"site_config.#":                     "1",
							"site_config.0.application_stack.#": "1",
							"site_config.0.application_stack.0.powershell_core_version": "7.6",
							"site_config.0.application_stack.0.dotnet_version":          "v4.0",
							"site_config.0.use_32_bit_worker":                           "false",
						},
					}
					config := terraform.NewResourceConfigRaw(map[string]interface{}{
						"site_config": []interface{}{map[string]interface{}{
							"application_stack": []interface{}{map[string]interface{}{"powershell_core_version": "7.6"}},
						}},
					})
					_, err := r.Diff(context.Background(), state, config, nil)
					if tc.wantErr {
						const message = "PowerShell 7.6 requires a 64-bit worker; remove `use_32_bit_worker` or set it to `false`"
						if err == nil || !strings.Contains(err.Error(), message) {
							t.Fatalf("expected %q, got %v", message, err)
						}
					} else if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
				})
			}
		})
	}
}
