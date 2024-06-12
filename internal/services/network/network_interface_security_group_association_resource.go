// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceSecurityGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceSecurityGroupAssociationCreate,
		Read:   resourceNetworkInterfaceSecurityGroupAssociationRead,
		Delete: resourceNetworkInterfaceSecurityGroupAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateNetworkInterfaceID,
			},

			"network_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
			},
		},
	}
}

func resourceNetworkInterfaceSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	nicId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(nicId.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(nicId.NetworkInterfaceName, networkInterfaceResourceName)

	nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("network_security_group_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(nsgId.NetworkSecurityGroupName, networkSecurityGroupResourceName)
	defer locks.UnlockByName(nsgId.NetworkSecurityGroupName, networkSecurityGroupResourceName)

	read, err := client.Get(ctx, *nicId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", *nicId)
		}
		return fmt.Errorf("retrieving %s: %+v", *nicId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", nicId)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", nicId)
	}

	id := commonids.NewCompositeResourceID(nicId, nsgId)

	if read.Model.Properties.NetworkSecurityGroup != nil {
		return tf.ImportAsExistsError("azurerm_network_interface_security_group_association", id.ID())
	}

	read.Model.Properties.NetworkSecurityGroup = &networkinterfaces.NetworkSecurityGroup{
		Id: pointer.To(nsgId.ID()),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *nicId, *read.Model); err != nil {
		return fmt.Errorf("updating Security Group Association for %s: %+v", *nicId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceSecurityGroupAssociationRead(d, meta)
}

func resourceNetworkInterfaceSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("%s was not found - removing from state!", id.First)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.NetworkSecurityGroup == nil || props.NetworkSecurityGroup.Id == nil {
				log.Printf("%s doesn't have a Security Group attached - removing from state!", id.First)
				d.SetId("")
				return nil
			}
		}
	}

	d.Set("network_interface_id", id.First.ID())
	d.Set("network_security_group_id", id.Second.ID())

	return nil
}

func resourceNetworkInterfaceSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
	if err != nil {
		return err
	}

	locks.ByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf(" %s was not found!", id.First)
		}

		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	read.Model.Properties.NetworkSecurityGroup = nil

	if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *read.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id.First, err)
	}

	return nil
}
