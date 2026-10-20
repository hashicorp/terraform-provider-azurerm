// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_scope_map

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryScopeMapUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.ScopeMaps
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scopemaps.ParseScopeMapID(d.Id())
	if err != nil {
		return err
	}

	parameters := scopemaps.ScopeMapUpdateParameters{
		Properties: &scopemaps.ScopeMapPropertiesUpdateParameters{
			Description: pointer.To(d.Get("description").(string)),
			Actions:     pluginsdk.ExpandStringSlice(d.Get("actions").([]any)),
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryScopeMapRead(d, meta)
}
