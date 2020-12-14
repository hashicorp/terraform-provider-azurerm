package firewall

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/firewall/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFirewallPolicyRuleCollectionGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallPolicyRuleCollectionGroupCreateUpdate,
		Read:   resourceArmFirewallPolicyRuleCollectionGroupRead,
		Update: resourceArmFirewallPolicyRuleCollectionGroupCreateUpdate,
		Delete: resourceArmFirewallPolicyRuleCollectionGroupDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.FirewallPolicyRuleCollectionGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallPolicyRuleCollectionGroupName(),
			},

			"firewall_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallPolicyID,
			},

			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"application_rule_collection": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.FirewallPolicyFilterRuleCollectionActionTypeAllow),
								string(network.FirewallPolicyFilterRuleCollectionActionTypeDeny),
							}, false),
						},
						"rule": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.FirewallPolicyRuleName(),
									},
									"protocols": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(network.FirewallPolicyRuleApplicationProtocolTypeHTTP),
														string(network.FirewallPolicyRuleApplicationProtocolTypeHTTPS),
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
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_fqdns": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_fqdn_tags": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},
			},

			"network_rule_collection": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.FirewallPolicyFilterRuleCollectionActionTypeAllow),
								string(network.FirewallPolicyFilterRuleCollectionActionTypeDeny),
							}, false),
						},
						"rule": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.FirewallPolicyRuleName(),
									},
									"protocols": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(network.FirewallPolicyRuleNetworkProtocolAny),
												string(network.FirewallPolicyRuleNetworkProtocolTCP),
												string(network.FirewallPolicyRuleNetworkProtocolUDP),
												string(network.FirewallPolicyRuleNetworkProtocolICMP),
											}, false),
										},
									},
									"source_addresses": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_addresses": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											// Can be IP address, CIDR, "*", or service tag
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_ip_groups": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_fqdns": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_ports": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validate.FirewallPolicyRulePort,
										},
									},
								},
							},
						},
					},
				},
			},

			"nat_rule_collection": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								// Hardcode to using `Dnat` instead of the one defined in Swagger (i.e. network.DNAT) because of: https://github.com/Azure/azure-rest-api-specs/issues/9986
								// Setting `StateFunc: state.IgnoreCase` will cause other issues, as tracked by: https://github.com/hashicorp/terraform-plugin-sdk/issues/485
								// Another solution is to customize the hash function for the containing block, but as there are a couple of properties here, especially
								// has property whose type is another nested block (Set), so the implementation is nontrivial and error-prone.
								"Dnat",
							}, false),
						},
						"rule": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.FirewallPolicyRuleName(),
									},
									"protocols": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(network.FirewallPolicyRuleNetworkProtocolTCP),
												string(network.FirewallPolicyRuleNetworkProtocolUDP),
											}, false),
										},
									},
									"source_addresses": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_address": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IsIPAddress,
											validation.IsCIDR,
										),
									},
									"destination_ports": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validate.FirewallPolicyRulePort,
										},
									},
									"translated_address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"translated_port": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IsPortNumber,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmFirewallPolicyRuleCollectionGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyId, err := parse.FirewallPolicyID(d.Get("firewall_policy_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		resp, err := client.Get(ctx, policyId.ResourceGroup, policyId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q): %+v", name, policyId.ResourceGroup, policyId.Name, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall_policy_rule_collection_group", *resp.ID)
		}
	}

	locks.ByName(policyId.Name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(policyId.Name, azureFirewallPolicyResourceName)

	param := network.FirewallPolicyRuleCollectionGroup{
		FirewallPolicyRuleCollectionGroupProperties: &network.FirewallPolicyRuleCollectionGroupProperties{
			Priority: utils.Int32(int32(d.Get("priority").(int))),
		},
	}
	var rulesCollections []network.BasicFirewallPolicyRuleCollection
	rulesCollections = append(rulesCollections, expandAzureRmFirewallPolicyRuleCollectionApplication(d.Get("application_rule_collection").(*schema.Set).List())...)
	rulesCollections = append(rulesCollections, expandAzureRmFirewallPolicyRuleCollectionNetwork(d.Get("network_rule_collection").(*schema.Set).List())...)
	rulesCollections = append(rulesCollections, expandAzureRmFirewallPolicyRuleCollectionNat(d.Get("nat_rule_collection").(*schema.Set).List())...)
	param.FirewallPolicyRuleCollectionGroupProperties.RuleCollections = &rulesCollections

	future, err := client.CreateOrUpdate(ctx, policyId.ResourceGroup, policyId.Name, name, param)
	if err != nil {
		return fmt.Errorf("creating Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q): %+v", name, policyId.ResourceGroup, policyId.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q): %+v", name, policyId.ResourceGroup, policyId.Name, err)
	}

	resp, err := client.Get(ctx, policyId.ResourceGroup, policyId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q): %+v", name, policyId.ResourceGroup, policyId.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q) ID", name, policyId.ResourceGroup, policyId.Name)
	}
	id, err := parse.FirewallPolicyRuleCollectionGroupID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID())

	return resourceArmFirewallPolicyRuleCollectionGroupRead(d, meta)
}

