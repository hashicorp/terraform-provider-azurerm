// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkSecurityRuleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.SecurityRules
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := securityrules.NewSecurityRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_security_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_network_security_rule", id.ID())
		}
	}

	rule := securityrules.SecurityRule{
		Name: &id.SecurityRuleName,
		Properties: &securityrules.SecurityRulePropertiesFormat{
			SourcePortRange:          pointer.To(d.Get("source_port_range").(string)),
			DestinationPortRange:     pointer.To(d.Get("destination_port_range").(string)),
			SourceAddressPrefix:      pointer.To(d.Get("source_address_prefix").(string)),
			DestinationAddressPrefix: pointer.To(d.Get("destination_address_prefix").(string)),
			Priority:                 int64(d.Get("priority").(int)),
			Access:                   securityrules.SecurityRuleAccess(d.Get("access").(string)),
			Direction:                securityrules.SecurityRuleDirection(d.Get("direction").(string)),
			Protocol:                 securityrules.SecurityRuleProtocol(d.Get("protocol").(string)),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Properties.Description = pointer.To(v.(string))
	}

	if r, ok := d.GetOk("source_port_ranges"); ok {
		var sourcePortRanges []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourcePortRanges = append(sourcePortRanges, s)
		}
		rule.Properties.SourcePortRanges = &sourcePortRanges
	}

	if r, ok := d.GetOk("destination_port_ranges"); ok {
		var destinationPortRanges []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationPortRanges = append(destinationPortRanges, s)
		}
		rule.Properties.DestinationPortRanges = &destinationPortRanges
	}

	if r, ok := d.GetOk("source_address_prefixes"); ok {
		var sourceAddressPrefixes []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourceAddressPrefixes = append(sourceAddressPrefixes, s)
		}
		rule.Properties.SourceAddressPrefixes = &sourceAddressPrefixes
	}

	if r, ok := d.GetOk("destination_address_prefixes"); ok {
		var destinationAddressPrefixes []string
		r := r.(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationAddressPrefixes = append(destinationAddressPrefixes, s)
		}
		rule.Properties.DestinationAddressPrefixes = &destinationAddressPrefixes
	}

	if r, ok := d.GetOk("source_application_security_group_ids"); ok {
		var sourceApplicationSecurityGroups []securityrules.ApplicationSecurityGroup
		for _, v := range r.(*pluginsdk.Set).List() {
			sg := securityrules.ApplicationSecurityGroup{
				Id: pointer.To(v.(string)),
			}
			sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		}
		rule.Properties.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
	}

	if r, ok := d.GetOk("destination_application_security_group_ids"); ok {
		var destinationApplicationSecurityGroups []securityrules.ApplicationSecurityGroup
		for _, v := range r.(*pluginsdk.Set).List() {
			sg := securityrules.ApplicationSecurityGroup{
				Id: pointer.To(v.(string)),
			}
			destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
		}
		rule.Properties.DestinationApplicationSecurityGroups = &destinationApplicationSecurityGroups
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, rule, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkSecurityRuleRead(d, meta)
}
