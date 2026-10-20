// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_token

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryTokenRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Tokens
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tokens.ParseTokenID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.TokenName)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			status := true
			if v := props.Status; v != nil && *v == tokens.TokenStatusDisabled {
				status = false
			}
			d.Set("enabled", status)

			d.Set("scope_map_id", pointer.From(props.ScopeMapId))
		}
	}

	return nil
}
