package eventhub

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/authorizationrulesnamespaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/namespaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/networkrulesets"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	eventHubNamespaceResourceName             = "azurerm_eventhub_namespace"
)

func resourceEventHubNamespace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceCreateUpdate,
		Read:   resourceEventHubNamespaceRead,
		Update: resourceEventHubNamespaceCreateUpdate,
		Delete: resourceEventHubNamespaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NamespaceID(id)
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
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(eventhub.Basic),
					string(eventhub.Standard),
				}, true),
			},

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1,
			},

			"auto_inflate_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"dedicated_cluster_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterID,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.SystemAssigned),
							}, false),
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"maximum_throughput_units": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 20),
			},

			"network_rulesets": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.Allow),
								string(eventhub.Deny),
							}, false),
						},

						"trusted_service_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						// 128 limit per https://docs.microsoft.com/azure/event-hubs/event-hubs-quotas
						"virtual_network_rule": {
							Type:       pluginsdk.TypeList,
							Optional:   true,
							MaxItems:   128,
							ConfigMode: pluginsdk.SchemaConfigModeAttr,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{

									// the API returns the subnet ID's resource group name in lowercase
									// https://github.com/Azure/azure-sdk-for-go/issues/5855
									"subnet_id": {
										Type:             pluginsdk.TypeString,
										Required:         true,
										ValidateFunc:     azure.ValidateResourceID,
										DiffSuppressFunc: suppress.CaseDifference,
									},

									"ignore_missing_virtual_network_service_endpoint": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
								},
							},
						},

						// 128 limit per https://docs.microsoft.com/azure/event-hubs/event-hubs-quotas
						"ip_rule": {
							Type:       pluginsdk.TypeList,
							Optional:   true,
							MaxItems:   128,
							ConfigMode: pluginsdk.SchemaConfigModeAttr,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"ip_mask": {
										Type:     pluginsdk.TypeString,
										Required: true,
									},

									"action": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(eventhub.NetworkRuleIPActionAllow),
										ValidateFunc: validation.StringInSlice([]string{
											string(eventhub.NetworkRuleIPActionAllow),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			"default_primary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceEventHubNamespaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace creation.")

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]interface{})
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)
	zoneRedundant := d.Get("zone_redundant").(bool)

	parameters := namespaces.EHNamespace{
		Location: &location,
		Sku: &namespaces.Sku{
			Name: namespaces.SkuName(sku),
			Tier: func() *namespaces.SkuTier {
				v := namespaces.SkuTier(sku)
				return &v
			}(),
			Capacity: utils.Int64(int64(capacity)),
		},
		Identity: expandEventHubIdentity(d.Get("identity").([]interface{})),
		Properties: &namespaces.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			ZoneRedundant:        utils.Bool(zoneRedundant),
		},
		Tags: expandTags(t),
	}

	if v := d.Get("dedicated_cluster_id").(string); v != "" {
		parameters.Properties.ClusterArmId = utils.String(v)
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.Properties.MaximumThroughputUnits = utils.Int64(int64(v.(int)))
	}

	// @favoretti: if we are downgrading from Standard to Basic SKU and namespace had both autoInflate enabled and
	// maximumThroughputUnits set - we need to force throughput units back to 0, otherwise downgrade fails
	//
	// See: https://github.com/terraform-providers/terraform-provider-azurerm/issues/10244
	//
	if *parameters.Sku.Tier == namespaces.SkuTierBasic && !autoInflateEnabled {
		parameters.Properties.MaximumThroughputUnits = utils.Int64(0)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		rulesets := networkrulesets.NetworkRuleSet{
			Properties: expandEventHubNamespaceNetworkRuleset(ruleSets.([]interface{})),
		}

		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name != namespaces.SkuNameBasic {
			ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
			namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroup, id.Name)
			if _, err := ruleSetsClient.NamespacesCreateOrUpdateNetworkRuleSet(ctx, namespaceId, rulesets); err != nil {
				return fmt.Errorf("setting network ruleset properties for %s: %+v", id, err)
			}
		} else if rulesets.Properties != nil {
			props := rulesets.Properties
			// so if the user has specified the non default rule sets throw a validation error
			if *props.DefaultAction != networkrulesets.DefaultActionDeny ||
				(props.IpRules != nil && len(*props.IpRules) > 0) ||
				(props.VirtualNetworkRules != nil && len(*props.VirtualNetworkRules) > 0) {
				return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
			}
		}
	}

	return resourceEventHubNamespaceRead(d, meta)
}

func resourceEventHubNamespaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	authorizationKeysClient := meta.(*clients.Client).Eventhub.NamespaceAuthorizationRulesClient
	ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
			d.Set("capacity", sku.Capacity)
		}

		if err := d.Set("identity", flattenEventHubIdentity(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
			d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
			d.Set("zone_redundant", props.ZoneRedundant)
			d.Set("dedicated_cluster_id", props.ClusterArmId)
		}

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroup, id.Name)
	ruleset, err := ruleSetsClient.NamespacesGetNetworkRuleSet(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving Network Rule Sets for %s: %+v", *id, err)
	}

	if err := d.Set("network_rulesets", flattenEventHubNamespaceNetworkRuleset(ruleset)); err != nil {
		return fmt.Errorf("Error setting `network_ruleset` for Evenhub Namespace %s: %v", id.Name, err)
	}

	authorizationRuleId := authorizationrulesnamespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroup, id.Name, eventHubNamespaceDefaultAuthorizationRule)
	keys, err := authorizationKeysClient.NamespacesListKeys(ctx, authorizationRuleId)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q: %+v", id.Name, err)
	}

	if model := keys.Model; model != nil {
		d.Set("default_primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string_alias", model.AliasSecondaryConnectionString)
		d.Set("default_primary_connection_string", model.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", model.SecondaryConnectionString)
		d.Set("default_primary_key", model.PrimaryKey)
		d.Set("default_secondary_key", model.SecondaryKey)
	}

	return nil
}

func resourceEventHubNamespaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(future.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return waitForEventHubNamespaceToBeDeleted(ctx, client, *id)
}

func waitForEventHubNamespaceToBeDeleted(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: eventHubNamespaceStateStatusCodeRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func eventHubNamespaceStateStatusCodeRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if res.HttpResponse != nil {
			log.Printf("Retrieving %s returned Status %d", id, res.HttpResponse.StatusCode)
		}

		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}

func expandEventHubNamespaceNetworkRuleset(input []interface{}) *networkrulesets.NetworkRuleSetProperties {
	if len(input) == 0 {
		return nil
	}

	block := input[0].(map[string]interface{})

	ruleset := networkrulesets.NetworkRuleSetProperties{
		DefaultAction: func() *networkrulesets.DefaultAction {
			v := networkrulesets.DefaultAction(block["default_action"].(string))
			return &v
		}(),
	}

	if v, ok := block["trusted_service_access_enabled"]; ok {
		ruleset.TrustedServiceAccessEnabled = utils.Bool(v.(bool))
	}

	if v, ok := block["virtual_network_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []networkrulesets.NWRuleSetVirtualNetworkRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, networkrulesets.NWRuleSetVirtualNetworkRules{
					Subnet: &networkrulesets.Subnet{
						Id: utils.String(rblock["subnet_id"].(string)),
					},
					IgnoreMissingVnetServiceEndpoint: utils.Bool(rblock["ignore_missing_virtual_network_service_endpoint"].(bool)),
				})
			}

			ruleset.VirtualNetworkRules = &rules
		}
	}

	if v, ok := block["ip_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []networkrulesets.NWRuleSetIpRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, networkrulesets.NWRuleSetIpRules{
					IpMask: utils.String(rblock["ip_mask"].(string)),
					Action: func() *networkrulesets.NetworkRuleIPAction {
						v := networkrulesets.NetworkRuleIPAction(rblock["action"].(string))
						return &v
					}(),
				})
			}

			ruleset.IpRules = &rules
		}
	}

	return &ruleset
}

func flattenEventHubNamespaceNetworkRuleset(ruleset networkrulesets.NamespacesGetNetworkRuleSetResponse) []interface{} {
	if ruleset.Model == nil || ruleset.Model.Properties == nil {
		return nil
	}

	vnetBlocks := make([]interface{}, 0)
	if vnetRules := ruleset.Model.Properties.VirtualNetworkRules; vnetRules != nil {
		for _, vnetRule := range *vnetRules {
			block := make(map[string]interface{})

			if s := vnetRule.Subnet; s != nil {
				if v := s.Id; v != nil {
					block["subnet_id"] = *v
				}
			}

			if v := vnetRule.IgnoreMissingVnetServiceEndpoint; v != nil {
				block["ignore_missing_virtual_network_service_endpoint"] = *v
			}

			vnetBlocks = append(vnetBlocks, block)
		}
	}
	ipBlocks := make([]interface{}, 0)
	if ipRules := ruleset.Model.Properties.IpRules; ipRules != nil {
		for _, ipRule := range *ipRules {
			block := make(map[string]interface{})

			action := ""
			if ipRule.Action != nil {
				action = string(*ipRule.Action)
			}

			block["action"] = action

			if v := ipRule.IpMask; v != nil {
				block["ip_mask"] = *v
			}

			ipBlocks = append(ipBlocks, block)
		}
	}

	// TODO: fix this

	return []interface{}{map[string]interface{}{
		"default_action":                 string(*ruleset.Model.Properties.DefaultAction),
		"virtual_network_rule":           vnetBlocks,
		"ip_rule":                        ipBlocks,
		"trusted_service_access_enabled": ruleset.Model.Properties.TrustedServiceAccessEnabled,
	}}
}

func expandEventHubIdentity(input []interface{}) *namespaces.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &namespaces.Identity{
		Type: func() *namespaces.ManagedServiceIdentityType {
			v := namespaces.ManagedServiceIdentityType(v["type"].(string))
			return &v
		}(),
	}
}

func flattenEventHubIdentity(input *namespaces.Identity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	identityType := ""
	if input.Type != nil {
		identityType = string(*input.Type)
	}

	principalID := ""
	if input.PrincipalId != nil {
		principalID = *input.PrincipalId
	}

	tenantID := ""
	if input.TenantId != nil {
		tenantID = *input.TenantId
	}

	return []interface{}{
		map[string]interface{}{
			"type":         identityType,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
