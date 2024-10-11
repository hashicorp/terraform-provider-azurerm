// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/firewallpolicyrulecollectiongroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFirewallPolicyRuleCollectionGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyRuleCollectionGroupCreateUpdate,
		Read:   resourceFirewallPolicyRuleCollectionGroupRead,
		Update: resourceFirewallPolicyRuleCollectionGroupCreateUpdate,
		Delete: resourceFirewallPolicyRuleCollectionGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := firewallpolicyrulecollectiongroups.ParseRuleCollectionGroupID(id)
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
				ValidateFunc: validate.FirewallPolicyRuleCollectionGroupName(),
			},

			"firewall_policy_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 65000),
			},

			"application_rule_collection": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionTypeAllow),
								string(firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionTypeDeny),
							}, false),
						},
						"rule": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"description": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"protocols": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"type": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocolTypeHTTP),
														string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocolTypeHTTPS),
														"Mssql",
													}, false),
												},
												"port": {
													Type:         pluginsdk.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(0, 64000),
												},
											},
										},
									},
									"http_headers": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"value": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
									"source_addresses": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsIPv4Range,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_addresses": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsIPv4Range,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"destination_fqdns": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_urls": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_fqdn_tags": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"terminate_tls": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"web_categories": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionTypeAllow),
								string(firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionTypeDeny),
							}, false),
						},
						"rule": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"description": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"protocols": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolAny),
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolTCP),
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolUDP),
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolICMP),
											}, false),
										},
									},
									"source_addresses": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsIPv4Range,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_addresses": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											// Can be IP address, CIDR, "*", or service tag
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_ip_groups": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_fqdns": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_ports": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validate.PortOrPortRangeWithin(1, 65535),
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
								},
							},
						},
					},
				},
			},

			"nat_rule_collection": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"priority": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 65000),
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								// Hardcode to using `Dnat` instead of the one defined in Swagger (i.e. firewallpolicyrulecollectiongroups.DNAT) because of: https://github.com/Azure/azure-rest-api-specs/issues/9986
								// Setting `StateFunc: state.IgnoreCase` will cause other issues, as tracked by: https://github.com/hashicorp/terraform-plugin-sdk/issues/485
								"Dnat",
							}, false),
						},
						"rule": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"description": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"protocols": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolTCP),
												string(firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocolUDP),
											}, false),
										},
									},
									"source_addresses": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsIPAddress,
												validation.IsIPv4Range,
												validation.IsCIDR,
												validation.StringInSlice([]string{`*`}, false),
											),
										},
									},
									"source_ip_groups": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
									"destination_address": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IsIPAddress,
											validation.IsCIDR,
										),
									},
									"destination_ports": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										// only support 1 destination port in one DNAT rule
										MaxItems: 1,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validate.PortOrPortRangeWithin(1, 64000),
										},
									},
									"translated_address": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPAddress,
									},
									"translated_port": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IsPortNumber,
									},
									"translated_fqdn": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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

func resourceFirewallPolicyRuleCollectionGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleCollectionGroups
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	policyId, err := firewallpolicies.ParseFirewallPolicyID(d.Get("firewall_policy_id").(string))
	if err != nil {
		return err
	}

	id := firewallpolicyrulecollectiongroups.NewRuleCollectionGroupID(policyId.SubscriptionId, policyId.ResourceGroupName, policyId.FirewallPolicyName, d.Get("name").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if resp.Model != nil {
			return tf.ImportAsExistsError("azurerm_firewall_policy_rule_collection_group", id.ID())
		}
	}

	locks.ByName(policyId.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(policyId.FirewallPolicyName, AzureFirewallPolicyResourceName)

	param := firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroup{
		Properties: &firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollectionGroupProperties{
			Priority: utils.Int64(int64(d.Get("priority").(int))),
		},
	}
	var rulesCollections []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection
	rulesCollections = append(rulesCollections, expandFirewallPolicyRuleCollectionApplication(d.Get("application_rule_collection").([]interface{}))...)
	rulesCollections = append(rulesCollections, expandFirewallPolicyRuleCollectionNetwork(d.Get("network_rule_collection").([]interface{}))...)

	natRules, err := expandFirewallPolicyRuleCollectionNat(d.Get("nat_rule_collection").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding NAT rule collection: %w", err)
	}
	rulesCollections = append(rulesCollections, natRules...)

	param.Properties.RuleCollections = &rulesCollections

	if err = client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceFirewallPolicyRuleCollectionGroupRead(d, meta)
}

func resourceFirewallPolicyRuleCollectionGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleCollectionGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallpolicyrulecollectiongroups.ParseRuleCollectionGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found- removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RuleCollectionGroupName)
	d.Set("firewall_policy_id", firewallpolicies.NewFirewallPolicyID(id.SubscriptionId, id.ResourceGroupName, id.FirewallPolicyName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("priority", props.Priority)

			applicationRuleCollections, networkRuleCollections, natRuleCollections, err := flattenFirewallPolicyRuleCollection(props.RuleCollections)
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
		}
	}

	return nil
}

func resourceFirewallPolicyRuleCollectionGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleCollectionGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := firewallpolicyrulecollectiongroups.ParseRuleCollectionGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandFirewallPolicyRuleCollectionApplication(input []interface{}) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	return expandFirewallPolicyFilterRuleCollection(input, expandFirewallPolicyRuleApplication)
}

func expandFirewallPolicyRuleCollectionNetwork(input []interface{}) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	return expandFirewallPolicyFilterRuleCollection(input, expandFirewallPolicyRuleNetwork)
}

func expandFirewallPolicyRuleCollectionNat(input []interface{}) ([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, error) {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]interface{})
		rules, err := expandFirewallPolicyRuleNat(rule["rule"].([]interface{}))
		if err != nil {
			return nil, err
		}
		output := &firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollection{
			Name:     utils.String(rule["name"].(string)),
			Priority: utils.Int64(int64(rule["priority"].(int))),
			Action: &firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollectionAction{
				Type: pointer.To(firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollectionActionType(rule["action"].(string))),
			},
			Rules: rules,
		}
		result = append(result, output)
	}
	return result, nil
}

func expandFirewallPolicyFilterRuleCollection(input []interface{}, f func(input []interface{}) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) []firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection, 0)
	for _, e := range input {
		rule := e.(map[string]interface{})
		output := &firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollection{
			Action: &firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionAction{
				Type: pointer.To(firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollectionActionType(rule["action"].(string))),
			},
			Name:     utils.String(rule["name"].(string)),
			Priority: utils.Int64(int64(rule["priority"].(int))),
			Rules:    f(rule["rule"].([]interface{})),
		}
		result = append(result, output)
	}
	return result
}

func expandFirewallPolicyRuleApplication(input []interface{}) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocol
		for _, p := range condition["protocols"].([]interface{}) {
			proto := p.(map[string]interface{})
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocol{
				ProtocolType: pointer.To(firewallpolicyrulecollectiongroups.FirewallPolicyRuleApplicationProtocolType(proto["type"].(string))),
				Port:         utils.Int64(int64(proto["port"].(int))),
			})
		}

		var httpHeader []firewallpolicyrulecollectiongroups.FirewallPolicyHTTPHeaderToInsert
		for _, h := range condition["http_headers"].([]interface{}) {
			header := h.(map[string]interface{})
			httpHeader = append(httpHeader, firewallpolicyrulecollectiongroups.FirewallPolicyHTTPHeaderToInsert{
				HeaderName:  pointer.To(header["name"].(string)),
				HeaderValue: pointer.To(header["value"].(string)),
			})
		}

		output := &firewallpolicyrulecollectiongroups.ApplicationRule{
			Name:                 utils.String(condition["name"].(string)),
			Description:          utils.String(condition["description"].(string)),
			Protocols:            &protocols,
			HTTPHeadersToInsert:  &httpHeader,
			SourceAddresses:      utils.ExpandStringSlice(condition["source_addresses"].([]interface{})),
			SourceIPGroups:       utils.ExpandStringSlice(condition["source_ip_groups"].([]interface{})),
			DestinationAddresses: utils.ExpandStringSlice(condition["destination_addresses"].([]interface{})),
			TargetFqdns:          utils.ExpandStringSlice(condition["destination_fqdns"].([]interface{})),
			TargetURLs:           utils.ExpandStringSlice(condition["destination_urls"].([]interface{})),
			FqdnTags:             utils.ExpandStringSlice(condition["destination_fqdn_tags"].([]interface{})),
			TerminateTLS:         utils.Bool(condition["terminate_tls"].(bool)),
			WebCategories:        utils.ExpandStringSlice(condition["web_categories"].([]interface{})),
		}
		result = append(result, output)
	}
	return &result
}

