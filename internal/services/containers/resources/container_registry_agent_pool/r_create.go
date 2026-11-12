// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_agent_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryAgentPoolCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2019_06_01_preview.AgentPools
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := agentpools.NewAgentPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_agent_pool", id.ID())
		}
	}

	parameters := agentpools.AgentPool{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &agentpools.AgentPoolProperties{
			// @favoretti: Only Linux is supported
			Os:    pointer.To(agentpools.OSLinux),
			Count: pointer.To(int64(d.Get("instance_count").(int))),
			Tier:  pointer.To(d.Get("tier").(string)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("virtual_network_subnet_id"); ok {
		parameters.Properties.VirtualNetworkSubnetResourceId = pointer.To(v.(string))
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryAgentPoolRead(d, meta)
}
