// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securityrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkSecurityRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkSecurityRuleCreateUpdate,
		Read:   resourceNetworkSecurityRuleRead,
		Update: resourceNetworkSecurityRuleCreateUpdate,
		Delete: resourceNetworkSecurityRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securityrules.ParseSecurityRuleID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringDoesNotContainAny("/\\?%"),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"network_security_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 140),
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityrules.SecurityRuleProtocolAny),
					string(securityrules.SecurityRuleProtocolTcp),
					string(securityrules.SecurityRuleProtocolUdp),
					string(securityrules.SecurityRuleProtocolIcmp),
					string(securityrules.SecurityRuleProtocolAh),
					string(securityrules.SecurityRuleProtocolEsp),
				}, false),
			},

			"source_port_range": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_port_ranges"},
			},

			"source_port_ranges": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"source_port_range"},
			},

			"destination_port_range": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_port_ranges"},
			},

			"destination_port_ranges": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"destination_port_range"},
			},

			"source_address_prefix": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source_address_prefixes"},
			},

			"source_address_prefixes": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"source_address_prefix"},
			},

			"destination_address_prefix": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"destination_address_prefixes"},
			},

			"destination_address_prefixes": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:           pluginsdk.HashString,
				ConflictsWith: []string{"destination_address_prefix"},
			},

			//lintignore:S018
			"source_application_security_group_ids": {
				Type:     pluginsdk.TypeSet,
				MaxItems: 10,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			//lintignore:S018
			"destination_application_security_group_ids": {
				Type:     pluginsdk.TypeSet,
				MaxItems: 10,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"access": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityrules.SecurityRuleAccessAllow),
					string(securityrules.SecurityRuleAccessDeny),
				}, false),
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 4096),
			},

			"direction": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityrules.SecurityRuleDirectionInbound),
					string(securityrules.SecurityRuleDirectionOutbound),
				}, false),
			},
		},
	}
}

func resourceNetworkSecurityRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRules
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := securityrules.NewSecurityRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_security_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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
		description := v.(string)
		rule.Properties.Description = &description
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
				Id: utils.String(v.(string)),
			}
			sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, sg)
		}
		rule.Properties.SourceApplicationSecurityGroups = &sourceApplicationSecurityGroups
	}

	if r, ok := d.GetOk("destination_application_security_group_ids"); ok {
		var destinationApplicationSecurityGroups []securityrules.ApplicationSecurityGroup
		for _, v := range r.(*pluginsdk.Set).List() {
			sg := securityrules.ApplicationSecurityGroup{
				Id: utils.String(v.(string)),
			}
			destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, sg)
		}
		rule.Properties.DestinationApplicationSecurityGroups = &destinationApplicationSecurityGroups
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, rule); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkSecurityRuleRead(d, meta)
}

func resourceNetworkSecurityRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securityrules.ParseSecurityRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.SecurityRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("network_security_group_name", id.NetworkSecurityGroupName)

	// For fixing the case insensitive issue for the NSR protocol in Azure
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/16092
	protocolMap := map[string]securityrules.SecurityRuleProtocol{}
	for _, protocol := range securityrules.PossibleValuesForSecurityRuleProtocol() {
		protocolMap[strings.ToLower(protocol)] = securityrules.SecurityRuleProtocol(protocol)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("protocol", string(protocolMap[strings.ToLower(string(props.Protocol))]))
			d.Set("destination_address_prefix", props.DestinationAddressPrefix)
			d.Set("destination_address_prefixes", props.DestinationAddressPrefixes)
			d.Set("destination_port_range", props.DestinationPortRange)
			d.Set("destination_port_ranges", props.DestinationPortRanges)
			d.Set("source_address_prefix", props.SourceAddressPrefix)
			d.Set("source_address_prefixes", props.SourceAddressPrefixes)
			d.Set("source_port_range", props.SourcePortRange)
			d.Set("source_port_ranges", props.SourcePortRanges)
			d.Set("access", string(props.Access))
			d.Set("priority", int(props.Priority))
			d.Set("direction", string(props.Direction))

			if err := d.Set("source_application_security_group_ids", flattenApplicationSecurityGroupIds(props.SourceApplicationSecurityGroups)); err != nil {
				return fmt.Errorf("setting `source_application_security_group_ids`: %+v", err)
			}

			if err := d.Set("destination_application_security_group_ids", flattenApplicationSecurityGroupIds(props.DestinationApplicationSecurityGroups)); err != nil {
				return fmt.Errorf("setting `source_application_security_group_ids`: %+v", err)
			}
		}
	}

	return nil
}

func resourceNetworkSecurityRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityRules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securityrules.ParseSecurityRuleID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenApplicationSecurityGroupIds(groups *[]securityrules.ApplicationSecurityGroup) []string {
	ids := make([]string, 0)

	if groups != nil {
		for _, v := range *groups {
			ids = append(ids, *v.Id)
		}
	}

	return ids
}
