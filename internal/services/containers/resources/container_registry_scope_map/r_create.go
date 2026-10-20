// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_scope_map

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryScopeMapCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.ScopeMaps
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scopemaps.NewScopeMapID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_scope_map", id.ID())
		}
	}

	parameters := scopemaps.ScopeMap{
		Properties: &scopemaps.ScopeMapProperties{
			Description: pointer.To(d.Get("description").(string)),
			Actions:     pointer.From(pluginsdk.ExpandStringSlice(d.Get("actions").([]any))),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryScopeMapRead(d, meta)
}
