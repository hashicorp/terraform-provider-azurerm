// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_security_group_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.First)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if model := read.Model; model != nil {
		nicProps := read.Model.Properties
		if nicProps == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
		}

		info := helpers.ParseFieldsFromNetworkInterface(*nicProps)
		exists := false
		for _, groupId := range info.ApplicationSecurityGroupIDs {
			if groupId == id.Second.ID() {
				exists = true
			}
		}

		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.First, id.Second)
			d.SetId("")
			return nil
		}
	}

	d.Set("application_security_group_id", id.Second.ID())
	d.Set("network_interface_id", id.First.ID())

	return nil
}