func expandFirewallPolicyRuleNetwork(input []interface{}) *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].([]interface{}) {
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}
		output := &firewallpolicyrulecollectiongroups.NetworkRule{
			Name:                 utils.String(condition["name"].(string)),
			IPProtocols:          &protocols,
			SourceAddresses:      utils.ExpandStringSlice(condition["source_addresses"].([]interface{})),
			SourceIPGroups:       utils.ExpandStringSlice(condition["source_ip_groups"].([]interface{})),
			DestinationAddresses: utils.ExpandStringSlice(condition["destination_addresses"].([]interface{})),
			DestinationIPGroups:  utils.ExpandStringSlice(condition["destination_ip_groups"].([]interface{})),
			DestinationFqdns:     utils.ExpandStringSlice(condition["destination_fqdns"].([]interface{})),
			DestinationPorts:     utils.ExpandStringSlice(condition["destination_ports"].([]interface{})),
			Description:          pointer.To(condition["description"].(string)),
		}
		result = append(result, output)
	}
	return &result
}

func expandFirewallPolicyRuleNat(input []interface{}) (*[]firewallpolicyrulecollectiongroups.FirewallPolicyRule, error) {
	result := make([]firewallpolicyrulecollectiongroups.FirewallPolicyRule, 0)
	for _, e := range input {
		condition := e.(map[string]interface{})
		var protocols []firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol
		for _, p := range condition["protocols"].([]interface{}) {
			protocols = append(protocols, firewallpolicyrulecollectiongroups.FirewallPolicyRuleNetworkProtocol(p.(string)))
		}
		destinationAddresses := []string{condition["destination_address"].(string)}

		// Exactly one of `translated_address` and `translated_fqdn` should be set.
		if condition["translated_address"].(string) != "" && condition["translated_fqdn"].(string) != "" {
			return nil, fmt.Errorf("can't specify both `translated_address` and `translated_fqdn` in rule %s", condition["name"].(string))
		}
		if condition["translated_address"].(string) == "" && condition["translated_fqdn"].(string) == "" {
			return nil, fmt.Errorf("should specify either `translated_address` or `translated_fqdn` in rule %s", condition["name"].(string))
		}
		output := &firewallpolicyrulecollectiongroups.NatRule{
			Name:                 utils.String(condition["name"].(string)),
			IPProtocols:          &protocols,
			SourceAddresses:      utils.ExpandStringSlice(condition["source_addresses"].([]interface{})),
			SourceIPGroups:       utils.ExpandStringSlice(condition["source_ip_groups"].([]interface{})),
			DestinationAddresses: &destinationAddresses,
			DestinationPorts:     utils.ExpandStringSlice(condition["destination_ports"].([]interface{})),
			TranslatedPort:       utils.String(strconv.Itoa(condition["translated_port"].(int))),
			Description:          pointer.To(condition["description"].(string)),
		}
		if condition["translated_address"].(string) != "" {
			output.TranslatedAddress = utils.String(condition["translated_address"].(string))
		}
		if condition["translated_fqdn"].(string) != "" {
			output.TranslatedFqdn = utils.String(condition["translated_fqdn"].(string))
		}
		result = append(result, output)
	}
	return &result, nil
}

