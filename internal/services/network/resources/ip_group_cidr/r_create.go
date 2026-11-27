// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ip_group_cidr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupCidrCreate(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).IPGroups
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	cidr := d.Get("cidr").(string)
	cidrName := strings.ReplaceAll(cidr, "/", "_")
	ipGroupId, err := ipgroups.ParseIPGroupID(d.Get("ip_group_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewIpGroupCidrID(ipGroupId.SubscriptionId, ipGroupId.ResourceGroupName, ipGroupId.IpGroupName, cidrName)

	locks.ByID(ipGroupId.ID())
	defer locks.UnlockByID(ipGroupId.ID())

	existing, err := client.Get(ctx, *ipGroupId, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", ipGroupId, err)
		}
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", ipGroupId)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", ipGroupId)
	}

	exists := false
	if slices.Contains(*existing.Model.Properties.IPAddresses, cidr) {
		exists = true
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_ip_group_cidr", id.ID())
		}
	}

	ipAddresses := make([]string, 0)
	if existing.Model.Properties.IPAddresses != nil {
		ipAddresses = *existing.Model.Properties.IPAddresses
	}
	if !exists {
		ipAddresses = append(ipAddresses, cidr)
	}

	params := ipgroups.IPGroup{
		Name:     &ipGroupId.IpGroupName,
		Location: existing.Model.Location,
		Tags:     existing.Model.Tags,
		Properties: &ipgroups.IPGroupPropertiesFormat{
			IPAddresses: &ipAddresses,
		},
	}

	// TODO: implement callback, requires migrating to an ID implementing `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *ipGroupId, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupCidrRead(d, meta)
}
