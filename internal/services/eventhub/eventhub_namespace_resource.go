// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	eventHubNamespaceResourceName             = "azurerm_eventhub_namespace"
)

func resourceEventHubNamespace() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceEventHubNamespaceCreate,
		Read:   resourceEventHubNamespaceRead,
		Update: resourceEventHubNamespaceUpdate,
		Delete: resourceEventHubNamespaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namespaces.ParseNamespaceID(id)
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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(namespaces.SkuNameBasic),
					string(namespaces.SkuNameStandard),
					string(namespaces.SkuNamePremium),
				}, false),
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

			"dedicated_cluster_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: eventhubsclusters.ValidateClusterID,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"maximum_throughput_units": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 40),
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
								string(networkrulesets.DefaultActionAllow),
								string(networkrulesets.DefaultActionDeny),
							}, false),
						},

						"public_network_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"trusted_service_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						// Returned value of the `virtual_network_rule` array does not honor the input order,
						// possibly a service design, thus changed to TypeSet
						"virtual_network_rule": {
							Type:       pluginsdk.TypeSet,
							Optional:   true,
							ConfigMode: pluginsdk.SchemaConfigModeAttr,
							Set:        resourceVnetRuleHash,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									// the API returns the subnet ID's resource group name in lowercase
									// https://github.com/Azure/azure-sdk-for-go/issues/5855
									"subnet_id": {
										Type:             pluginsdk.TypeString,
										Required:         true,
										ValidateFunc:     commonids.ValidateSubnetID,
										DiffSuppressFunc: suppress.CaseDifference,
									},

									"ignore_missing_virtual_network_service_endpoint": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
								},
							},
						},

						"ip_rule": {
							Type:       pluginsdk.TypeList,
							Optional:   true,
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
										Default:  string(networkrulesets.NetworkRuleIPActionAllow),
										ValidateFunc: validation.StringInSlice([]string{
											string(networkrulesets.NetworkRuleIPActionAllow),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(namespaces.TlsVersionOnePointTwo),
				ValidateFunc: validation.StringInSlice([]string{
					string(namespaces.TlsVersionOnePointTwo),
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				oldSku, newSku := d.GetChange("sku")
				if d.HasChange("sku") {
					if strings.EqualFold(newSku.(string), string(namespaces.SkuNamePremium)) || strings.EqualFold(oldSku.(string), string(namespaces.SkuTierPremium)) {
						log.Printf("[DEBUG] cannot migrate a namespace from or to Premium SKU")
						d.ForceNew("sku")
					}
				}
				return nil
			}),
			pluginsdk.CustomizeDiffShim(eventhubTLSVersionDiff),
		),
	}

	if !features.FivePointOh() {
		resource.Schema["minimum_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(namespaces.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.TlsVersionOnePointZero),
				string(namespaces.TlsVersionOnePointOne),
				string(namespaces.TlsVersionOnePointTwo),
			}, false),
		}
	}
	return resource
}

func resourceEventHubNamespaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace creation.")

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if existing.Model != nil {
		return tf.ImportAsExistsError("azurerm_eventhub_namespace", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]interface{})
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
	}

	disableLocalAuth := false
	if !d.Get("local_authentication_enabled").(bool) {
		disableLocalAuth = true
	}

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
		Identity: identity,
		Properties: &namespaces.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			DisableLocalAuth:     utils.Bool(disableLocalAuth),
			PublicNetworkAccess:  &publicNetworkEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v := d.Get("dedicated_cluster_id").(string); v != "" {
		parameters.Properties.ClusterArmId = utils.String(v)
	}

	if tlsValue := d.Get("minimum_tls_version").(string); tlsValue != "" {
		minimumTls := namespaces.TlsVersion(tlsValue)
		parameters.Properties.MinimumTlsVersion = &minimumTls
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.Properties.MaximumThroughputUnits = utils.Int64(int64(v.(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name == namespaces.SkuNameBasic {
			return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
		}

		rulesets := networkrulesets.NetworkRuleSet{
			Properties: expandEventHubNamespaceNetworkRuleset(ruleSets.([]interface{})),
		}

		if !strings.EqualFold(string(*rulesets.Properties.PublicNetworkAccess), string(*parameters.Properties.PublicNetworkAccess)) {
			return fmt.Errorf("the value of public network access of namespace should be the same as of the network rulesets")
		}

		ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
		namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
		if _, err := ruleSetsClient.NamespacesCreateOrUpdateNetworkRuleSet(ctx, namespaceId, rulesets); err != nil {
			return fmt.Errorf("setting network ruleset properties for %s: %+v", id, err)
		}
	}

	return resourceEventHubNamespaceRead(d, meta)
}

func resourceEventHubNamespaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace update.")

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]interface{})
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)

	publicNetworkEnabled := namespaces.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkEnabled = namespaces.PublicNetworkAccessDisabled
	}

	disableLocalAuth := false
	if !d.Get("local_authentication_enabled").(bool) {
		disableLocalAuth = true
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

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
		Identity: identity,
		Properties: &namespaces.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			DisableLocalAuth:     utils.Bool(disableLocalAuth),
			PublicNetworkAccess:  &publicNetworkEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v := d.Get("dedicated_cluster_id").(string); v != "" {
		parameters.Properties.ClusterArmId = utils.String(v)
	}

	if tlsValue := d.Get("minimum_tls_version").(string); tlsValue != "" {
		minimumTls := namespaces.TlsVersion(tlsValue)
		parameters.Properties.MinimumTlsVersion = &minimumTls
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.Properties.MaximumThroughputUnits = utils.Int64(int64(v.(int)))
	}

	// @favoretti: if we are downgrading from Standard to Basic SKU and namespace had both autoInflate enabled and
	// maximumThroughputUnits set - we need to force throughput units back to 0, otherwise downgrade fails
	//
	// See: https://github.com/hashicorp/terraform-provider-azurerm/issues/10244
	//
	if *parameters.Sku.Tier == namespaces.SkuTierBasic && !autoInflateEnabled {
		parameters.Properties.MaximumThroughputUnits = utils.Int64(0)
	}

	if _, err = client.Update(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	deadline, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"Activating", "ActivatingIdentity", "Updating", "Pending"},
		Target:       []string{"Succeeded"},
		Refresh:      eventHubNamespaceProvisioningStateRefreshFunc(ctx, client, id),
		Timeout:      time.Until(deadline),
		PollInterval: 10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be updated: %+v", id, err)
	}

	if d.HasChange("network_rulesets") {
		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name == namespaces.SkuNameBasic {
			return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
		}

		ruleSets := d.Get("network_rulesets")
		rulesets := networkrulesets.NetworkRuleSet{
			Properties: expandEventHubNamespaceNetworkRuleset(ruleSets.([]interface{})),
		}

		if !strings.EqualFold(string(*rulesets.Properties.PublicNetworkAccess), string(*parameters.Properties.PublicNetworkAccess)) {
			return fmt.Errorf("the value of public network access of namespace should be the same as of the network rulesets")
		}

		ruleSetsClient := meta.(*clients.Client).Eventhub.NetworkRuleSetsClient
		namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
		if _, err := ruleSetsClient.NamespacesCreateOrUpdateNetworkRuleSet(ctx, namespaceId, rulesets); err != nil {
			return fmt.Errorf("setting network ruleset properties for %s: %+v", id, err)
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

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(sku.Name))
			d.Set("capacity", sku.Capacity)
		}

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
			d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
			d.Set("dedicated_cluster_id", props.ClusterArmId)

			localAuthDisabled := false
			if props.DisableLocalAuth != nil {
				localAuthDisabled = *props.DisableLocalAuth
			}
			d.Set("local_authentication_enabled", !localAuthDisabled)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == namespaces.PublicNetworkAccessDisabled {
				publicNetworkAccess = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("minimum_tls_version", string(pointer.From(props.MinimumTlsVersion)))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	namespaceId := networkrulesets.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	ruleset, err := ruleSetsClient.NamespacesGetNetworkRuleSet(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving Network Rule Sets for %s: %+v", *id, err)
	}

	networkRuleSets, err := flattenEventHubNamespaceNetworkRuleset(ruleset)
	if err != nil {
		return fmt.Errorf("flattening `network_rule` for %s: %+v", id, err)
	}
	if err := d.Set("network_rulesets", networkRuleSets); err != nil {
		return fmt.Errorf("setting `network_ruleset` for Evenhub Namespace %s: %v", id.NamespaceName, err)
	}

	authorizationRuleId := authorizationrulesnamespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, eventHubNamespaceDefaultAuthorizationRule)
	keys, err := authorizationKeysClient.NamespacesListKeys(ctx, authorizationRuleId)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q: %+v", id.NamespaceName, err)
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

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandEventHubNamespaceNetworkRuleset(input []interface{}) *networkrulesets.NetworkRuleSetProperties {
	if len(input) == 0 {
		return nil
	}

	block := input[0].(map[string]interface{})

	publicNetworkAccess := networkrulesets.PublicNetworkAccessFlagEnabled
	if !block["public_network_access_enabled"].(bool) {
		publicNetworkAccess = networkrulesets.PublicNetworkAccessFlagDisabled
	}

	ruleset := networkrulesets.NetworkRuleSetProperties{
		DefaultAction: func() *networkrulesets.DefaultAction {
			v := networkrulesets.DefaultAction(block["default_action"].(string))
			return &v
		}(),
		PublicNetworkAccess: &publicNetworkAccess,
	}

	if v, ok := block["trusted_service_access_enabled"]; ok {
		ruleset.TrustedServiceAccessEnabled = utils.Bool(v.(bool))
	}

	if v, ok := block["virtual_network_rule"]; ok {
		value := v.(*pluginsdk.Set).List()
		if len(value) > 0 {
			var rules []networkrulesets.NWRuleSetVirtualNetworkRules
			for _, r := range value {
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
			var rules []networkrulesets.NWRuleSetIPRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, networkrulesets.NWRuleSetIPRules{
					IPMask: utils.String(rblock["ip_mask"].(string)),
					Action: func() *networkrulesets.NetworkRuleIPAction {
						v := networkrulesets.NetworkRuleIPAction(rblock["action"].(string))
						return &v
					}(),
				})
			}

			ruleset.IPRules = &rules
		}
	}

	return &ruleset
}

