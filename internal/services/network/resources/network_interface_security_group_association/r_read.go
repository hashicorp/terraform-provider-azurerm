// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_interface_security_group_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NetworkInterfaces
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