func resourceArmFirewallPolicyRuleCollectionGroupRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyRuleCollectionGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Firewall Policy Rule Collection Group %q was not found in Resource Group %q - removing from state!", id.RuleCollectionGroupName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q): %+v", id.RuleCollectionGroupName, id.ResourceGroup, id.FirewallPolicyName, err)
	}

	d.Set("name", resp.Name)
	d.Set("priority", resp.Priority)
	d.Set("firewall_policy_id", parse.NewFirewallPolicyID(subscriptionId, id.ResourceGroup, id.FirewallPolicyName).ID())

	applicationRuleCollections, networkRuleCollections, natRuleCollections, err := flattenAzureRmFirewallPolicyRuleCollection(resp.RuleCollections)
	if err != nil {
		return fmt.Errorf("flattening Firewall Policy Rule Collections: %+v", err)
	}

	if err := d.Set("application_rule_collection", applicationRuleCollections); err != nil {
		return fmt.Errorf("setting `application_rule_collection`: %+v", err)
	}
	if err := d.Set("network_rule_collection", networkRuleCollections); err != nil {
		return fmt.Errorf("setting `network_rule_collection`: %+v", err)
	}
	if err := d.Set("nat_rule_collection", natRuleCollections); err != nil {
		return fmt.Errorf("setting `nat_rule_collection`: %+v", err)
	}

	return nil
}

func resourceArmFirewallPolicyRuleCollectionGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyRuleCollectionGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FirewallPolicyName, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, azureFirewallPolicyResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
	if err != nil {
		return fmt.Errorf("deleting Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q): %+v", id.RuleCollectionGroupName, id.ResourceGroup, id.FirewallPolicyName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting %q (Resource Group %q / Policy: %q): %+v", id.RuleCollectionGroupName, id.ResourceGroup, id.FirewallPolicyName, err)
		}
	}

	return nil
}

func expandAzureRmFirewallPolicyRuleCollectionApplication(input []interface{}) []network.BasicFirewallPolicyRuleCollection {
	return expandAzureRmFirewallPolicyFilterRuleCollection(input, expandAzureRmFirewallPolicyRuleApplication)
}

func expandAzureRmFirewallPolicyRuleCollectionNetwork(input []interface{}) []network.BasicFirewallPolicyRuleCollection {
	return expandAzureRmFirewallPolicyFilterRuleCollection(input, expandAzureRmFirewallPolicyRuleNetwork)
}

func expandAzureRmFirewallPolicyRuleCollectionNat(input []interface{}) []network.BasicFirewallPolicyRuleCollection {
	result := make([]network.BasicFirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]interface{})
		output := &network.FirewallPolicyNatRuleCollection{
			RuleCollectionType: network.RuleCollectionTypeFirewallPolicyNatRuleCollection,
			Name:               utils.String(rule["name"].(string)),
			Priority:           utils.Int32(int32(rule["priority"].(int))),
			Action: &network.FirewallPolicyNatRuleCollectionAction{
				Type: network.FirewallPolicyNatRuleCollectionActionType(rule["action"].(string)),
			},
			Rules: expandAzureRmFirewallPolicyRuleNat(rule["rule"].(*schema.Set).List()),
		}
		result = append(result, output)
	}
	return result
}

func expandAzureRmFirewallPolicyFilterRuleCollection(input []interface{}, f func(input []interface{}) *[]network.BasicFirewallPolicyRule) []network.BasicFirewallPolicyRuleCollection {
	result := make([]network.BasicFirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]interface{})
		output := &network.FirewallPolicyFilterRuleCollection{
			Action: &network.FirewallPolicyFilterRuleCollectionAction{
				Type: network.FirewallPolicyFilterRuleCollectionActionType(rule["action"].(string)),
			},
			Name:               utils.String(rule["name"].(string)),
			Priority:           utils.Int32(int32(rule["priority"].(int))),
			RuleCollectionType: network.RuleCollectionTypeFirewallPolicyFilterRuleCollection,
			Rules:              f(rule["rule"].(*schema.Set).List()),
		}
		result = append(result, output)
	}
	return result
}

