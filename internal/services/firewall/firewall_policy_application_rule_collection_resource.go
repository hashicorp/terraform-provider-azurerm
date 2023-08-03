package firewall

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceFirewallPolicyApplicationRuleCollection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFirewallPolicyApplicationRuleCollectionCreateUpdate,
		Read:   resourceFirewallPolicyApplicationRuleCollectionRead,
		Update: resourceFirewallPolicyApplicationRuleCollectionCreateUpdate,
		Delete: resourceFirewallPolicyApplicationRuleCollectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallPolicyRuleCollectionID(id)
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
				ValidateFunc: validate.FirewallPolicyRuleCollectionName(),
			},

			"rule_collection_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallPolicyRuleCollectionGroupID,
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
					string(network.FirewallPolicyFilterRuleCollectionActionTypeAllow),
					string(network.FirewallPolicyFilterRuleCollectionActionTypeDeny),
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
											string(network.FirewallPolicyRuleApplicationProtocolTypeHTTP),
											string(network.FirewallPolicyRuleApplicationProtocolTypeHTTPS),
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
						"source_addresses": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									validation.IsIPAddress,
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
	}
}

func resourceFirewallPolicyApplicationRuleCollectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	ruleCollectionGroupId, err := parse.FirewallPolicyRuleCollectionGroupID(d.Get("rule_collection_group_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q): %+v", ruleCollectionGroupId.RuleCollectionGroupName, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, err)
		}
	}

	idx := indexFunc(*resp.RuleCollections, func(ruleCollection network.BasicFirewallPolicyRuleCollection) bool {
		info, _ := ruleCollection.AsFirewallPolicyFilterRuleCollection()
		return *info.Name == name
	})

	if d.IsNewResource() && resp.ID != nil && *resp.ID != "" && idx != -1 {
		return tf.ImportAsExistsError("azurerm_firewall_policy_application_rule_collection", *resp.ID)
	}

	locks.ByName(ruleCollectionGroupId.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(ruleCollectionGroupId.FirewallPolicyName, AzureFirewallPolicyResourceName)

	rules := expandFirewallPolicyRuleApplication(d.Get("rule").([]interface{}))
	ruleCollection := &network.FirewallPolicyFilterRuleCollection{
		Action: &network.FirewallPolicyFilterRuleCollectionAction{
			Type: network.FirewallPolicyFilterRuleCollectionActionType(d.Get("action").(string)),
		},
		Name:               utils.String(d.Get("name").(string)),
		Priority:           utils.Int32(int32(d.Get("priority").(int))),
		RuleCollectionType: network.RuleCollectionTypeFirewallPolicyFilterRuleCollection,
		Rules:              rules,
	}

	if idx == -1 {
		ruleCollections := append(*resp.RuleCollections, ruleCollection)
		resp.RuleCollections = &ruleCollections
	} else {
		(*resp.RuleCollections)[idx] = ruleCollection
	}

	param := network.FirewallPolicyRuleCollectionGroup{
		FirewallPolicyRuleCollectionGroupProperties: &network.FirewallPolicyRuleCollectionGroupProperties{
			Priority:        resp.Priority,
			RuleCollections: resp.RuleCollections,
		},
	}

	future, err := client.CreateOrUpdate(ctx, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName, param)
	if err != nil {
		return fmt.Errorf("creating Firewall Policy Rule Collection %q (Resource Group %q / Policy: %q / Rule Collection Group: %q): %+v", name, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q / Rule Collection Group: %q): %+v", name, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName, err)
	}

	resp, err = client.Get(ctx, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName)
	if err != nil {
		return fmt.Errorf("retrieving Firewall Policy Rule Collection Group %q (Resource Group %q / Policy: %q: %+v", ruleCollectionGroupId.RuleCollectionGroupName, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, err)
	}

	idx = indexFunc(*resp.RuleCollections, func(ruleCollection network.BasicFirewallPolicyRuleCollection) bool {
		info, _ := ruleCollection.AsFirewallPolicyFilterRuleCollection()
		return *info.Name == name
	})

	if resp.ID == nil || *resp.ID == "" || idx == -1 {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy Rule Collection %q (Resource Group %q / Policy: %q / Rule Collection Group: %q) ID", name, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName)
	}
	ruleCollectionGroupId, err = parse.FirewallPolicyRuleCollectionGroupID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(parse.NewFirewallPolicyRuleCollectionID(ruleCollectionGroupId.SubscriptionId, ruleCollectionGroupId.ResourceGroup, ruleCollectionGroupId.FirewallPolicyName, ruleCollectionGroupId.RuleCollectionGroupName, name).ID())

	return resourceFirewallPolicyApplicationRuleCollectionRead(d, meta)
}

