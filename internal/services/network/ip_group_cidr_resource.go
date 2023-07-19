// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceIpGroupCidr() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIpGroupCidrCreateUpdate,
		Read:   resourceIpGroupCidrRead,
		Delete: resourceIpGroupCidrDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IpGroupCidrID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"ip_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IpGroupID,
			},
			"cidr": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceIpGroupCidrCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	cidr := d.Get("cidr").(string)
	cidrName := strings.ReplaceAll(cidr, "/", "_")
	ipGroupId, err := parse.IpGroupID(d.Get("ip_group_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewIpGroupCidrID(subscriptionId, ipGroupId.ResourceGroup, ipGroupId.Name, cidrName)

	locks.ByID(ipGroupId.ID())
	defer locks.UnlockByID(ipGroupId.ID())

	existing, err := client.Get(ctx, ipGroupId.ResourceGroup, ipGroupId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", ipGroupId, err)
		}
	}

	if d.IsNewResource() {
		if utils.SliceContainsValue(*existing.IPAddresses, cidr) {
			return tf.ImportAsExistsError("azurerm_ip_group_cidr", id.ID())
		}
	}

	ipAddresses := make([]string, 0)
	if existing.IPAddresses != nil {
		ipAddresses = *existing.IPAddresses
	}
	ipAddresses = append(ipAddresses, cidr)

	params := network.IPGroup{
		Name:     &ipGroupId.Name,
		Location: existing.Location,
		Tags:     existing.Tags,
		IPGroupPropertiesFormat: &network.IPGroupPropertiesFormat{
			IPAddresses: &ipAddresses,
		},
	}

	future, err := client.CreateOrUpdate(ctx, ipGroupId.ResourceGroup, ipGroupId.Name, params)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupCidrRead(d, meta)
}

func resourceIpGroupCidrRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupCidrID(d.Id())
	if err != nil {
		return err
	}
	ipGroupId := parse.NewIpGroupID(id.SubscriptionId, id.ResourceGroup, id.IpGroupName)
	cidr := strings.ReplaceAll(id.CidrName, "_", "/")

	resp, err := client.Get(ctx, id.ResourceGroup, id.IpGroupName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("making Read request on IP Group %q (Resource Group %q): %+v", ipGroupId.Name, ipGroupId.ResourceGroup, err)
		}
		if !utils.SliceContainsValue(*resp.IPAddresses, cidr) {
			d.SetId("")
			return nil
		}
	}

	d.Set("ip_group_id", ipGroupId.ID())
	d.Set("cidr", cidr)

	return nil
}

func resourceIpGroupCidrDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	cidr := d.Get("cidr").(string)
	ipGroupId, err := parse.IpGroupID(d.Get("ip_group_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(ipGroupId.ID())
	defer locks.UnlockByID(ipGroupId.ID())

	existing, err := client.Get(ctx, ipGroupId.ResourceGroup, ipGroupId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("reading existing %s: %s", ipGroupId, err)
		}
	}

	ipAddresses := *existing.IPAddresses
	ipAddresses = utils.RemoveFromStringArray(ipAddresses, cidr)

	params := network.IPGroup{
		Name:     &ipGroupId.Name,
		Location: existing.Location,
		Tags:     existing.Tags,
		IPGroupPropertiesFormat: &network.IPGroupPropertiesFormat{
			IPAddresses: &ipAddresses,
		},
	}

	future, err := client.CreateOrUpdate(ctx, ipGroupId.ResourceGroup, ipGroupId.Name, params)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", ipGroupId.ID(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("deleting IP Group CIDR %q (IP Group %q - Resource Group %q): %+v", cidr, ipGroupId.Name, ipGroupId.ResourceGroup, err)
	}

	return err
}
