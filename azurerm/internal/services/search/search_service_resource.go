package search

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2020-03-13/search"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/search/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSearchService() *schema.Resource {
	return &schema.Resource{
		Create: resourceSearchServiceCreateUpdate,
		Read:   resourceSearchServiceRead,
		Update: resourceSearchServiceCreateUpdate,
		Delete: resourceSearchServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SearchServiceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(search.Free),
					string(search.Basic),
					string(search.Standard),
					string(search.Standard2),
					string(search.Standard3),
					string(search.StorageOptimizedL1),
					string(search.StorageOptimizedL2),
				}, false),
			},

			"replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"partition_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtMost(12),
			},

			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"query_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allowed_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validate.IPv4Address,
						validate.CIDR,
					),
				},
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
								string(search.SystemAssigned),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceSearchServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	skuName := d.Get("sku").(string)

	publicNetworkAccess := search.Enabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = search.Disabled
	}

	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, nil)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Search Service %q (ResourceGroup %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_search_service", *existing.ID)
		}
	}

	properties := search.Service{
		Location: utils.String(location),
		Sku: &search.Sku{
			Name: search.SkuName(skuName),
		},
		ServiceProperties: &search.ServiceProperties{
			PublicNetworkAccess: publicNetworkAccess,
			NetworkRuleSet: &search.NetworkRuleSet{
				IPRules: expandSearchServiceIPRules(d.Get("allowed_ips").([]interface{})),
			},
		},
		Identity: expandSearchServiceIdentity(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(t),
	}

	if v, ok := d.GetOk("replica_count"); ok {
		replicaCount := int32(v.(int))
		properties.ServiceProperties.ReplicaCount = utils.Int32(replicaCount)
	}

	if v, ok := d.GetOk("partition_count"); ok {
		partitionCount := int32(v.(int))
		properties.ServiceProperties.PartitionCount = utils.Int32(partitionCount)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties, nil)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Search Service %q (ResourceGroup %q): %s", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for the completion of the creating/updating of Search Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Search Service %q (ResourceGroup %q): %s", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceSearchServiceRead(d, meta)
}

func resourceSearchServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Search Service %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Search Service: %+v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.ServiceProperties; props != nil {
		if count := props.PartitionCount; count != nil {
			d.Set("partition_count", int(*count))
		}

		if count := props.ReplicaCount; count != nil {
			d.Set("replica_count", int(*count))
		}

		d.Set("public_network_access_enabled", props.PublicNetworkAccess != "Disabled")

		d.Set("allowed_ips", flattenSearchServiceIPRules(props.NetworkRuleSet))
	}

	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysResp, err := adminKeysClient.Get(ctx, id.ResourceGroup, id.Name, nil)
	if err == nil {
		d.Set("primary_key", adminKeysResp.PrimaryKey)
		d.Set("secondary_key", adminKeysResp.SecondaryKey)
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, id.ResourceGroup, id.Name, nil)
	if err == nil {
		d.Set("query_keys", flattenSearchQueryKeys(queryKeysResp.Values()))
	}

	if err := d.Set("identity", flattenSearchServiceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSearchServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name, nil)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Search Service %q (resource group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func flattenSearchQueryKeys(input []search.QueryKey) []interface{} {
	results := make([]interface{}, 0)

	for _, v := range input {
		result := make(map[string]interface{})

		if v.Name != nil {
			result["name"] = *v.Name
		}
		result["key"] = *v.Key

		results = append(results, result)
	}

	return results
}

func expandSearchServiceIPRules(input []interface{}) *[]search.IPRule {
	output := make([]search.IPRule, 0)
	if input == nil {
		return &output
	}

	for _, rule := range input {
		if rule != nil {
			output = append(output, search.IPRule{
				Value: utils.String(rule.(string)),
			})
		}
	}

	return &output
}

func flattenSearchServiceIPRules(input *search.NetworkRuleSet) []interface{} {
	if input == nil || *input.IPRules == nil || len(*input.IPRules) == 0 {
		return nil
	}
	result := make([]interface{}, 0)
	for _, rule := range *input.IPRules {
		result = append(result, rule.Value)
	}
	return result
}

func expandSearchServiceIdentity(input []interface{}) *search.Identity {
	if len(input) == 0 || input[0] == nil {
		return &search.Identity{
			Type: search.None,
		}
	}
	identity := input[0].(map[string]interface{})
	return &search.Identity{
		Type: search.IdentityType(identity["type"].(string)),
	}
}

func flattenSearchServiceIdentity(identity *search.Identity) []interface{} {
	if identity == nil || identity.Type == search.None {
		return make([]interface{}, 0)
	}

	principalId := ""
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	tenantId := ""
	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"principal_id": principalId,
			"tenant_id":    tenantId,
			"type":         string(identity.Type),
		},
	}
}
