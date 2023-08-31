// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	firewallParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceIpGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIpGroupCreate,
		Read:   resourceIpGroupRead,
		Update: resourceIpGroupUpdate,
		Delete: resourceIpGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IpGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"firewall_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"firewall_policy_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"cidrs": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: pluginsdk.HashString,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIpGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]interface{}) {
		id, _ := firewallParse.FirewallID(fw.(string))
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]interface{}) {
		id, _ := firewallParse.FirewallPolicyID(fwpol.(string))
		locks.ByName(id.Name, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.Name, firewall.AzureFirewallPolicyResourceName)
	}

	id := parse.NewIpGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_ip_group", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	ipAddresses := d.Get("cidrs").(*pluginsdk.Set).List()

	sg := network.IPGroup{
		Name:     &id.Name,
		Location: &location,
		IPGroupPropertiesFormat: &network.IPGroupPropertiesFormat{
			IPAddresses: utils.ExpandStringSlice(ipAddresses),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, sg)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupRead(d, meta)
}

func resourceIpGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.IPGroupPropertiesFormat; props != nil {
		if props.IPAddresses == nil {
			return fmt.Errorf("list of ipAddresses returned is nil")
		}
		if err := d.Set("cidrs", props.IPAddresses); err != nil {
			return fmt.Errorf("setting `cidrs`: %+v", err)
		}
	}

	d.Set("firewall_ids", getIds(resp.Firewalls))
	d.Set("firewall_policy_ids", getIds(resp.FirewallPolicies))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIpGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]interface{}) {
		id, _ := firewallParse.FirewallID(fw.(string))
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]interface{}) {
		id, _ := firewallParse.FirewallPolicyID(fwpol.(string))
		locks.ByName(id.Name, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.Name, firewall.AzureFirewallPolicyResourceName)
	}

	id := parse.NewIpGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	exisiting, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(exisiting.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if d.HasChange("cidrs") {
		if exisiting.IPGroupPropertiesFormat != nil {
			exisiting.IPGroupPropertiesFormat.IPAddresses = utils.ExpandStringSlice(d.Get("cidrs").(*pluginsdk.Set).List())
		}
	}

	if d.HasChange("tags") {
		exisiting.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, exisiting)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupRead(d, meta)
}

func getIds(subResource *[]network.SubResource) []string {
	if subResource == nil {
		return nil
	}

	ids := make([]string, 0)
	for _, v := range *subResource {
		if v.ID == nil {
			continue
		}

		ids = append(ids, *v.ID)
	}

	return ids
}

func resourceIpGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.IPGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IpGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] IP Group %q was not found in Resource Group %q - assuming removed!", id.Name, id.ResourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving ip group %s : %+v", *id, err)
	}

	for _, fw := range *read.Firewalls {
		id, _ := firewallParse.FirewallID(*fw.ID)
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range *read.FirewallPolicies {
		id, _ := firewallParse.FirewallPolicyID(*fwpol.ID)
		locks.ByName(id.Name, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.Name, firewall.AzureFirewallPolicyResourceName)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("deleting IP Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return err
}