func flattenEventHubNamespaceNetworkRuleset(ruleset networkrulesets.NamespacesGetNetworkRuleSetOperationResponse) ([]interface{}, error) {
	if ruleset.Model == nil || ruleset.Model.Properties == nil {
		return nil, nil
	}

	vnetBlocks := make([]interface{}, 0)
	if vnetRules := ruleset.Model.Properties.VirtualNetworkRules; vnetRules != nil {
		for _, vnetRule := range *vnetRules {
			block := make(map[string]interface{})

			if s := vnetRule.Subnet; s != nil {
				if v := s.Id; v != nil {
					// the API returns the subnet ID's resource group name in lowercase
					// https://github.com/Azure/azure-sdk-for-go/issues/5855
					// for some reason the DiffSuppressFunc for `subnet_id` isn't working as intended, so we'll also flatten the id insensitively
					subnetId, err := commonids.ParseSubnetIDInsensitively(*v)
					if err != nil {
						return nil, fmt.Errorf("parsing `subnet_id`: %+v", err)
					}
					block["subnet_id"] = subnetId.ID()
				}
			}

			if v := vnetRule.IgnoreMissingVnetServiceEndpoint; v != nil {
				block["ignore_missing_virtual_network_service_endpoint"] = *v
			}

			vnetBlocks = append(vnetBlocks, block)
		}
	}
	ipBlocks := make([]interface{}, 0)
	if ipRules := ruleset.Model.Properties.IPRules; ipRules != nil {
		for _, ipRule := range *ipRules {
			block := make(map[string]interface{})

			action := ""
			if ipRule.Action != nil {
				action = string(*ipRule.Action)
			}

			block["action"] = action

			if v := ipRule.IPMask; v != nil {
				block["ip_mask"] = *v
			}

			ipBlocks = append(ipBlocks, block)
		}
	}

	// TODO: fix this

	publicNetworkAccess := true
	if ruleset.Model.Properties.PublicNetworkAccess != nil && *ruleset.Model.Properties.PublicNetworkAccess == networkrulesets.PublicNetworkAccessFlagDisabled {
		publicNetworkAccess = false
	}
	return []interface{}{map[string]interface{}{
		"default_action":                 string(*ruleset.Model.Properties.DefaultAction),
		"public_network_access_enabled":  publicNetworkAccess,
		"virtual_network_rule":           vnetBlocks,
		"ip_rule":                        ipBlocks,
		"trusted_service_access_enabled": ruleset.Model.Properties.TrustedServiceAccessEnabled,
	}}, nil
}

// The resource id of subnet_id that's being returned by API is always lower case &
// the default caseDiff suppress func is not working in TypeSet
func resourceVnetRuleHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["subnet_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
		}
		if v, ok := m["ignore_missing_virtual_network_service_endpoint"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}
	return pluginsdk.HashString(buf.String())
}

func eventhubTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	if old != "" && new == "" {
		err = fmt.Errorf("`minimum_tls_version` has been set before, please set a valid value for this property ")
	}
	return
}

func eventHubNamespaceProvisioningStateRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)

		provisioningState := "Pending"
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, provisioningState, nil
			}
			return nil, "Error", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			provisioningState = *res.Model.Properties.ProvisioningState
		}

		return res, provisioningState, nil
	}
}
