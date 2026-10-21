// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkSecurityRuleUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.SecurityRules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securityrules.ParseSecurityRuleID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("protocol") {
		payload.Properties.Protocol = securityrules.SecurityRuleProtocol(d.Get("protocol").(string))
	}

	if d.HasChange("source_port_range") {
		payload.Properties.SourcePortRange = pointer.To(d.Get("source_port_range").(string))
	}

	if d.HasChange("source_port_ranges") {
		var sourcePortRanges []string
		r := d.Get("source_port_ranges").(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourcePortRanges = append(sourcePortRanges, s)
		}
		payload.Properties.SourcePortRanges = pointer.To(sourcePortRanges)
	}

	if d.HasChange("destination_port_range") {
		payload.Properties.DestinationPortRange = pointer.To(d.Get("destination_port_range").(string))
	}

	if d.HasChange("destination_port_ranges") {
		var destinationPortRanges []string
		r := d.Get("destination_port_ranges").(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationPortRanges = append(destinationPortRanges, s)
		}
		payload.Properties.DestinationPortRanges = pointer.To(destinationPortRanges)
	}

	if d.HasChange("source_address_prefix") {
		payload.Properties.SourceAddressPrefix = pointer.To(d.Get("source_address_prefix").(string))
	}

	if d.HasChange("source_address_prefixes") {
		var sourceAddressPrefixes []string
		r := d.Get("source_address_prefixes").(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			sourceAddressPrefixes = append(sourceAddressPrefixes, s)
		}
		payload.Properties.SourceAddressPrefixes = pointer.To(sourceAddressPrefixes)
	}

	if d.HasChange("destination_address_prefix") {
		payload.Properties.DestinationAddressPrefix = pointer.To(d.Get("destination_address_prefix").(string))
	}

	if d.HasChange("destination_address_prefixes") {
		var destinationAddressPrefixes []string
		r := d.Get("destination_address_prefixes").(*pluginsdk.Set).List()
		for _, v := range r {
			s := v.(string)
			destinationAddressPrefixes = append(destinationAddressPrefixes, s)
		}
		payload.Properties.DestinationAddressPrefixes = pointer.To(destinationAddressPrefixes)
	}

	if d.HasChange("source_application_security_group_ids") {
		var sourceApplicationSecurityGroups []securityrules.ApplicationSecurityGroup
		for _, v := range d.Get("source_application_security_group_ids").(*pluginsdk.Set).List() {
			sg := securityrules.ApplicationSecurityGroup{
				Id: pointer.To(v.(string)),
			}
			sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		}
		payload.Properties.SourceApplicationSecurityGroups = pointer.To(sourceApplicationSecurityGroups)
	}

	if d.HasChange("destination_application_security_group_ids") {
		var destinationApplicationSecurityGroups []securityrules.ApplicationSecurityGroup
		for _, v := range d.Get("destination_application_security_group_ids").(*pluginsdk.Set).List() {
			sg := securityrules.ApplicationSecurityGroup{
				Id: pointer.To(v.(string)),
			}
			destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
		}
		payload.Properties.DestinationApplicationSecurityGroups = pointer.To(destinationApplicationSecurityGroups)
	}

	if d.HasChange("access") {
		payload.Properties.Access = securityrules.SecurityRuleAccess(d.Get("access").(string))
	}

	if d.HasChange("priority") {
		payload.Properties.Priority = int64(d.Get("priority").(int))
	}

	if d.HasChange("direction") {
		payload.Properties.Direction = securityrules.SecurityRuleDirection(d.Get("direction").(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkSecurityRuleRead(d, meta)
}
