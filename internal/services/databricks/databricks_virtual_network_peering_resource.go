// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/vnetpeering"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const databricksVnetPeeringsResourceType string = "azurerm_databricks_virtual_network_peering"

func resourceDatabricksVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabricksVirtualNetworkPeeringCreate,
		Read:   resourceDatabricksVirtualNetworkPeeringRead,
		Update: resourceDatabricksVirtualNetworkPeeringUpdate,
		Delete: resourceDatabricksVirtualNetworkPeeringDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := vnetpeering.ParseVirtualNetworkPeeringID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabricksVirtualNetworkPeeringName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"remote_address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.CIDRIsIPv4OrIPv6,
				},
			},

			"remote_virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"allow_virtual_network_access": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allow_forwarded_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allow_gateway_transit": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"use_remote_gateways": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDatabricksVirtualNetworkPeeringCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM databricks virtual network peering creation.")
	var id vnetpeering.VirtualNetworkPeeringId

	// I need to include the workspace ID in the properties because I need the name
	// of the workspace to create the peerings ID
	workspaceId, err := vnetpeering.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse 'workspace_id': %+v", err)
	}

	id = vnetpeering.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), workspaceId.WorkspaceName, d.Get("name").(string))

	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Databricks %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_databricks_virtual_network_peering", id.ID())
	}

	allowForwardedTraffic := d.Get("allow_forwarded_traffic").(bool)
	allowGatewayTransit := d.Get("allow_gateway_transit").(bool)
	allowVirtualNetworkAccess := d.Get("allow_virtual_network_access").(bool)
	useRemoteGateways := d.Get("use_remote_gateways").(bool)
	remoteVirtualNetwork := d.Get("remote_virtual_network_id").(string)
	databricksAddressSpace := utils.ExpandStringSlice(d.Get("address_space_prefixes").([]interface{}))
	remoteAddressSpace := utils.ExpandStringSlice(d.Get("remote_address_space_prefixes").([]interface{}))

	props := vnetpeering.VirtualNetworkPeeringPropertiesFormat{
		DatabricksAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: databricksAddressSpace,
		},
		// The RP always creates the same vNet ID for the Databricks internal vNet in the below format:
		// '/subscriptions/{subscription}/resourceGroups/{group1}/providers/Microsoft.Network/virtualNetworks/workers-vnet'
		DatabricksVirtualNetwork: &vnetpeering.VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork{
			Id: utils.String(commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, "workers-vnet").ID()),
		},
		RemoteAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: remoteAddressSpace,
		},
		RemoteVirtualNetwork: vnetpeering.VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork{
			Id: utils.String(remoteVirtualNetwork),
		},
		AllowForwardedTraffic:     &allowForwardedTraffic,
		AllowGatewayTransit:       &allowGatewayTransit,
		AllowVirtualNetworkAccess: &allowVirtualNetworkAccess,
		UseRemoteGateways:         &useRemoteGateways,
	}

	peer := vnetpeering.VirtualNetworkPeering{
		Name:       &id.VirtualNetworkPeeringName,
		Properties: props,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, peer); err != nil {
		return fmt.Errorf("creating Databricks %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDatabricksVirtualNetworkPeeringRead(d, meta)
}

func resourceDatabricksVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("workspace_id", vnetpeering.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		d.Set("allow_virtual_network_access", model.Properties.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", model.Properties.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", model.Properties.AllowGatewayTransit)
		d.Set("use_remote_gateways", model.Properties.UseRemoteGateways)

		addressSpacePrefixes := make([]string, 0)
		if model.Properties.DatabricksAddressSpace != nil && model.Properties.DatabricksAddressSpace.AddressPrefixes != nil {
			addressSpacePrefixes = *model.Properties.DatabricksAddressSpace.AddressPrefixes
		}
		d.Set("address_space_prefixes", addressSpacePrefixes)

		remoteAddressSpacePrefixes := make([]string, 0)
		if model.Properties.RemoteAddressSpace != nil && model.Properties.RemoteAddressSpace.AddressPrefixes != nil {
			remoteAddressSpacePrefixes = *model.Properties.RemoteAddressSpace.AddressPrefixes
		}
		d.Set("remote_address_space_prefixes", remoteAddressSpacePrefixes)

		databricksVirtualNetworkId := ""
		if model.Properties.DatabricksVirtualNetwork != nil && model.Properties.DatabricksVirtualNetwork.Id != nil {
			databricksVirtualNetworkId = *model.Properties.DatabricksVirtualNetwork.Id
		}
		d.Set("virtual_network_id", databricksVirtualNetworkId)

		remoteVirtualNetworkId := ""
		if model.Properties.RemoteVirtualNetwork.Id != nil {
			remoteVirtualNetworkId = *model.Properties.RemoteVirtualNetwork.Id
		}
		d.Set("remote_virtual_network_id", remoteVirtualNetworkId)
	}

	return nil
}

func resourceDatabricksVirtualNetworkPeeringUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM databricks virtual network peering update.")

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
	}

	// these are the only updatable values, so everything else in the existing.Model should still be unchanged...
	if d.HasChange("allow_forwarded_traffic") {
		existing.Model.Properties.AllowForwardedTraffic = pointer.To(d.Get("allow_forwarded_traffic").(bool))
	}

	if d.HasChange("allow_gateway_transit") {
		existing.Model.Properties.AllowGatewayTransit = pointer.To(d.Get("allow_gateway_transit").(bool))
	}

	if d.HasChange("allow_virtual_network_access") {
		existing.Model.Properties.AllowVirtualNetworkAccess = pointer.To(d.Get("allow_virtual_network_access").(bool))
	}

	if d.HasChange("use_remote_gateways") {
		existing.Model.Properties.UseRemoteGateways = pointer.To(d.Get("use_remote_gateways").(bool))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating Databricks %s: %+v", *id, err)
	}

	return resourceDatabricksVirtualNetworkPeeringRead(d, meta)
}

func resourceDatabricksVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	// Block all changes to any resource of this type...
	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting Databricks %s: %+v", *id, err)
	}

	return nil
}
