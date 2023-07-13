// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	managedvirtualnetwork "github.com/tombuildsstuff/kermit/sdk/synapse/2019-06-01-preview/synapse"
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
		},
	}
}

func resourceSynapseManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse workspace %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, err)
	}
	if workspace.WorkspaceProperties == nil || workspace.WorkspaceProperties.ManagedVirtualNetwork == nil {
		return fmt.Errorf("empty or nil `ManagedVirtualNetwork` for Synapse workspace %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, err)
	}
	virtualNetworkName := *workspace.WorkspaceProperties.ManagedVirtualNetwork

	id := parse.NewManagedPrivateEndpointID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, virtualNetworkName, d.Get("name").(string))
	client, err := synapseClient.ManagedPrivateEndpointsClient(workspaceId.Name, *synapseDomainSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %+v", id, err)
	}

	// check exist
	existing, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_synapse_managed_private_endpoint", id.ID())
	}

	managedPrivateEndpoint := managedvirtualnetwork.ManagedPrivateEndpoint{
		Properties: &managedvirtualnetwork.ManagedPrivateEndpointProperties{
			PrivateLinkResourceID: utils.String(d.Get("target_resource_id").(string)),
			GroupID:               utils.String(d.Get("subresource_name").(string)),
		},
	}
	if _, err := client.Create(ctx, id.ManagedVirtualNetworkName, id.Name, managedPrivateEndpoint); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSynapseManagedPrivateEndpointRead(d, meta)
}

func resourceSynapseManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %v", *id, err)
	}

	resp, err := client.Get(ctx, id.ManagedVirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID()
	d.Set("synapse_workspace_id", workspaceId)
	d.Set("name", id.Name)

	if props := resp.Properties; props != nil {
		d.Set("target_resource_id", props.PrivateLinkResourceID)
		d.Set("subresource_name", props.GroupID)
	}
	return nil
}

func resourceSynapseManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.ManagedPrivateEndpointsClient(id.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return fmt.Errorf("building Client for %s: %v", *id, err)
	}

	if _, err := client.Delete(ctx, id.ManagedVirtualNetworkName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
