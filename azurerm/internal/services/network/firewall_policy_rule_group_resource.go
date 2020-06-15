package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var azureFirewallPolicyRuleGroupResourceName = "azurerm_firewall_policy_rule_group"

// https://docs.microsoft.com/en-us/rest/api/virtualnetwork/firewallpolicyrulegroups/createorupdate

func resourceArmFirewallPolicyRuleGroup() *schema.Resource {
	application_condition := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"destination_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"fqdn_tags": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"protocols": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.FirewallPolicyFilterRuleActionTypeAllow),
									string(network.FirewallPolicyFilterRuleActionTypeDeny),
								}, false),
							},
							"port": {
								Type:         schema.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 64000),
							},
						},
					},
				},

				"source_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"source_ip_groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"target_fqdns": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	nat_condition := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"destination_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"destination_ports": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"ip_protocols": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"source_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"source_ip_groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	network_condition := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"destination_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"destination_ip_groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"destination_ports": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"ip_protocols": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"source_addresses": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"source_ip_groups": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
	return &schema.Resource{
		Create: resourceArmFirewallPolicyRuleGroupCreateUpdate,
		Read:   resourceArmFirewallPolicyRuleGroupRead,
		Update: resourceArmFirewallPolicyRuleGroupCreateUpdate,
		Delete: resourceArmFirewallPolicyRuleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"firewall_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"filter_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},

						"action_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.FirewallPolicyFilterRuleActionTypeAllow),
								string(network.FirewallPolicyFilterRuleActionTypeDeny),
							}, false),
						},

						"application_condition": application_condition,
						"nat_condition":         nat_condition,
						"network_condition":     network_condition,
					},
				},
			},

			"nat_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},

						"translated_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"translated_port": {
							Type:     schema.TypeString,
							Required: true,
						},

						"application_condition": application_condition,
						"nat_condition":         nat_condition,
						"network_condition":     network_condition,
					},
				},
			},
		},
	}
}

func resourceArmFirewallPolicyRuleGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall Policy Rule Group creation")

	name := d.Get("name").(string)
	firewallPolicyName := d.Get("firewall_policy_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, firewallPolicyName, name)

		if err != nil {
			// TODO currently API returns 400 instead of 404 when not found
			if !utils.ResponseWasNotFound(existing.Response) && !utils.ResponseWasStatusCode(existing.Response, 400) {
				return fmt.Errorf("Error checking for presence of existing Firewall Policy Rule Group %q (Resource Group %q, Policy %q): %s", name, resourceGroup, firewallPolicyName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall_policy_rule_group", *existing.ID)
		}
	}

	locks.ByName(firewallPolicyName, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(firewallPolicyName, azureFirewallPolicyResourceName)

	locks.ByName(name, azureFirewallPolicyRuleGroupResourceName)
	defer locks.UnlockByName(name, azureFirewallPolicyRuleGroupResourceName)

	parameters := network.FirewallPolicyRuleGroup{
		FirewallPolicyRuleGroupProperties: &network.FirewallPolicyRuleGroupProperties{},
	}

	priority := int32(d.Get("priority").(int))
	if priority > 0 {
		parameters.Priority = &priority
	}

	rules := make([]network.BasicFirewallPolicyRule, 0)

	rawRulesFilter := d.Get("filter_rule").([]interface{})
	for _, rawRule := range rawRulesFilter {
		data := rawRule.(map[string]interface{})
		rule, err := mapFirewallPolicyFilterRuleToSDK(data)
		if err != nil {
			return err
		}
		rules = append(rules, rule)
	}

	rawRulesNat := d.Get("nat_rule").([]interface{})
	for _, rawRule := range rawRulesNat {
		data := rawRule.(map[string]interface{})
		rule, err := mapFirewallPolicyNatRuleToSDK(data)
		if err != nil {
			return err
		}
		rules = append(rules, rule)
	}

	parameters.Rules = &rules

	future, err := client.CreateOrUpdate(ctx, resourceGroup, firewallPolicyName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall Policy Rule Group Name %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall Policy Rule Group Name %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, firewallPolicyName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall Policy Rule Group Name %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall Policy Rule Group Name %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmFirewallPolicyRuleGroupRead(d, meta)
}

func mapFirewallPolicyFilterRuleToSDK(data map[string]interface{}) (rule *network.FirewallPolicyFilterRule, err error) {
	rule = &network.FirewallPolicyFilterRule{
		RuleType: network.RuleTypeFirewallPolicyFilterRule,
	}

	if name := data["name"].(string); name != "" {
		rule.Name = &name
	}

	if priority := int32(data["priority"].(int)); priority > 0 {
		rule.Priority = &priority
	}

	if actionType := network.FirewallPolicyFilterRuleActionType(data["action_type"].(string)); actionType != "" {
		rule.Action = &network.FirewallPolicyFilterRuleAction{
			Type: actionType,
		}
	}

	rule.RuleConditions = mapFirewallPolicyRuleConditionsToSDK(data)

	return rule, nil
}

func mapFirewallPolicyNatRuleToSDK(data map[string]interface{}) (rule *network.FirewallPolicyNatRule, err error) {
	rule = &network.FirewallPolicyNatRule{
		RuleType: network.RuleTypeFirewallPolicyNatRule,
	}

	if name := data["name"].(string); name != "" {
		rule.Name = &name
	}

	if priority := int32(data["priority"].(int)); priority > 0 {
		rule.Priority = &priority
	}

	if actionType := network.FirewallPolicyNatRuleActionType(data["action_type"].(string)); actionType != "" {
		rule.Action = &network.FirewallPolicyNatRuleAction{
			Type: actionType,
		}
	}

	conditions := *mapFirewallPolicyRuleConditionsToSDK(data)
	if len(conditions) != 1 {
		return nil, fmt.Errorf("Only a single condition is allowed for NAT rules.")
	}
	rule.RuleCondition = conditions[0]

	return rule, nil
}

func mapFirewallPolicyRuleConditionsToSDK(data map[string]interface{}) *[]network.BasicFirewallPolicyRuleCondition {
	conditions := make([]network.BasicFirewallPolicyRuleCondition, 0)

	rawApplicationConditions := data["application_condition"].([]interface{})
	for _, rawCondition := range rawApplicationConditions {
		data := rawCondition.(map[string]interface{})
		condition := mapFirewallPolicyRuleApplicationConditionsToSDK(data)
		conditions = append(conditions, condition)
	}

	rawNatConditions := data["nat_condition"].([]interface{})
	for _, rawCondition := range rawNatConditions {
		data := rawCondition.(map[string]interface{})
		condition := mapFirewallPolicyRuleNatConditionsToSDK(data)
		conditions = append(conditions, condition)
	}

	rawNetworkConditions := data["network_condition"].([]interface{})
	for _, rawCondition := range rawNetworkConditions {
		data := rawCondition.(map[string]interface{})
		condition := mapFirewallPolicyRuleNetworkConditionsToSDK(data)
		conditions = append(conditions, condition)
	}

	return &conditions
}

func mapFirewallPolicyRuleApplicationConditionsToSDK(data map[string]interface{}) *network.ApplicationRuleCondition {
	condition := &network.ApplicationRuleCondition{
		RuleConditionType: network.RuleConditionTypeApplicationRuleCondition,
	}

	if name := data["name"].(string); name != "" {
		condition.Name = &name
	}

	if destinationAddressesRaw := data["destination_addresses"].([]interface{}); len(destinationAddressesRaw) > 0 {
		destinationAddresses := make([]string, len(destinationAddressesRaw))
		for index, destinationAddressRaw := range destinationAddressesRaw {
			destinationAddresses[index] = destinationAddressRaw.(string)
		}
		condition.DestinationAddresses = &destinationAddresses
	}

	if fqdnTagsRaw := data["fqdn_tags"].([]interface{}); len(fqdnTagsRaw) > 0 {
		fqdnTags := make([]string, len(fqdnTagsRaw))
		for index, fqdnTagRaw := range fqdnTagsRaw {
			fqdnTags[index] = fqdnTagRaw.(string)
		}
		condition.FqdnTags = &fqdnTags
	}

	if protocolsRaw := data["protocols"].([]interface{}); len(protocolsRaw) > 0 {
		protocols := make([]network.FirewallPolicyRuleConditionApplicationProtocol, 0)
		for _, protocolRaw := range protocolsRaw {
			data := protocolRaw.(map[string]interface{})
			port := int32(data["port"].(int))
			protocol := network.FirewallPolicyRuleConditionApplicationProtocol{
				ProtocolType: network.FirewallPolicyRuleConditionApplicationProtocolType(data["type"].(string)),
				Port:         &port,
			}
			protocols = append(protocols, protocol)
		}
		condition.Protocols = &protocols
	}

	if sourceAddressesRaw := data["source_addresses"].([]interface{}); len(sourceAddressesRaw) > 0 {
		sourceAddresses := make([]string, len(sourceAddressesRaw))
		for index, sourceAddressRaw := range sourceAddressesRaw {
			sourceAddresses[index] = sourceAddressRaw.(string)
		}
		condition.SourceAddresses = &sourceAddresses
	}

	if sourceIpGroupsRaw := data["source_ip_groups"].([]interface{}); len(sourceIpGroupsRaw) > 0 {
		sourceIpGroups := make([]string, len(sourceIpGroupsRaw))
		for index, sourceIpGroupRaw := range sourceIpGroupsRaw {
			sourceIpGroups[index] = sourceIpGroupRaw.(string)
		}
		condition.SourceIPGroups = &sourceIpGroups
	}

	if targetFqdnsRaw := data["target_fqdns"].([]interface{}); len(targetFqdnsRaw) > 0 {
		targetFqdns := make([]string, len(targetFqdnsRaw))
		for index, targetFqdnRaw := range targetFqdnsRaw {
			targetFqdns[index] = targetFqdnRaw.(string)
		}
		condition.TargetFqdns = &targetFqdns
	}

	return condition
}

func mapFirewallPolicyRuleNatConditionsToSDK(data map[string]interface{}) *network.NatRuleCondition {
	condition := &network.NatRuleCondition{
		RuleConditionType: network.RuleConditionTypeNatRuleCondition,
	}

	if name := data["name"].(string); name != "" {
		condition.Name = &name
	}

	if destinationAddressesRaw := data["destination_addresses"].([]interface{}); len(destinationAddressesRaw) > 0 {
		destinationAddresses := make([]string, len(destinationAddressesRaw))
		for index, destinationAddresseRaw := range destinationAddressesRaw {
			destinationAddresses[index] = destinationAddresseRaw.(string)
		}
		condition.DestinationAddresses = &destinationAddresses
	}

	if destinationPortsRaw := data["destination_ports"].([]interface{}); len(destinationPortsRaw) > 0 {
		destinationPorts := make([]string, len(destinationPortsRaw))
		for index, destinationPortRaw := range destinationPortsRaw {
			destinationPorts[index] = destinationPortRaw.(string)
		}
		condition.DestinationPorts = &destinationPorts
	}

	if ipProtocolsRaw := data["ip_protocols"].([]interface{}); len(ipProtocolsRaw) > 0 {
		ipProtocols := make([]network.FirewallPolicyRuleConditionNetworkProtocol, len(ipProtocolsRaw))
		for index, ipProtocolRaw := range ipProtocolsRaw {
			ipProtocols[index] = network.FirewallPolicyRuleConditionNetworkProtocol(ipProtocolRaw.(string))
		}
		condition.IPProtocols = &ipProtocols
	}

	if sourceAddressesRaw := data["source_addresses"].([]interface{}); len(sourceAddressesRaw) > 0 {
		sourceAddresses := make([]string, len(sourceAddressesRaw))
		for index, sourceAddressesRaw := range sourceAddressesRaw {
			sourceAddresses[index] = sourceAddressesRaw.(string)
		}
		condition.SourceAddresses = &sourceAddresses
	}

	if sourceIpGroupsRaw := data["source_ip_groups"].([]interface{}); len(sourceIpGroupsRaw) > 0 {
		sourceIpGroups := make([]string, len(sourceIpGroupsRaw))
		for index, sourceIpGroupRaw := range sourceIpGroupsRaw {
			sourceIpGroups[index] = sourceIpGroupRaw.(string)
		}
		condition.SourceIPGroups = &sourceIpGroups
	}

	return condition
}

func mapFirewallPolicyRuleNetworkConditionsToSDK(data map[string]interface{}) *network.RuleCondition {
	condition := &network.RuleCondition{
		RuleConditionType: network.RuleConditionTypeNetworkRuleCondition,
	}

	if name := data["name"].(string); name != "" {
		condition.Name = &name
	}

	if destinationAddressesRaw := data["destination_addresses"].([]interface{}); len(destinationAddressesRaw) > 0 {
		destinationAddresses := make([]string, len(destinationAddressesRaw))
		for index, destinationAddresseRaw := range destinationAddressesRaw {
			destinationAddresses[index] = destinationAddresseRaw.(string)
		}
		condition.DestinationAddresses = &destinationAddresses
	}

	if destinationPortsRaw := data["destination_ports"].([]interface{}); len(destinationPortsRaw) > 0 {
		destinationPorts := make([]string, len(destinationPortsRaw))
		for index, destinationPortRaw := range destinationPortsRaw {
			destinationPorts[index] = destinationPortRaw.(string)
		}
		condition.DestinationPorts = &destinationPorts
	}

	if ipProtocolsRaw := data["ip_protocols"].([]interface{}); len(ipProtocolsRaw) > 0 {
		ipProtocols := make([]network.FirewallPolicyRuleConditionNetworkProtocol, len(ipProtocolsRaw))
		for index, ipProtocolRaw := range ipProtocolsRaw {
			ipProtocols[index] = network.FirewallPolicyRuleConditionNetworkProtocol(ipProtocolRaw.(string))
		}
		condition.IPProtocols = &ipProtocols
	}

	if sourceAddressesRaw := data["source_addresses"].([]interface{}); len(sourceAddressesRaw) > 0 {
		sourceAddresses := make([]string, len(sourceAddressesRaw))
		for index, sourceAddressesRaw := range sourceAddressesRaw {
			sourceAddresses[index] = sourceAddressesRaw.(string)
		}
		condition.SourceAddresses = &sourceAddresses
	}

	if sourceIpGroupsRaw := data["source_ip_groups"].([]interface{}); len(sourceIpGroupsRaw) > 0 {
		sourceIpGroups := make([]string, len(sourceIpGroupsRaw))
		for index, sourceIpGroupRaw := range sourceIpGroupsRaw {
			sourceIpGroups[index] = sourceIpGroupRaw.(string)
		}
		condition.SourceIPGroups = &sourceIpGroups
	}

	return condition
}

func resourceArmFirewallPolicyRuleGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ruleGroups"]
	firewallPolicyName := id.Path["firewallPolicies"]

	read, err := client.Get(ctx, resourceGroup, firewallPolicyName, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Firewall Policy Group Name %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall Policy Rule Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("firewall_policy_name", firewallPolicyName)
	d.Set("resource_group_name", resourceGroup)

	if priority := read.Priority; priority != nil {
		d.Set("priority", read.Priority)
	}

	if rules := read.Rules; rules != nil {
		filterRulesRaw := make([]interface{}, 0)
		natRulesRaw := make([]interface{}, 0)

		for _, rule := range *rules {
			if filterRule, correct := rule.AsFirewallPolicyFilterRule(); correct {
				filterRuleRaw := mapFirewallPolicyFilterRuleFromSDK(filterRule)
				filterRulesRaw = append(filterRulesRaw, filterRuleRaw)
			}
			if natRule, correct := rule.AsFirewallPolicyNatRule(); correct {
				natRuleRaw := mapFirewallPolicyNatRuleFromSDK(natRule)
				natRulesRaw = append(natRulesRaw, natRuleRaw)
			}
		}

		if len(filterRulesRaw) > 0 {
			d.Set("filter_rule", filterRulesRaw)
		}
		if len(natRulesRaw) > 0 {
			d.Set("nat_rule", natRulesRaw)
		}
	}

	return nil
}

func mapFirewallPolicyFilterRuleFromSDK(rule *network.FirewallPolicyFilterRule) *map[string]interface{} {
	raw := make(map[string]interface{})
	if rule.Name != nil {
		raw["name"] = *rule.Name
	}
	if rule.Priority != rule.Priority {
		raw["priority"] = int(*rule.Priority)
	}
	raw["action_type"] = string(rule.Action.Type)

	applicationConditionsRaw := make([]interface{}, 0)
	natConditionsRaw := make([]interface{}, 0)
	networkConditionsRaw := make([]interface{}, 0)
	for _, condition := range *rule.RuleConditions {
		if applicationCondition, correct := condition.AsApplicationRuleCondition(); correct {
			applicationConditionRaw := mapFirewallPolicyApplicationRuleConditionFromSDK(applicationCondition)
			applicationConditionsRaw = append(applicationConditionsRaw, applicationConditionRaw)
		}
		if natCondition, correct := condition.AsNatRuleCondition(); correct {
			natConditionRaw := mapFirewallPolicyNatRuleConditionFromSDK(natCondition)
			natConditionsRaw = append(natConditionsRaw, natConditionRaw)
		}
		if networkCondition, correct := condition.AsRuleCondition(); correct {
			networkConditionRaw := mapFirewallPolicyNetworkRuleConditionFromSDK(networkCondition)
			networkConditionsRaw = append(networkConditionsRaw, networkConditionRaw)
		}
	}
	raw["application_condition"] = applicationConditionsRaw
	raw["nat_condition"] = natConditionsRaw
	raw["network_condition"] = networkConditionsRaw

	return &raw
}

func mapFirewallPolicyNatRuleFromSDK(rule *network.FirewallPolicyNatRule) *map[string]interface{} {
	raw := make(map[string]interface{})
	if rule.Name != nil {
		raw["name"] = *rule.Name
	}
	if rule.Priority != nil {
		raw["priority"] = int(*rule.Priority)
	}
	raw["action_type"] = string(rule.Action.Type)

	raw["translated_address"] = *rule.TranslatedAddress
	raw["translated_port"] = *rule.TranslatedPort

	applicationConditionsRaw := make([]interface{}, 0)
	natConditionsRaw := make([]interface{}, 0)
	networkConditionsRaw := make([]interface{}, 0)
	if applicationCondition, correct := rule.RuleCondition.AsApplicationRuleCondition(); correct {
		applicationConditionRaw := mapFirewallPolicyApplicationRuleConditionFromSDK(applicationCondition)
		applicationConditionsRaw = append(applicationConditionsRaw, applicationConditionRaw)
	}
	if natCondition, correct := rule.RuleCondition.AsNatRuleCondition(); correct {
		natConditionRaw := mapFirewallPolicyNatRuleConditionFromSDK(natCondition)
		natConditionsRaw = append(natConditionsRaw, natConditionRaw)
	}
	if networkCondition, correct := rule.RuleCondition.AsRuleCondition(); correct {
		networkConditionRaw := mapFirewallPolicyNetworkRuleConditionFromSDK(networkCondition)
		networkConditionsRaw = append(networkConditionsRaw, networkConditionRaw)
	}
	raw["application_condition"] = applicationConditionsRaw
	raw["nat_condition"] = natConditionsRaw
	raw["network_condition"] = networkConditionsRaw

	return &raw
}

func mapFirewallPolicyApplicationRuleConditionFromSDK(condition *network.ApplicationRuleCondition) *map[string]interface{} {
	raw := make(map[string]interface{})

	if condition.Name != nil {
		raw["name"] = *condition.Name
	}

	if condition.DestinationAddresses != nil {
		destinationAddressesRaw := make([]string, 0)
		for _, destinationAddress := range *condition.DestinationAddresses {
			destinationAddressesRaw = append(destinationAddressesRaw, destinationAddress)
		}
		raw["destination_addresses"] = destinationAddressesRaw
	}

	if condition.FqdnTags != nil {
		fqdnTagsRaw := make([]string, 0)
		for _, fqdnTag := range *condition.FqdnTags {
			fqdnTagsRaw = append(fqdnTagsRaw, fqdnTag)
		}
		raw["fqdn_tags"] = fqdnTagsRaw
	}

	if condition.Protocols != nil {
		protocolsRaw := make([]interface{}, 0)
		for _, protocol := range *condition.Protocols {
			protocolRaw := make(map[string]interface{}, 0)
			protocolRaw["type"] = string(protocol.ProtocolType)
			protocolRaw["port"] = *protocol.Port
			protocolsRaw = append(protocolsRaw, protocolRaw)
		}
		raw["protocols"] = protocolsRaw
	}

	if condition.SourceAddresses != nil {
		sourceAddressesRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceAddresses {
			sourceAddressesRaw = append(sourceAddressesRaw, sourceAddressRaw)
		}
		raw["source_addresses"] = sourceAddressesRaw
	}

	if condition.SourceIPGroups != nil {
		sourceIPGroupsRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceIPGroups {
			sourceIPGroupsRaw = append(sourceIPGroupsRaw, sourceAddressRaw)
		}
		raw["source_ip_groups"] = sourceIPGroupsRaw
	}

	if condition.TargetFqdns != nil {
		targetFqdnsRaw := make([]string, 0)
		for _, targetFqdnRaw := range *condition.TargetFqdns {
			targetFqdnsRaw = append(targetFqdnsRaw, targetFqdnRaw)
		}
		raw["target_fqdns"] = targetFqdnsRaw
	}

	return &raw
}

func mapFirewallPolicyNatRuleConditionFromSDK(condition *network.NatRuleCondition) *map[string]interface{} {
	raw := make(map[string]interface{})

	if condition.Name != nil {
		raw["name"] = *condition.Name
	}

	if condition.DestinationAddresses != nil {
		destinationAddressesRaw := make([]string, 0)
		for _, destinationAddress := range *condition.DestinationAddresses {
			destinationAddressesRaw = append(destinationAddressesRaw, destinationAddress)
		}
		raw["destination_addresses"] = destinationAddressesRaw
	}

	if condition.DestinationAddresses != nil {
		destinationAddressesRaw := make([]string, 0)
		for _, destinationAddress := range *condition.DestinationAddresses {
			destinationAddressesRaw = append(destinationAddressesRaw, destinationAddress)
		}
		raw["destination_addresses"] = destinationAddressesRaw
	}

	if condition.DestinationPorts != nil {
		destinationPortsRaw := make([]string, 0)
		for _, destinationPortRaw := range *condition.DestinationPorts {
			destinationPortsRaw = append(destinationPortsRaw, destinationPortRaw)
		}
		raw["destination_ports"] = destinationPortsRaw
	}

	if condition.IPProtocols != nil {
		ipProtocolsRaw := make([]string, 0)
		for _, ipProtocolRaw := range *condition.IPProtocols {
			ipProtocolsRaw = append(ipProtocolsRaw, string(ipProtocolRaw))
		}
		raw["ip_protocols"] = ipProtocolsRaw
	}

	if condition.SourceAddresses != nil {
		sourceAddressesRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceAddresses {
			sourceAddressesRaw = append(sourceAddressesRaw, sourceAddressRaw)
		}
		raw["source_addresses"] = sourceAddressesRaw
	}

	if condition.SourceIPGroups != nil {
		sourceIPGroupsRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceIPGroups {
			sourceIPGroupsRaw = append(sourceIPGroupsRaw, sourceAddressRaw)
		}
		raw["source_ip_groups"] = sourceIPGroupsRaw
	}

	return &raw
}

func mapFirewallPolicyNetworkRuleConditionFromSDK(condition *network.RuleCondition) *map[string]interface{} {
	raw := make(map[string]interface{})

	if condition.Name != nil {
		raw["name"] = *condition.Name
	}

	if condition.DestinationAddresses != nil {
		destinationAddressesRaw := make([]string, 0)
		for _, destinationAddress := range *condition.DestinationAddresses {
			destinationAddressesRaw = append(destinationAddressesRaw, destinationAddress)
		}
		raw["destination_addresses"] = destinationAddressesRaw
	}

	if condition.DestinationIPGroups != nil {
		destinationIPGroupsRaw := make([]string, 0)
		for _, destinationIPGroupsaw := range *condition.DestinationIPGroups {
			destinationIPGroupsRaw = append(destinationIPGroupsRaw, destinationIPGroupsaw)
		}
		raw["destination_ip_groups"] = destinationIPGroupsRaw
	}

	if condition.DestinationPorts != nil {
		destinationPortsRaw := make([]string, 0)
		for _, destinationPortRaw := range *condition.DestinationPorts {
			destinationPortsRaw = append(destinationPortsRaw, destinationPortRaw)
		}
		raw["destination_ports"] = destinationPortsRaw
	}

	if condition.IPProtocols != nil {
		ipProtocolsRaw := make([]string, 0)
		for _, ipProtocolRaw := range *condition.IPProtocols {
			ipProtocolsRaw = append(ipProtocolsRaw, string(ipProtocolRaw))
		}
		raw["ip_protocols"] = ipProtocolsRaw
	}

	if condition.SourceAddresses != nil {
		sourceAddressesRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceAddresses {
			sourceAddressesRaw = append(sourceAddressesRaw, sourceAddressRaw)
		}
		raw["source_addresses"] = sourceAddressesRaw
	}

	if condition.SourceIPGroups != nil {
		sourceIPGroupsRaw := make([]string, 0)
		for _, sourceAddressRaw := range *condition.SourceIPGroups {
			sourceIPGroupsRaw = append(sourceIPGroupsRaw, sourceAddressRaw)
		}
		raw["source_ip_groups"] = sourceIPGroupsRaw
	}

	return &raw
}

func resourceArmFirewallPolicyRuleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ruleGroups"]
	firewallPolicyName := id.Path["firewallPolicies"]

	read, err := client.Get(ctx, resourceGroup, firewallPolicyName, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall Policy Rule Group %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Firewall Policy Rule Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	locks.ByName(name, azureFirewallPolicyRuleGroupResourceName)
	defer locks.UnlockByName(name, azureFirewallPolicyRuleGroupResourceName)

	locks.ByName(firewallPolicyName, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(firewallPolicyName, azureFirewallPolicyResourceName)

	future, err := client.Delete(ctx, resourceGroup, firewallPolicyName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall Policy Rule Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Azure Firewall Policy Rule Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}
