// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_scope_map

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	containersClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryScopeMapDelete(d *pluginsdk.ResourceData, meta any) error {
	client := containersClient.From(meta).ContainerRegistryClient.ScopeMaps
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := scopemaps.ParseScopeMapID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
