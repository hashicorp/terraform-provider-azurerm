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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var networkSecurityGroupResourceName = "azurerm_network_security_group"

func resourceNetworkSecurityGroup() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceNetworkSecurityGroupCreate,
		Read:   resourceNetworkSecurityGroupRead,
		Update: resourceNetworkSecurityGroupUpdate,
		Delete: resourceNetworkSecurityGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := networksecuritygroups.ParseNetworkSecurityGroupID(id)
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

			"security_rule": {
				Type:       pluginsdk.TypeSet,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
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
								string(networksecuritygroups.SecurityRuleProtocolAny),
								string(networksecuritygroups.SecurityRuleProtocolTcp),
								string(networksecuritygroups.SecurityRuleProtocolUdp),
								string(networksecuritygroups.SecurityRuleProtocolIcmp),
								string(networksecuritygroups.SecurityRuleProtocolAh),
								string(networksecuritygroups.SecurityRuleProtocolEsp),
							}, false),
						},

						"source_port_range": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"source_port_ranges": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_port_range": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"destination_port_ranges": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"source_address_prefix": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"source_address_prefixes": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_address_prefix": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"destination_address_prefixes": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"destination_application_security_group_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"source_application_security_group_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"access": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(networksecuritygroups.SecurityRuleAccessAllow),
								string(networksecuritygroups.SecurityRuleAccessDeny),
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
								string(networksecuritygroups.SecurityRuleDirectionInbound),
								string(networksecuritygroups.SecurityRuleDirectionOutbound),
							}, false),
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}

	return resource
}

func resourceNetworkSecurityGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkSecurityGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := networksecuritygroups.NewNetworkSecurityGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, networksecuritygroups.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_security_group", id.ID())
	}

	sgRules, sgErr := expandSecurityRules(d)
	if sgErr != nil {
		return fmt.Errorf("building list of Network Security Group Rules: %+v", sgErr)
	}

	locks.ByName(id.NetworkSecurityGroupName, networkSecurityGroupResourceName)
	defer locks.UnlockByName(id.NetworkSecurityGroupName, networkSecurityGroupResourceName)

	sg := networksecuritygroups.NetworkSecurityGroup{
		Name:     pointer.To(id.NetworkSecurityGroupName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &networksecuritygroups.NetworkSecurityGroupPropertiesFormat{
			SecurityRules: &sgRules,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, sg); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkSecurityGroupRead(d, meta)
}

func resourceNetworkSecurityGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkSecurityGroups
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, networksecuritygroups.DefaultGetOperationOptions())
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

	if d.HasChange("security_rule") {
		sgRules, sgErr := expandSecurityRules(d)
		if sgErr != nil {
			return fmt.Errorf("building list of Network Security Group Rules: %+v", sgErr)
		}

		payload.Properties.SecurityRules = pointer.To(sgRules)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	locks.ByName(id.NetworkSecurityGroupName, networkSecurityGroupResourceName)
	defer locks.UnlockByName(id.NetworkSecurityGroupName, networkSecurityGroupResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkSecurityGroupRead(d, meta)
}

func resourceNetworkSecurityGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkSecurityGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networksecuritygroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.NetworkSecurityGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			flattenedRules := flattenNetworkSecurityRules(props.SecurityRules)
			if err := d.Set("security_rule", flattenedRules); err != nil {
				return fmt.Errorf("setting `security_rule`: %+v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetworkSecurityGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkSecurityGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return err
}

func expandSecurityRules(d *pluginsdk.ResourceData) ([]networksecuritygroups.SecurityRule, error) {
	sgRules := d.Get("security_rule").(*pluginsdk.Set).List()
	rules := make([]networksecuritygroups.SecurityRule, 0)

	for _, sgRaw := range sgRules {
		sgRule := sgRaw.(map[string]interface{})

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

func flattenNetworkSecurityRules(rules *[]networksecuritygroups.SecurityRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	// For fixing the case insensitive issue for the NSR protocol in Azure
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/16092
	protocolMap := map[string]string{}
	for _, protocol := range networksecuritygroups.PossibleValuesForSecurityRuleProtocol() {
		protocolMap[strings.ToLower(protocol)] = protocol
	}

	if rules != nil {
		for _, rule := range *rules {
			sgRule := make(map[string]interface{})
			sgRule["name"] = *rule.Name

			if props := rule.Properties; props != nil {
				if props.Description != nil {
					sgRule["description"] = *props.Description
				}

				if props.DestinationAddressPrefix != nil {
					sgRule["destination_address_prefix"] = *props.DestinationAddressPrefix
				}
				if props.DestinationAddressPrefixes != nil {
					sgRule["destination_address_prefixes"] = set.FromStringSlice(*props.DestinationAddressPrefixes)
				}
				if props.DestinationPortRange != nil {
					sgRule["destination_port_range"] = *props.DestinationPortRange
				}
				if props.DestinationPortRanges != nil {
					sgRule["destination_port_ranges"] = set.FromStringSlice(*props.DestinationPortRanges)
				}

				destinationApplicationSecurityGroups := make([]string, 0)
				if props.DestinationApplicationSecurityGroups != nil {
					for _, g := range *props.DestinationApplicationSecurityGroups {
						destinationApplicationSecurityGroups = append(destinationApplicationSecurityGroups, *g.Id)
					}
				}
				sgRule["destination_application_security_group_ids"] = set.FromStringSlice(destinationApplicationSecurityGroups)

				if props.SourceAddressPrefix != nil {
					sgRule["source_address_prefix"] = *props.SourceAddressPrefix
				}
				if props.SourceAddressPrefixes != nil {
					sgRule["source_address_prefixes"] = set.FromStringSlice(*props.SourceAddressPrefixes)
				}

				sourceApplicationSecurityGroups := make([]string, 0)
				if props.SourceApplicationSecurityGroups != nil {
					for _, g := range *props.SourceApplicationSecurityGroups {
						sourceApplicationSecurityGroups = append(sourceApplicationSecurityGroups, *g.Id)
					}
				}
				sgRule["source_application_security_group_ids"] = set.FromStringSlice(sourceApplicationSecurityGroups)

				if props.SourcePortRange != nil {
					sgRule["source_port_range"] = *props.SourcePortRange
				}
				if props.SourcePortRanges != nil {
					sgRule["source_port_ranges"] = set.FromStringSlice(*props.SourcePortRanges)
				}

				sgRule["protocol"] = protocolMap[strings.ToLower(string(props.Protocol))]
				sgRule["priority"] = int(props.Priority)
				sgRule["access"] = string(props.Access)
				sgRule["direction"] = string(props.Direction)
			}

			result = append(result, sgRule)
		}
	}

	return result
}

func validateSecurityRule(sgRule map[string]interface{}) error {
	var err *multierror.Error

	sourcePortRange := sgRule["source_port_range"].(string)
	sourcePortRanges := sgRule["source_port_ranges"].(*pluginsdk.Set)
	destinationPortRange := sgRule["destination_port_range"].(string)
	destinationPortRanges := sgRule["destination_port_ranges"].(*pluginsdk.Set)
	sourceAddressPrefix := sgRule["source_address_prefix"].(string)
	sourceAddressPrefixes := sgRule["source_address_prefixes"].(*pluginsdk.Set)
	destinationAddressPrefix := sgRule["destination_address_prefix"].(string)
	destinationAddressPrefixes := sgRule["destination_address_prefixes"].(*pluginsdk.Set)

	if sourcePortRange != "" && sourcePortRanges.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"source_port_range\" and \"source_port_ranges\" can be used per security rule"))
	}
	if destinationPortRange != "" && destinationPortRanges.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"destination_port_range\" and \"destination_port_ranges\" can be used per security rule"))
	}
	if sourceAddressPrefix != "" && sourceAddressPrefixes.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"source_address_prefix\" and \"source_address_prefixes\" can be used per security rule"))
	}
	if destinationAddressPrefix != "" && destinationAddressPrefixes.Len() > 0 {
		err = multierror.Append(err, fmt.Errorf(
			"only one of \"destination_address_prefix\" and \"destination_address_prefixes\" can be used per security rule"))
	}

	return err.ErrorOrNil()
}
