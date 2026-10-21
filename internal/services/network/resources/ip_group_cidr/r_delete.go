// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group_cidr

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupCidrDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupCidrID(d.Id())
	if err != nil {
		return err
	}

	// TODO this resource should use a composite resource ID to remove this instance of d.Get() in the Delete
	// this file can then be removed from the exceptions list in the gradually-deprecated.sh script
	cidr := d.Get("cidr").(string)
	ipGroupId := ipgroups.NewIPGroupID(id.SubscriptionId, id.ResourceGroup, id.IpGroupName)

	locks.ByID(ipGroupId.ID())
	defer locks.UnlockByID(ipGroupId.ID())

	existing, err := client.Get(ctx, ipGroupId, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving %s: %s", ipGroupId, err)
		}
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", ipGroupId)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", ipGroupId)
	}

	ipAddresses := *existing.Model.Properties.IPAddresses
	ipAddresses = slices.DeleteFunc(ipAddresses, func(s string) bool { return s == cidr })

	params := ipgroups.IPGroup{
		Name:     &ipGroupId.IpGroupName,
		Location: existing.Model.Location,
		Tags:     existing.Model.Tags,
		Properties: &ipgroups.IPGroupPropertiesFormat{
			IPAddresses: &ipAddresses,
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, ipGroupId, params); err != nil {
		return fmt.Errorf("updating %s: %+v", ipGroupId.ID(), err)
	}

	return err
}