func flattenFirewallPolicyRuleCollection(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRuleCollection) ([]interface{}, []interface{}, []interface{}, error) {
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
		case firewallpolicyrulecollectiongroups.FirewallPolicyFilterRuleCollection:
			var name string
			if rule.Name != nil {
				name = *rule.Name
			}
			var priority int64
			if rule.Priority != nil {
				priority = *rule.Priority
			}

			var action string
			if rule.Action != nil {
				action = string(pointer.From(rule.Action.Type))
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
			case firewallpolicyrulecollectiongroups.ApplicationRule:
				appRules, err := flattenFirewallPolicyRuleApplication(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = appRules

				applicationRuleCollection = append(applicationRuleCollection, result)

			case firewallpolicyrulecollectiongroups.FirewallPolicyRule:
				networkRules, err := flattenFirewallPolicyRuleNetwork(rule.Rules)
				if err != nil {
					return nil, nil, nil, err
				}
				result["rule"] = networkRules

				networkRuleCollection = append(networkRuleCollection, result)

			default:
				return nil, nil, nil, fmt.Errorf("unknown rule condition type %+v", (*rule.Rules)[0])
			}
		case firewallpolicyrulecollectiongroups.FirewallPolicyNatRuleCollection:
			var name string
			if rule.Name != nil {
				name = *rule.Name
			}
			var priority int64
			if rule.Priority != nil {
				priority = *rule.Priority
			}

			var action string
			if rule.Action != nil {
				// todo 4.0 change this from DNAT to Dnat
				// doing this because we hardcode Dnat for https://github.com/Azure/azure-rest-api-specs/issues/9986
				if strings.EqualFold(string(pointer.From(rule.Action.Type)), "Dnat") {
					action = "Dnat"
				} else {
					action = string(pointer.From(rule.Action.Type))
				}
			}

			rules, err := flattenFirewallPolicyRuleNat(rule.Rules)
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

func flattenFirewallPolicyRuleApplication(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.ApplicationRule)
		if !ok {
			return nil, fmt.Errorf("unexpected non-application rule: %+v", e)
		}

		var name string
		if rule.Name != nil {
			name = *rule.Name
		}

		var description string
		if rule.Description != nil {
			description = *rule.Description
		}

		var terminate_tls bool
		if rule.TerminateTLS != nil {
			terminate_tls = *rule.TerminateTLS
		}

		protocols := make([]interface{}, 0)
		if rule.Protocols != nil {
			for _, protocol := range *rule.Protocols {
				var port int
				if protocol.Port != nil {
					port = int(*protocol.Port)
				}
				protocols = append(protocols, map[string]interface{}{
					"type": string(pointer.From(protocol.ProtocolType)),
					"port": port,
				})
			}
		}

		httpHeaders := make([]interface{}, 0)
		for _, header := range pointer.From(rule.HTTPHeadersToInsert) {
			httpHeaders = append(httpHeaders, map[string]interface{}{
				"name":  pointer.From(header.HeaderName),
				"value": pointer.From(header.HeaderValue),
			})
		}

		output = append(output, map[string]interface{}{
			"name":                  name,
			"description":           description,
			"protocols":             protocols,
			"http_headers":          httpHeaders,
			"source_addresses":      utils.FlattenStringSlice(rule.SourceAddresses),
			"source_ip_groups":      utils.FlattenStringSlice(rule.SourceIPGroups),
			"destination_addresses": utils.FlattenStringSlice(rule.DestinationAddresses),
			"destination_urls":      utils.FlattenStringSlice(rule.TargetURLs),
			"destination_fqdns":     utils.FlattenStringSlice(rule.TargetFqdns),
			"destination_fqdn_tags": utils.FlattenStringSlice(rule.FqdnTags),
			"terminate_tls":         terminate_tls,
			"web_categories":        utils.FlattenStringSlice(rule.WebCategories),
		})
	}

	return output, nil
}

func flattenFirewallPolicyRuleNetwork(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.NetworkRule)
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
			"description":           pointer.From(rule.Description),
		})
	}
	return output, nil
}

func flattenFirewallPolicyRuleNat(input *[]firewallpolicyrulecollectiongroups.FirewallPolicyRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}
	output := make([]interface{}, 0)
	for _, e := range *input {
		rule, ok := e.(firewallpolicyrulecollectiongroups.NatRule)
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

		translatedAddress := ""
		if rule.TranslatedAddress != nil {
			translatedAddress = *rule.TranslatedAddress
		}

		translatedFQDN := ""
		if rule.TranslatedFqdn != nil {
			translatedFQDN = *rule.TranslatedFqdn
		}

		output = append(output, map[string]interface{}{
			"name":                name,
			"protocols":           protocols,
			"source_addresses":    utils.FlattenStringSlice(rule.SourceAddresses),
			"source_ip_groups":    utils.FlattenStringSlice(rule.SourceIPGroups),
			"destination_address": destinationAddr,
			"destination_ports":   utils.FlattenStringSlice(rule.DestinationPorts),
			"translated_address":  translatedAddress,
			"translated_port":     translatedPort,
			"translated_fqdn":     translatedFQDN,
			"description":         pointer.From(rule.Description),
		})
	}
	return output, nil
}
