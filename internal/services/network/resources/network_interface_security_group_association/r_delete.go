// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_interface_security_group_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &networksecuritygroups.NetworkSecurityGroupId{})
	if err != nil {
		return err
	}

	locks.ByID(strings.ToLower(id.First.ID()))
	defer locks.UnlockByID(strings.ToLower(id.First.ID()))

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", id.First)
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
