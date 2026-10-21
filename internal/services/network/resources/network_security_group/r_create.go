// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkSecurityGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkSecurityGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := networksecuritygroups.NewNetworkSecurityGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, networksecuritygroups.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_network_security_group", id.ID())
		}
	}

	sgRules, sgErr := expandSecurityRules(d)
	if sgErr != nil {
		return fmt.Errorf("building list of Network Security Group Rules: %+v", sgErr)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	sg := networksecuritygroups.NetworkSecurityGroup{
		Name:     pointer.To(id.NetworkSecurityGroupName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &networksecuritygroups.NetworkSecurityGroupPropertiesFormat{
			SecurityRules: &sgRules,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, sg, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkSecurityGroupRead(d, meta)
}

func expandSecurityRules(d *pluginsdk.ResourceData) ([]networksecuritygroups.SecurityRule, error) {
	sgRules := d.Get("security_rule").(*pluginsdk.Set).List()
	rules := make([]networksecuritygroups.SecurityRule, 0)

	for _, sgRaw := range sgRules {
		sgRule := sgRaw.(map[string]any)

		if err := validateSecurityRule(sgRule); err != nil {
			return nil, err
		}

		properties := networksecuritygroups.SecurityRulePropertiesFormat{
			SourcePortRange:          pointer.To(sgRule["source_port_range"].(string)),
			DestinationPortRange:     pointer.To(sgRule["destination_port_range"].(string)),
			SourceAddressPrefix:      pointer.To(sgRule["source_address_prefix"].(string)),
			DestinationAddressPrefix: pointer.To(sgRule["destination_address_prefix"].(string)),
			Priority:                 int64(sgRule["priority"].(int)),
			Access:                   networksecuritygroups.SecurityRuleAccess(sgRule["access"].(string)),
			Direction:                networksecuritygroups.SecurityRuleDirection(sgRule["direction"].(string)),
			Protocol:                 networksecuritygroups.SecurityRuleProtocol(sgRule["protocol"].(string)),
		}

		if v := sgRule["description"].(string); v != "" {
			properties.Description = &v
		}

		if r, ok := sgRule["source_port_ranges"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var sourcePortRanges []string
			for _, v := range r.List() {
				s := v.(string)
				sourcePortRanges = append(sourcePortRanges, s)
			}
			properties.SourcePortRanges = &sourcePortRanges
		}

		if r, ok := sgRule["destination_port_ranges"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var destinationPortRanges []string
			for _, v := range r.List() {
				s := v.(string)
				destinationPortRanges = append(destinationPortRanges, s)
			}
			properties.DestinationPortRanges = &destinationPortRanges
		}

		if r, ok := sgRule["source_address_prefixes"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var sourceAddressPrefixes []string
			for _, v := range r.List() {
				s := v.(string)
				sourceAddressPrefixes = append(sourceAddressPrefixes, s)
			}
			properties.SourceAddressPrefixes = &sourceAddressPrefixes
		}

		if r, ok := sgRule["destination_address_prefixes"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var destinationAddressPrefixes []string
			for _, v := range r.List() {
				s := v.(string)
				destinationAddressPrefixes = append(destinationAddressPrefixes, s)
			}
			properties.DestinationAddressPrefixes = &destinationAddressPrefixes
		}

		if r, ok := sgRule["source_application_security_group_ids"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var sourceApplicationSecurityGroups []networksecuritygroups.ApplicationSecurityGroup
			for _, v := range r.List() {
				sg := networksecuritygroups.ApplicationSecurityGroup{
					Id: pointer.To(v.(string)),
				}
				sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
			}
			properties.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
		}

		if r, ok := sgRule["destination_application_security_group_ids"].(*pluginsdk.Set); ok && r.Len() > 0 {
			var destinationApplicationSecurityGroups []networksecuritygroups.ApplicationSecurityGroup
			for _, v := range r.List() {
				sg := networksecuritygroups.ApplicationSecurityGroup{
					Id: pointer.To(v.(string)),
				}
				destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
			}
			properties.DestinationApplicationSecurityGroups = &destinationApplicationSecurityGroups
		}

		rules = append(rules, networksecuritygroups.SecurityRule{
			Name:       pointer.To(sgRule["name"].(string)),
			Properties: &properties,
		})
	}

	return rules, nil
}
