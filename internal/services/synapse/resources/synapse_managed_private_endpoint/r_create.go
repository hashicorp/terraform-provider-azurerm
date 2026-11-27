// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_managed_private_endpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	workspaceClient := meta.(*clients.Client).Synapse.WorkspacesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	baseURI, err := NewSynapseWorkspaceBaseURI(meta.(*clients.Client).Account.Environment, workspaceId.WorkspaceName)
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, *workspaceId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", workspaceId, err)
	}

	if workspace.Model == nil || workspace.Model.Properties == nil || workspace.Model.Properties.ManagedVirtualNetwork == nil {
		return fmt.Errorf("empty or nil `ManagedVirtualNetwork` for %s", workspaceId)
	}
	virtualNetworkName := *workspace.Model.Properties.ManagedVirtualNetwork

	// TODO: migrate this to a data plane ID, will require a state migration
	id := parse.NewManagedPrivateEndpointID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, virtualNetworkName, d.Get("name").(string))

	dataPlaneID := managedprivateendpoints.NewManagedPrivateEndpointID(baseURI, id.ManagedVirtualNetworkName, id.Name)
	client := meta.(*clients.Client).Synapse.ManagedPrivateEndpointsClient.Clone(dataPlaneID.BaseURI)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, dataPlaneID)
		if !response.WasNotFound(existing.HttpResponse) {
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %w", id, err)
			}
			return tf.ImportAsExistsError("azurerm_synapse_managed_private_endpoint", id.ID())
		}
	}

	payload := managedprivateendpoints.ManagedPrivateEndpoint{
		Properties: &managedprivateendpoints.ManagedPrivateEndpointProperties{
			Fqdns:                 pluginsdk.ExpandStringSlice(d.Get("fully_qualified_domain_names").([]any)),
			GroupId:               pointer.To(d.Get("subresource_name").(string)),
			PrivateLinkResourceId: pointer.To(d.Get("target_resource_id").(string)),
		},
	}

	if _, err := client.Create(ctx, dataPlaneID, payload); err != nil {
		return fmt.Errorf("creating %s: %w", id, err)
	}

	d.SetId(id.ID())

	// The Create operation returns a 200 but behaves as if it is async, thus requiring polling on `provisioningState`
	pollerType := custompollers.NewSynapseManagedPrivateEndpointCreatePoller(client, dataPlaneID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling %s: %w", id, err)
	}

	return resourceSynapseManagedPrivateEndpointRead(d, meta)
}
