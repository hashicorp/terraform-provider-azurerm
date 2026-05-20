// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseManagedPrivateEndpointCreate,
		Read:   resourceSynapseManagedPrivateEndpointRead,
		Delete: resourceSynapseManagedPrivateEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"subresource_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},

			"fully_qualified_domain_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceSynapseManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	existing, err := client.Get(ctx, dataPlaneID)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %w", id, err)
		}
		return tf.ImportAsExistsError("azurerm_synapse_managed_private_endpoint", id.ID())
	}

	payload := managedprivateendpoints.ManagedPrivateEndpoint{
		Properties: &managedprivateendpoints.ManagedPrivateEndpointProperties{
			Fqdns:                 utils.ExpandStringSlice(d.Get("fully_qualified_domain_names").([]any)),
			GroupId:               pointer.To(d.Get("subresource_name").(string)),
			PrivateLinkResourceId: pointer.To(d.Get("target_resource_id").(string)),
		},
	}

	if _, err := client.Create(ctx, dataPlaneID, payload); err != nil {
		return fmt.Errorf("creating %s: %w", id, err)
	}

	// The Create operation returns a 200 but behaves as if it is async, thus requiring polling on `provisioningState`
	pollerType := custompollers.NewSynapseManagedPrivateEndpointCreatePoller(client, dataPlaneID)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling %s: %w", id, err)
	}

	d.SetId(id.ID())
	return resourceSynapseManagedPrivateEndpointRead(d, meta)
}

func resourceSynapseManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceSynapseManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
