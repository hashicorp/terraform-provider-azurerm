// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_build_pack_binding

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudBuildPackBindingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.BuildPackBindingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseBuildPackBindingIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.BuildPackBindingName)
	d.Set("spring_cloud_builder_id", appplatform_rm.NewBuilderID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.BuildServiceName, id.BuilderName).ID())
	if props := resp.Properties; props != nil {
		d.Set("binding_type", props.BindingType)
		if err := d.Set("launch", flattenBuildPackBindingBuildPackBindingLaunchProperties(props.LaunchProperties, d.Get("launch").([]any))); err != nil {
			return fmt.Errorf("setting `launch`: %+v", err)
		}
	}
	return nil
}

func flattenBuildPackBindingBuildPackBindingLaunchProperties(input *appplatform.BuildpackBindingLaunchProperties, old []any) []any {
	if input == nil {
		return make([]any, 0)
	}

	props := make(map[string]any)
	if input.Properties != nil {
		props = pluginsdk.FlattenMapStringPtrString(input.Properties)
	}
	secrets := make(map[string]any)
	if len(old) != 0 {
		v := old[0].(map[string]any)
		if secretsRaw, ok := v["secrets"]; ok && secretsRaw != nil {
			secrets = secretsRaw.(map[string]any)
		}
	}
	return []any{
		map[string]any{
			"properties": props,
			"secrets":    secrets,
		},
	}
}
