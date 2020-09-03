package eventhub

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
var eventHubNamespaceResourceName = "azurerm_eventhub_namespace"

func resourceArmEventHubNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubNamespaceCreateUpdate,
		Read:   resourceArmEventHubNamespaceRead,
		Update: resourceArmEventHubNamespaceCreateUpdate,
		Delete: resourceArmEventHubNamespaceDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NamespaceID(id)
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
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(eventhub.Basic),
					string(eventhub.Standard),
				}, true),
			},

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"auto_inflate_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"dedicated_cluster_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubDedicatedClusterID,
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.SystemAssigned),
							}, false),
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"maximum_throughput_units": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 20),
			},

			"network_rulesets": {
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"default_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.Allow),
								string(eventhub.Deny),
							}, false),
						},

						// 128 limit per https://docs.microsoft.com/azure/event-hubs/event-hubs-quotas
						"virtual_network_rule": {
							Type:       schema.TypeList,
							Optional:   true,
							MaxItems:   128,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									// the API returns the subnet ID's resource group name in lowercase
									// https://github.com/Azure/azure-sdk-for-go/issues/5855
									"subnet_id": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateFunc:     azure.ValidateResourceID,
										DiffSuppressFunc: suppress.CaseDifference,
									},

									"ignore_missing_virtual_network_service_endpoint": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},

						// 128 limit per https://docs.microsoft.com/azure/event-hubs/event-hubs-quotas
						"ip_rule": {
							Type:       schema.TypeList,
							Optional:   true,
							MaxItems:   128,
							ConfigMode: schema.SchemaConfigModeAttr,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_mask": {
										Type:     schema.TypeString,
										Required: true,
									},

									"action": {
										Type:     schema.TypeString,
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
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string_alias": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmEventHubNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	t := d.Get("tags").(map[string]interface{})
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)
	zoneRedundant := d.Get("zone_redundant").(bool)

	parameters := eventhub.EHNamespace{
		Location: &location,
		Sku: &eventhub.Sku{
			Name:     eventhub.SkuName(sku),
			Tier:     eventhub.SkuTier(sku),
			Capacity: &capacity,
		},
		Identity: expandEventHubIdentity(d.Get("identity").([]interface{})),
		EHNamespaceProperties: &eventhub.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			ZoneRedundant:        utils.Bool(zoneRedundant),
		},
		Tags: tags.Expand(t),
	}

	if v := d.Get("dedicated_cluster_id").(string); v != "" {
		parameters.EHNamespaceProperties.ClusterArmID = utils.String(v)
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.EHNamespaceProperties.MaximumThroughputUnits = utils.Int32(int32(v.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating eventhub namespace: %+v", err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		rulesets := eventhub.NetworkRuleSet{
			NetworkRuleSetProperties: expandEventHubNamespaceNetworkRuleset(ruleSets.([]interface{})),
		}

		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name != eventhub.Basic {
			if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resGroup, name, rulesets); err != nil {
				return fmt.Errorf("Error setting network ruleset properties for EventHub Namespace %q (resource group %q): %v", name, resGroup, err)
			}
		} else {
			// so if the user has specified the non default rule sets throw a validation error
			if rulesets.DefaultAction != eventhub.Deny ||
				(rulesets.IPRules != nil && len(*rulesets.IPRules) > 0) ||
				(rulesets.VirtualNetworkRules != nil && len(*rulesets.VirtualNetworkRules) > 0) {
				return fmt.Errorf("network_rulesets cannot be used when the SKU is basic")
			}
		}
	}

	return resourceArmEventHubNamespaceRead(d, meta)
}

func resourceArmEventHubNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on EventHub Namespace %q: %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
		d.Set("capacity", sku.Capacity)
	}

	if err := d.Set("identity", flattenEventHubIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if props := resp.EHNamespaceProperties; props != nil {
		d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
		d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
		d.Set("zone_redundant", props.ZoneRedundant)
		d.Set("dedicated_cluster_id", props.ClusterArmID)
	}

	ruleset, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on EventHub Namespace %q Network Ruleset: %+v", id.Name, err)
	}

	if err := d.Set("network_rulesets", flattenEventHubNamespaceNetworkRuleset(ruleset)); err != nil {
		return fmt.Errorf("Error setting `network_ruleset` for Evenhub Namespace %s: %v", id.Name, err)
	}

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.Name, eventHubNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q: %+v", id.Name, err)
	} else {
		d.Set("default_primary_connection_string_alias", keys.AliasPrimaryConnectionString)
		d.Set("default_secondary_connection_string_alias", keys.AliasSecondaryConnectionString)
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmEventHubNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request of EventHub Namespace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return waitForEventHubNamespaceToBeDeleted(ctx, client, id.ResourceGroup, id.Name, d)
}

func waitForEventHubNamespaceToBeDeleted(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string, d *schema.ResourceData) error {
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for EventHub Namespace (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: eventHubNamespaceStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for EventHub NameSpace (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func eventHubNamespaceStateStatusCodeRefreshFunc(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving EventHub Namespace %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the EventHub Namespace %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandEventHubNamespaceNetworkRuleset(input []interface{}) *eventhub.NetworkRuleSetProperties {
	if len(input) == 0 {
		return nil
	}

	block := input[0].(map[string]interface{})

	ruleset := eventhub.NetworkRuleSetProperties{
		DefaultAction: eventhub.DefaultAction(block["default_action"].(string)),
	}

	if v, ok := block["virtual_network_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []eventhub.NWRuleSetVirtualNetworkRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, eventhub.NWRuleSetVirtualNetworkRules{
					Subnet: &eventhub.Subnet{
						ID: utils.String(rblock["subnet_id"].(string)),
					},
					IgnoreMissingVnetServiceEndpoint: utils.Bool(rblock["ignore_missing_virtual_network_service_endpoint"].(bool)),
				})
			}

			ruleset.VirtualNetworkRules = &rules
		}
	}

	if v, ok := block["ip_rule"].([]interface{}); ok {
		if len(v) > 0 {
			var rules []eventhub.NWRuleSetIPRules
			for _, r := range v {
				rblock := r.(map[string]interface{})
				rules = append(rules, eventhub.NWRuleSetIPRules{
					IPMask: utils.String(rblock["ip_mask"].(string)),
					Action: eventhub.NetworkRuleIPAction(rblock["action"].(string)),
				})
			}

			ruleset.IPRules = &rules
		}
	}

	return &ruleset
}

func flattenEventHubNamespaceNetworkRuleset(ruleset eventhub.NetworkRuleSet) []interface{} {
	if ruleset.NetworkRuleSetProperties == nil {
		return nil
	}

	vnetBlocks := make([]interface{}, 0)
	if vnetRules := ruleset.NetworkRuleSetProperties.VirtualNetworkRules; vnetRules != nil {
		for _, vnetRule := range *vnetRules {
			block := make(map[string]interface{})

			if s := vnetRule.Subnet; s != nil {
				if v := s.ID; v != nil {
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
	if ipRules := ruleset.NetworkRuleSetProperties.IPRules; ipRules != nil {
		for _, ipRule := range *ipRules {
			block := make(map[string]interface{})

			block["action"] = string(ipRule.Action)

			if v := ipRule.IPMask; v != nil {
				block["ip_mask"] = *v
			}

			ipBlocks = append(ipBlocks, block)
		}
	}

	return []interface{}{map[string]interface{}{
		"default_action":       string(ruleset.DefaultAction),
		"virtual_network_rule": vnetBlocks,
		"ip_rule":              ipBlocks,
	}}
}

func expandEventHubIdentity(input []interface{}) *eventhub.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &eventhub.Identity{
		Type: eventhub.IdentityType(v["type"].(string)),
	}
}

func flattenEventHubIdentity(input *eventhub.Identity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	principalID := ""
	if input.PrincipalID != nil {
		principalID = *input.PrincipalID
	}

	tenantID := ""
	if input.TenantID != nil {
		tenantID = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
