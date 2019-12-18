package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmWebApplicationFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmWebApplicationFirewallPolicyCreateUpdate,
		Read:   resourceArmWebApplicationFirewallPolicyRead,
		Update: resourceArmWebApplicationFirewallPolicyCreateUpdate,
		Delete: resourceArmWebApplicationFirewallPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"custom_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.WebApplicationFirewallActionAllow),
								string(network.WebApplicationFirewallActionBlock),
								string(network.WebApplicationFirewallActionLog),
							}, false),
						},
						"match_conditions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"match_variables": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"variable_name": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(network.RemoteAddr),
														string(network.RequestMethod),
														string(network.QueryString),
														string(network.PostArgs),
														string(network.RequestURI),
														string(network.RequestHeaders),
														string(network.RequestBody),
														string(network.RequestCookies),
													}, false),
												},
												"selector": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.WebApplicationFirewallOperatorIPMatch),
											string(network.WebApplicationFirewallOperatorEqual),
											string(network.WebApplicationFirewallOperatorContains),
											string(network.WebApplicationFirewallOperatorLessThan),
											string(network.WebApplicationFirewallOperatorGreaterThan),
											string(network.WebApplicationFirewallOperatorLessThanOrEqual),
											string(network.WebApplicationFirewallOperatorGreaterThanOrEqual),
											string(network.WebApplicationFirewallOperatorBeginsWith),
											string(network.WebApplicationFirewallOperatorEndsWith),
											string(network.WebApplicationFirewallOperatorRegex),
										}, false),
									},
									"negation_condition": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.WebApplicationFirewallRuleTypeMatchRule),
								string(network.WebApplicationFirewallRuleTypeInvalid),
							}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"policy_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Prevention),
								string(network.Detection),
							}, false),
							Default: string(network.Prevention),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmWebApplicationFirewallPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_web_application_firewall_policy", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	customRules := d.Get("custom_rules").([]interface{})
	policySettings := d.Get("policy_settings").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.WebApplicationFirewallPolicy{
		Location: utils.String(location),
		WebApplicationFirewallPolicyPropertiesFormat: &network.WebApplicationFirewallPolicyPropertiesFormat{
			CustomRules:    expandArmWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(customRules),
			PolicySettings: expandArmWebApplicationFirewallPolicyPolicySettings(policySettings),
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Web Application Firewall Policy %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmWebApplicationFirewallPolicyRead(d, meta)
}

func resourceArmWebApplicationFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ApplicationGatewayWebApplicationFirewallPolicies"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Web Application Firewall Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if webApplicationFirewallPolicyPropertiesFormat := resp.WebApplicationFirewallPolicyPropertiesFormat; webApplicationFirewallPolicyPropertiesFormat != nil {
		if err := d.Set("custom_rules", flattenArmWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(webApplicationFirewallPolicyPropertiesFormat.CustomRules)); err != nil {
			return fmt.Errorf("Error setting `custom_rules`: %+v", err)
		}
		if err := d.Set("policy_settings", flattenArmWebApplicationFirewallPolicyPolicySettings(webApplicationFirewallPolicyPropertiesFormat.PolicySettings)); err != nil {
			return fmt.Errorf("Error setting `policy_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmWebApplicationFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.WebApplicationFirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ApplicationGatewayWebApplicationFirewallPolicies"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Web Application Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input []interface{}) *[]network.WebApplicationFirewallCustomRule {
	results := make([]network.WebApplicationFirewallCustomRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		priority := v["priority"].(int)
		ruleType := v["rule_type"].(string)
		matchConditions := v["match_conditions"].([]interface{})
		action := v["action"].(string)

		result := network.WebApplicationFirewallCustomRule{
			Action:          network.WebApplicationFirewallAction(action),
			MatchConditions: expandArmWebApplicationFirewallPolicyMatchCondition(matchConditions),
			Name:            utils.String(name),
			Priority:        utils.Int32(int32(priority)),
			RuleType:        network.WebApplicationFirewallRuleType(ruleType),
		}

		results = append(results, result)
	}
	return &results
}

func expandArmWebApplicationFirewallPolicyPolicySettings(input []interface{}) *network.PolicySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	enabled := network.WebApplicationFirewallEnabledStateDisabled
	if v["enabled"].(bool) {
		enabled = network.WebApplicationFirewallEnabledStateEnabled
	}
	mode := v["mode"].(string)

	result := network.PolicySettings{
		State: enabled,
		Mode:  network.WebApplicationFirewallMode(mode),
	}
	return &result
}

func expandArmWebApplicationFirewallPolicyMatchCondition(input []interface{}) *[]network.MatchCondition {
	results := make([]network.MatchCondition, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		matchVariables := v["match_variables"].([]interface{})
		operator := v["operator"].(string)
		negationCondition := v["negation_condition"].(bool)
		matchValues := v["match_values"].([]interface{})

		result := network.MatchCondition{
			MatchValues:      utils.ExpandStringSlice(matchValues),
			MatchVariables:   expandArmWebApplicationFirewallPolicyMatchVariable(matchVariables),
			NegationConditon: utils.Bool(negationCondition),
			Operator:         network.WebApplicationFirewallOperator(operator),
		}

		results = append(results, result)
	}
	return &results
}

func expandArmWebApplicationFirewallPolicyMatchVariable(input []interface{}) *[]network.MatchVariable {
	results := make([]network.MatchVariable, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		variableName := v["variable_name"].(string)
		selector := v["selector"].(string)

		result := network.MatchVariable{
			Selector:     utils.String(selector),
			VariableName: network.WebApplicationFirewallMatchVariable(variableName),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input *[]network.WebApplicationFirewallCustomRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		v["action"] = string(item.Action)
		v["match_conditions"] = flattenArmWebApplicationFirewallPolicyMatchCondition(item.MatchConditions)
		if priority := item.Priority; priority != nil {
			v["priority"] = int(*priority)
		}
		v["rule_type"] = string(item.RuleType)

		results = append(results, v)
	}

	return results
}

func flattenArmWebApplicationFirewallPolicyPolicySettings(input *network.PolicySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["enabled"] = input.State == network.WebApplicationFirewallEnabledStateDisabled
	result["mode"] = string(input.Mode)

	return []interface{}{result}
}

func flattenArmWebApplicationFirewallPolicyMatchCondition(input *[]network.MatchCondition) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["match_values"] = utils.FlattenStringSlice(item.MatchValues)
		v["match_variables"] = flattenArmWebApplicationFirewallPolicyMatchVariable(item.MatchVariables)
		if negationCondition := item.NegationConditon; negationCondition != nil {
			v["negation_condition"] = *negationCondition
		}
		v["operator"] = string(item.Operator)

		results = append(results, v)
	}

	return results
}

func flattenArmWebApplicationFirewallPolicyMatchVariable(input *[]network.MatchVariable) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if selector := item.Selector; selector != nil {
			v["selector"] = *selector
		}
		v["variable_name"] = string(item.VariableName)

		results = append(results, v)
	}

	return results
}