func resourceFirewallPolicyApplicationRuleCollectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyRuleCollectionID(d.Id())
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

	idx := indexFunc(*resp.RuleCollections, func(policy network.BasicFirewallPolicyRuleCollection) bool {
		info, _ := policy.AsFirewallPolicyFilterRuleCollection()
		return *info.Name == id.RuleCollectionName
	})

	if idx == -1 {
		log.Printf("[DEBUG] Firewall Policy Rule Collection %q was not found in Firewall Policy Rule Collection Group %q - removing from state!", id.RuleCollectionName, id.RuleCollectionGroupName)
		d.SetId("")
		return nil
	}

	ruleCollection, _ := (*resp.RuleCollections)[idx].AsFirewallPolicyFilterRuleCollection()

	d.Set("name", *ruleCollection.Name)
	d.Set("action", ruleCollection.Action.Type)
	d.Set("priority", *ruleCollection.Priority)
	d.Set("rule_collection_group_id", parse.NewFirewallPolicyRuleCollectionGroupID(id.SubscriptionId, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName).ID())

	rules, err := flattenFirewallPolicyRuleApplication(ruleCollection.Rules)
	if err != nil {
		return fmt.Errorf("flattening Firewall Policy Rule Collections: %+v", err)
	}

	if err := d.Set("rule", rules); err != nil {
		return fmt.Errorf("setting `rules`: %+v", err)
	}

	return nil
}

func resourceFirewallPolicyApplicationRuleCollectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallPolicyRuleCollectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q): %+v", id.RuleCollectionGroupName, id.ResourceGroup, id.FirewallPolicyName, err)
		}
	}

	idx := indexFunc(*resp.RuleCollections, func(policy network.BasicFirewallPolicyRuleCollection) bool {
		info, _ := policy.AsFirewallPolicyFilterRuleCollection()
		return *info.Name == id.RuleCollectionName
	})

	if idx == -1 {
		log.Printf("[DEBUG] Firewall Policy Rule Collection %q was not found in Firewall Policy Rule Collection Group %q - removing from state!", id.RuleCollectionName, id.RuleCollectionGroupName)
		return nil
	}

	param := network.FirewallPolicyRuleCollectionGroup{
		FirewallPolicyRuleCollectionGroupProperties: &network.FirewallPolicyRuleCollectionGroupProperties{
			Priority: resp.Priority,
		},
	}
	rulesCollections := *resp.RuleCollections
	rulesCollections = append(rulesCollections[:idx], rulesCollections[idx+1:]...)
	param.FirewallPolicyRuleCollectionGroupProperties.RuleCollections = &rulesCollections

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName, param)
	if err != nil {
		return fmt.Errorf("deleting Firewall Policy Rule Collection %q (Resource Group %q / Policy: %q / Rule Collection Group: %q): %+v", id.RuleCollectionName, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting %q (Resource Group %q / Policy: %q / Rule Collection Group: %q): %+v", id.RuleCollectionName, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName, err)
		}
	}

	return nil
}

func indexFunc[T any](s []T, f func(T) bool) int {
	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			return i
		}
	}
	return -1
}
