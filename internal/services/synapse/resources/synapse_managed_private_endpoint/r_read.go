// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_managed_private_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	baseURI, err := NewSynapseWorkspaceBaseURI(meta.(*clients.Client).Account.Environment, id.WorkspaceName)
	if err != nil {
		return err
	}

	dataPlaneID := managedprivateendpoints.NewManagedPrivateEndpointID(baseURI, id.ManagedVirtualNetworkName, id.Name)
	client := meta.(*clients.Client).Synapse.ManagedPrivateEndpointsClient.Clone(dataPlaneID.BaseURI)

	resp, err := client.Get(ctx, dataPlaneID)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("synapse_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("name", id.Name)

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("fully_qualified_domain_names", pointer.From(props.Fqdns))
			d.Set("subresource_name", props.GroupId)
			d.Set("target_resource_id", props.PrivateLinkResourceId)
		}
	}
	return nil
}
