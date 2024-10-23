// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIpGroupCidr() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIpGroupCidrCreate,
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
				ValidateFunc: ipgroups.ValidateIPGroupID,
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

func resourceIpGroupCidrCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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

	if utils.SliceContainsValue(*existing.Model.Properties.IPAddresses, cidr) {
		return tf.ImportAsExistsError("azurerm_ip_group_cidr", id.ID())
	}

	ipAddresses := make([]string, 0)
	if existing.Model.Properties.IPAddresses != nil {
		ipAddresses = *existing.Model.Properties.IPAddresses
	}
	ipAddresses = append(ipAddresses, cidr)

	params := ipgroups.IPGroup{
		Name:     &ipGroupId.IpGroupName,
		Location: existing.Model.Location,
		Tags:     existing.Model.Tags,
		Properties: &ipgroups.IPGroupPropertiesFormat{
			IPAddresses: &ipAddresses,
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *ipGroupId, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupCidrRead(d, meta)
}

func resourceIpGroupCidrRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
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

	if !utils.SliceContainsValue(pointer.From(resp.Model.Properties.IPAddresses), cidr) {
		d.SetId("")
		return nil
	}

	d.Set("ip_group_id", ipGroupId.ID())
	d.Set("cidr", cidr)

	return nil
}

func resourceIpGroupCidrDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupCidrID(d.Id())
	if err != nil {
		return err
	}

	// TODO this resource should use a composite resource ID to remove this instance of d.Get() in the Delete
	// this file can then be removed from the exceptions list in the run-gradually-deprecated.sh script
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
	ipAddresses = utils.RemoveFromStringArray(ipAddresses, cidr)

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
