// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_security_group_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
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

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	if props.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil)", id.First)
	}

	info := helpers.ParseFieldsFromNetworkInterface(*props)

	applicationSecurityGroupIds := make([]string, 0)
	for _, v := range info.ApplicationSecurityGroupIDs {
		if v != id.Second.ID() {
			applicationSecurityGroupIds = append(applicationSecurityGroupIds, v)
		}
	}
	info.ApplicationSecurityGroupIDs = applicationSecurityGroupIds
	props.IPConfigurations = helpers.MapFieldsToNetworkInterface(props.IPConfigurations, info)

	if err = client.CreateOrUpdateThenPoll(ctx, *id.First, *read.Model); err != nil {
		return fmt.Errorf("removing Application Security Group for %s: %+v", id.First, err)
	}

	return nil
}
