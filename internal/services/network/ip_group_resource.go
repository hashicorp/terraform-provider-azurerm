// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIpGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIpGroupCreate,
		Read:   resourceIpGroupRead,
		Update: resourceIpGroupUpdate,
		Delete: resourceIpGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := ipgroups.ParseIPGroupID(id)
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceIpGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]interface{}) {
		id, err := azurefirewalls.ParseAzureFirewallID(fw.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", fw, err)
		}
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]interface{}) {
		id, err := firewallpolicies.ParseFirewallPolicyID(fwpol.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", fwpol, err)
		}
		locks.ByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	id := ipgroups.NewIPGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_ip_group", id.ID())
	}

	ipAddresses := d.Get("cidrs").(*pluginsdk.Set).List()

	sg := ipgroups.IPGroup{
		Name:     &id.IpGroupName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &ipgroups.IPGroupPropertiesFormat{
			IPAddresses: utils.ExpandStringSlice(ipAddresses),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, sg); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupRead(d, meta)
}

func resourceIpGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.IpGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if props.IPAddresses == nil {
				return fmt.Errorf("list of ipAddresses returned is nil")
			}
			if err := d.Set("cidrs", props.IPAddresses); err != nil {
				return fmt.Errorf("setting `cidrs`: %+v", err)
			}

			firewallIDs := make([]string, 0)
			for _, idStr := range getIds(props.Firewalls) {
				firewallID, err := azurefirewalls.ParseAzureFirewallIDInsensitively(idStr)
				if err != nil {
					return fmt.Errorf("parsing Azure Firewall ID %q: %+v", idStr, err)
				}
				firewallIDs = append(firewallIDs, firewallID.ID())
			}
			d.Set("firewall_ids", firewallIDs)

			firewallPolicyIDs := make([]string, 0)
			for _, idStr := range getIds(props.FirewallPolicies) {
				policyID, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(idStr)
				if err != nil {
					return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", idStr, err)
				}
				firewallPolicyIDs = append(firewallPolicyIDs, policyID.ID())
			}
			d.Set("firewall_policy_ids", firewallPolicyIDs)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceIpGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]interface{}) {
		id, err := azurefirewalls.ParseAzureFirewallID(fw.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", fw, err)
		}
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]interface{}) {
		id, err := firewallpolicies.ParseFirewallPolicyID(fwpol.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", fwpol, err)
		}
		locks.ByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("cidrs") {
		payload.Properties.IPAddresses = utils.ExpandStringSlice(d.Get("cidrs").(*pluginsdk.Set).List())
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIpGroupRead(d, meta)
}

func getIds(subResource *[]ipgroups.SubResource) []string {
	if subResource == nil {
		return nil
	}

	ids := make([]string, 0)
	for _, v := range *subResource {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}

func resourceIpGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.IPGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ipgroups.ParseIPGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.Get(ctx, *id, ipgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			return nil
		}
		return fmt.Errorf("retrieving %s : %+v", *id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	for _, fw := range *resp.Model.Properties.Firewalls {
		fwID, err := azurefirewalls.ParseAzureFirewallID(pointer.From(fw.Id))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", pointer.From(fw.Id), err)
		}
		locks.ByName(fwID.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(fwID.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range *resp.Model.Properties.FirewallPolicies {
		polID, err := firewallpolicies.ParseFirewallPolicyID(pointer.From(fwpol.Id))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", *fwpol.Id, err)
		}
		locks.ByName(polID.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(polID.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return err
}
