// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group_cidr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupCidrRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupCidrID(d.Id())
	if err != nil {
		return err
	}
	ipGroupId := ipgroups.NewIPGroupID(id.SubscriptionId, id.ResourceGroup, id.IpGroupName)
	cidr := strings.ReplaceAll(id.CidrName, "_", "/")

	resp, err := client.Get(ctx, ipGroupId, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", ipGroupId, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", ipGroupId)
		}
		if resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil", ipGroupId)
		}
	}

	if !slices.Contains(pointer.From(resp.Model.Properties.IPAddresses), cidr) {
		d.SetId("")
		return nil
	}

	d.Set("ip_group_id", ipGroupId.ID())
	d.Set("cidr", cidr)

	return nil
}
