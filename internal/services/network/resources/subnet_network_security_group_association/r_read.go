// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_network_security_group_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNetworkSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("`model` was nil for %s", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("`properties` was nil for %s", *id)
	}

	securityGroup := props.NetworkSecurityGroup
	if securityGroup == nil {
		log.Printf("[DEBUG] %s doesn't have a Network Security Group - removing from state!", *id)
		d.SetId("")
		return nil
	}

	d.Set("subnet_id", model.Id)
	d.Set("network_security_group_id", securityGroup.Id)

	return nil
}
