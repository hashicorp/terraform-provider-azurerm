// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_managed_private_endpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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

	if _, err := client.Delete(ctx, dataPlaneID); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// The delete operation returns immediately but this doesn't necessarily mean the resource is gone
	pollerType := custompollers.NewSynapseManagedPrivateEndpointDeletePoller(client, dataPlaneID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling %s: %w", id, err)
	}

	return nil
}