func expandAzureRmFirewallPolicyRuleApplication(input []interface{}) *[]network.BasicFirewallPolicyRule {
	result := make([]network.BasicFirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []network.FirewallPolicyRuleApplicationProtocol
		for _, p := range condition["protocols"].(*schema.Set).List() {
			proto := p.(map[string]interface{})
			protocols = append(protocols, network.FirewallPolicyRuleApplicationProtocol{
				ProtocolType: network.FirewallPolicyRuleApplicationProtocolType(proto["type"].(string)),
				Port:         utils.Int32(int32(proto["port"].(int))),
			})
		}
		output := &network.ApplicationRule{
			Name:            utils.String(condition["name"].(string)),
			RuleType:        network.RuleTypeApplicationRule,
			Protocols:       &protocols,
			SourceAddresses: utils.ExpandStringSlice(condition["source_addresses"].(*schema.Set).List()),
			SourceIPGroups:  utils.ExpandStringSlice(condition["source_ip_groups"].(*schema.Set).List()),
			TargetFqdns:     utils.ExpandStringSlice(condition["destination_fqdns"].(*schema.Set).List()),
			FqdnTags:        utils.ExpandStringSlice(condition["destination_fqdn_tags"].(*schema.Set).List()),
		}
		result = append(result, output)
	}
	return &result
}

func expandAzureRmFirewallPolicyRuleNetwork(input []interface{}) *[]network.BasicFirewallPolicyRule {
	result := make([]network.BasicFirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []network.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].(*schema.Set).List() {
			protocols = append(protocols, network.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}
		output := &network.Rule{
			Name:                 utils.String(condition["name"].(string)),
			RuleType:             network.RuleTypeNetworkRule,
			IPProtocols:          &protocols,
			SourceAddresses:      utils.ExpandStringSlice(condition["source_addresses"].(*schema.Set).List()),
			SourceIPGroups:       utils.ExpandStringSlice(condition["source_ip_groups"].(*schema.Set).List()),
			DestinationAddresses: utils.ExpandStringSlice(condition["destination_addresses"].(*schema.Set).List()),
			DestinationIPGroups:  utils.ExpandStringSlice(condition["destination_ip_groups"].(*schema.Set).List()),
			DestinationFqdns:     utils.ExpandStringSlice(condition["destination_fqdns"].(*schema.Set).List()),
			DestinationPorts:     utils.ExpandStringSlice(condition["destination_ports"].(*schema.Set).List()),
		}
		result = append(result, output)
	}
	return &result
}

func expandAzureRmFirewallPolicyRuleNat(input []interface{}) *[]network.BasicFirewallPolicyRule {
	result := make([]network.BasicFirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []network.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].(*schema.Set).List() {
			protocols = append(protocols, network.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}
		destinationAddresses := []string{condition["destination_address"].(string)}
		output := &network.NatRule{
			Name:                 utils.String(condition["name"].(string)),
			RuleType:             network.RuleTypeNatRule,
			IPProtocols:          &protocols,
			SourceAddresses:      utils.ExpandStringSlice(condition["source_addresses"].(*schema.Set).List()),
			SourceIPGroups:       utils.ExpandStringSlice(condition["source_ip_groups"].(*schema.Set).List()),
			DestinationAddresses: &destinationAddresses,
			DestinationPorts:     utils.ExpandStringSlice(condition["destination_ports"].(*schema.Set).List()),
			TranslatedAddress:    utils.String(condition["translated_address"].(string)),
			TranslatedPort:       utils.String(strconv.Itoa(condition["translated_port"].(int))),
		}
		result = append(result, output)
	}
	return &result
}

