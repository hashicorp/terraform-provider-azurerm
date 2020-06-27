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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var eventHubNamespaceDedicatedResourceName = "azurerm_eventhub_namespace_dedicated"

func resourceArmEventHubNamespaceDedicated() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmeventHubNamespaceDedicatedCreateUpdate,
		Read:   resourceArmeventHubNamespaceDedicatedRead,
		Update: resourceArmeventHubNamespaceDedicatedCreateUpdate,
		Delete: resourceArmeventHubNamespaceDedicatedDelete,
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
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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

func resourceArmeventHubNamespaceDedicatedCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace creation.")

	name := d.Get("name").(string)
	clusterID := d.Get("cluster_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace %q (Resource Group %q): %s", name, resourceGroup, err)
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

	parameters := eventhub.EHNamespace{
		Location: &location,
		Sku: &eventhub.Sku{
			Name:     eventhub.SkuName(sku),
			Tier:     eventhub.SkuTier(sku),
			Capacity: &capacity,
		},
		EHNamespaceProperties: &eventhub.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			ClusterArmID:         &clusterID,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.EHNamespaceProperties.MaximumThroughputUnits = utils.Int32(int32(v.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating eventhub namespace: %+v", err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	ruleSets, hasRuleSets := d.GetOk("network_rulesets")
	if hasRuleSets {
		rulesets := eventhub.NetworkRuleSet{
			NetworkRuleSetProperties: expandEventHubNamespaceNetworkRuleset(ruleSets.([]interface{})),
		}

		// cannot use network rulesets with the basic SKU
		if parameters.Sku.Name != eventhub.Basic {
			if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, resourceGroup, name, rulesets); err != nil {
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

	return resourceArmeventHubNamespaceDedicatedRead(d, meta)
}

func resourceArmeventHubNamespaceDedicatedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on EventHub Namespace %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
		d.Set("capacity", sku.Capacity)
	}

	if props := resp.EHNamespaceProperties; props != nil {
		d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
		d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
		d.Set("cluster_id", props.ClusterArmID)
	}

	ruleset, err := client.GetNetworkRuleSet(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Dedicated EventHub Namespace %q Network Ruleset: %+v", name, err)
	}

	if err := d.Set("network_rulesets", flattenEventHubNamespaceNetworkRuleset(ruleset)); err != nil {
		return fmt.Errorf("Error setting `network_ruleset` for Dedicated Eventhub Namespace %s: %v", name, err)
	}

	keys, err := client.ListKeys(ctx, resGroup, name, eventHubNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for Dedicated EventHub Namespace %q: %+v", name, err)
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

func resourceArmeventHubNamespaceDedicatedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request of EventHub Namespace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return waitForeventHubNamespaceDedicatedToBeDeleted(ctx, client, resGroup, name, d)
}

func waitForeventHubNamespaceDedicatedToBeDeleted(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string, d *schema.ResourceData) error {
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for Dedicated EventHub Namespace (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: eventHubNamespaceDedicatedStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Dedicated EventHub NameSpace (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func eventHubNamespaceDedicatedStateStatusCodeRefreshFunc(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving Dedicated EventHub Namespace %q (RG %q, ClusterID %q) returned Status %d", name, resourceGroup, *res.ClusterArmID, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the Dedicated EventHub Namespace %q (RG: %q, ClusterID: %q): %+v", name, resourceGroup, *res.ClusterArmID, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