func flattenAzureRmFirewallPolicyRuleCollection(input *[]network.BasicFirewallPolicyRuleCollection) ([]interface{}, []interface{}, []interface{}, error) {
	var (
		applicationRuleCollection = []interface{}{}
		networkRuleCollection     = []interface{}{}
		natRuleCollection         = []interface{}{}
	)
	if input == nil {
		return applicationRuleCollection, networkRuleCollection, natRuleCollection, nil
	}

	for _, e := range *input {
		var result map[string]interface{}

		switch rule := e.(type) {
		case network.FirewallPolicyFilterRuleCollection:
			var name string
			if rule.Name != nil {
				name = *rule.Name
			}
			var priority int32
			if rule.Priority != nil {
				priority = *rule.Priority
			}

			var action string
			if rule.Action != nil {
				action = string(rule.Action.Type)
			}

			result = map[string]interface{}{
				"name":     name,
				"priority": priority,
				"action":   action,
			}

			if rule.Rules == nil || len(*rule.Rules) == 0 {
				continue
			}

			// Determine the rule type based on the first rule's type
			switch (*rule.Rules)[0].(type) {
			case network.ApplicationRule:
				appRules, err := flattenAzureRmFirewallPolicyRuleApplication(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = appRules

				applicationRuleCollection = append(applicationRuleCollection, result)

			case network.Rule:
				networkRules, err := flattenAzureRmFirewallPolicyRuleNetwork(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = networkRules

				networkRuleCollection = append(networkRuleCollection, result)

			default:
				return nil, nil, nil, fmt.Errorf("unknown rule condition type %+v", (*rule.Rules)[0])
			}
		case network.FirewallPolicyNatRuleCollection:
			var name string
			if rule.Name != nil {
				name = *rule.Name
			}
			var priority int32
			if rule.Priority != nil {
				priority = *rule.Priority
			}

			var action string
			if rule.Action != nil {
				action = string(rule.Action.Type)
			}

			rules, err := flattenAzureRmFirewallPolicyRuleNat(rule.Rules)
			if err != nil {
				return nil, nil, nil, err
			}
			result = map[string]interface{}{
				"name":     name,
				"priority": priority,
				"action":   action,
				"rule":     rules,
			}

			natRuleCollection = append(natRuleCollection, result)

		default:
			return nil, nil, nil, fmt.Errorf("unknown rule type %+v", rule)
		}
	}
	return applicationRuleCollection, networkRuleCollection, natRuleCollection, nil
}

func flattenAzureRmFirewallPolicyRuleApplication(input *[]network.BasicFirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(network.ApplicationRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-application rule: %+v", e)
		}

		var name string
		if rule.Name != nil {
			name = *rule.Name
		}

		protocols := make([]interface{}, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				var port int
				if protocol.Port != nil {
					port = int(*protocol.Port)
				}
				protocols = append(protocols, map[string]interface{}{
					"type": string(protocol.ProtocolType),
					"port": port,
				})
			}
		}

		output = append(output, map[string]interface{}{
			"name":                  name,
			"protocols":             protocols,
			"source_addresses":      utils.FlattenStringSlice(rule.SourceAddresses),
			"source_ip_groups":      utils.FlattenStringSlice(rule.SourceIPGroups),
			"destination_fqdns":     utils.FlattenStringSlice(rule.TargetFqdns),
			"destination_fqdn_tags": utils.FlattenStringSlice(rule.FqdnTags),
		})
	}

	return output, nil
}

func flattenAzureRmFirewallPolicyRuleNetwork(input *[]network.BasicFirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(network.Rule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-network rule: %+v", e)
		}

		var name string
		if rule.Name != nil {
			name = *rule.Name
		}

		protocols := make([]interface{}, 0)
		if rule.IPProtocols != nil {
			for _, protocol := range *rule.IPProtocols {
				protocols = append(protocols, string(protocol))
			}
		}

		output = append(output, map[string]interface{}{
			"name":                  name,
			"protocols":             protocols,
			"source_addresses":      utils.FlattenStringSlice(rule.SourceAddresses),
			"source_ip_groups":      utils.FlattenStringSlice(rule.SourceIPGroups),
			"destination_addresses": utils.FlattenStringSlice(rule.DestinationAddresses),
			"destination_ip_groups": utils.FlattenStringSlice(rule.DestinationIPGroups),
			"destination_fqdns":     utils.FlattenStringSlice(rule.DestinationFqdns),
			"destination_ports":     utils.FlattenStringSlice(rule.DestinationPorts),
		})
	}
	return output, nil
}

func flattenAzureRmFirewallPolicyRuleNat(input *[]network.BasicFirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(network.NatRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-nat rule: %+v", e)
		}

		var name string
		if rule.Name != nil {
			name = *rule.Name
		}

		protocols := make([]interface{}, 0)
		if rule.IPProtocols != nil {
			for _, protocol := range *rule.IPProtocols {
				protocols = append(protocols, string(protocol))
			}
		}
		destinationAddr := ""
		if rule.DestinationAddresses != nil && len(*rule.DestinationAddresses) != 0 {
			destinationAddr = (*rule.DestinationAddresses)[0]
		}

		translatedPort := 0
		if rule.TranslatedPort != nil {
			port, err := strconv.Atoi(*rule.TranslatedPort)
			if err != nil {
				return nil, fmt.Errorf(`The "translatedPort" property is not a valid integer (%s)`, *rule.TranslatedPort)
			}
			translatedPort = port
		}

		output = append(output, map[string]interface{}{
			"name":                name,
			"protocols":           protocols,
			"source_addresses":    utils.FlattenStringSlice(rule.SourceAddresses),
			"source_ip_groups":    utils.FlattenStringSlice(rule.SourceIPGroups),
			"destination_address": destinationAddr,
			"destination_ports":   utils.FlattenStringSlice(rule.DestinationPorts),
			"translated_address":  rule.TranslatedAddress,
			"translated_port":     &translatedPort,
		})
	}
	return output, nil
}
