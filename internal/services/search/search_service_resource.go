package search

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSearchService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSearchServiceCreate,
		Read:   resourceSearchServiceRead,
		Update: resourceSearchServiceUpdate,
		Delete: resourceSearchServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := services.ParseSearchServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			// NOTE: in the 2022-09-01 version of the API 'location'
			// is now just a string instead of a *string
			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.SkuNameFree),
					string(services.SkuNameBasic),
					string(services.SkuNameStandard),
					string(services.SkuNameStandardTwo),
					string(services.SkuNameStandardThree),
					string(services.SkuNameStorageOptimizedLOne),
					string(services.SkuNameStorageOptimizedLTwo),
				}, false),
			},

			"replica_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
			},

			"partition_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtMost(12),
			},

			"hosting_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  services.HostingModeDefault,
				ValidateFunc: validation.StringInSlice([]string{
					string(services.HostingModeDefault),
					string(services.HostingModeHighDensity),
				}, false),
			},

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"query_keys": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allowed_ips": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.Any(
						validate.IPv4Address,
						validate.CIDR,
					),
				},
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSearchServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := services.NewSearchServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, services.GetOperationOptions{})
	if err != nil && !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_search_service", id.ID())
	}

	properties, err := resourceSearchServiceCreateOrUpdateProperties(d)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, *properties, services.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSearchServiceRead(d, meta)
}

func resourceSearchServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Get(ctx, *id, services.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	properties, err := resourceSearchServiceCreateOrUpdateProperties(d)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	err = client.CreateOrUpdateThenPoll(ctx, *id, *properties, services.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceSearchServiceRead(d, meta)
}

func resourceSearchServiceCreateOrUpdateProperties(d *pluginsdk.ResourceData) (*services.SearchService, error) {
	// this might be broken because of how terraform treats ResourceData in the Create vs. the Update function... I think in
	// update it is pulling it from state instead of the config, like it does in Create...

	location := azure.NormalizeLocation(d.Get("location").(string))

	publicNetworkAccess := services.PublicNetworkAccessEnabled
	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		publicNetworkAccess = services.PublicNetworkAccessDisabled
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	skuName := services.SkuName(d.Get("sku").(string))
	hostingMode := services.HostingMode(d.Get("hosting_mode").(string))

	if skuName != services.SkuNameStandardThree && hostingMode == services.HostingModeHighDensity {
		return nil, fmt.Errorf("'hosting_mode' can only be set to 'highDensity' if the 'sku' is 'standard3', got %q", skuName)
	}

	searchService := services.SearchService{
		Location: location,
		Sku: &services.Sku{
			Name: &skuName,
		},
		Properties: &services.SearchServiceProperties{
			PublicNetworkAccess: &publicNetworkAccess,
			NetworkRuleSet: &services.NetworkRuleSet{
				IPRules: expandSearchServiceIPRules(d.Get("allowed_ips").([]interface{})),
			},
			HostingMode: &hostingMode,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("replica_count"); ok {
		replicaCount := int64(v.(int))
		searchService.Properties.ReplicaCount = utils.Int64(replicaCount)
	}

	// NOTE: 1 is now the "default" value for partitionCount...
	var partitionCount int64
	if v, ok := d.GetOk("partition_count"); ok {
		partitionCount = int64(v.(int))
		searchService.Properties.PartitionCount = utils.Int64(partitionCount)
	}

	// NOTE: 'partition_count' values greater than 1 are only valid for standard SKUs...
	if !strings.HasPrefix(strings.ToLower(string(skuName)), "standard") && partitionCount > 1 {
		return nil, fmt.Errorf("'partition_count' values greater than 1 are only valid for 'standard' SKUs, got (sku: %q, partition_count: %d)", skuName, partitionCount)
	}

	// NOTE: If SKU is 'standard3' and the 'hosting_mode' is set to 'highDensity' the maximum number of partitions allowed is 3
	// where if 'hosting_mode' is set to 'default' the maximum number of partitions is 12...
	if skuName == services.SkuNameStandardThree && partitionCount > 3 && hostingMode == services.HostingModeHighDensity {
		return nil, fmt.Errorf("'standard3' SKUs in 'highDensity' mode can have a maximum of 3 partitions, got %d", partitionCount)
	}

	return &searchService, nil
}

func resourceSearchServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, services.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SearchServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		skuName := ""
		if sku := model.Sku; sku != nil && sku.Name != nil {
			skuName = string(*sku.Name)
		}
		d.Set("sku", skuName)

		if props := model.Properties; props != nil {
			partitionCount := 0
			replicaCount := 0
			publicNetworkAccess := false
			hostingMode := services.HostingModeDefault

			if count := props.PartitionCount; count != nil {
				partitionCount = int(*count)
			}

			if count := props.ReplicaCount; count != nil {
				replicaCount = int(*count)
			}

			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess != services.PublicNetworkAccessDisabled
			}

			if props.HostingMode != nil {
				hostingMode = *props.HostingMode
			}

			d.Set("partition_count", partitionCount)
			d.Set("replica_count", replicaCount)
			d.Set("public_network_access_enabled", publicNetworkAccess)
			d.Set("hosting_mode", hostingMode)
			d.Set("allowed_ips", flattenSearchServiceIPRules(props.NetworkRuleSet))
		}

		if err = d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	adminKeysClient := meta.(*clients.Client).Search.AdminKeysClient
	adminKeysId, err := adminkeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	adminKeysResp, err := adminKeysClient.Get(ctx, *adminKeysId, adminkeys.GetOperationOptions{})
	if err == nil {
		if model := adminKeysResp.Model; model != nil {
			d.Set("primary_key", model.PrimaryKey)
			d.Set("secondary_key", model.SecondaryKey)
		}
	}

	queryKeysClient := meta.(*clients.Client).Search.QueryKeysClient
	queryKeysId, err := querykeys.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}
	queryKeysResp, err := queryKeysClient.ListBySearchService(ctx, *queryKeysId, querykeys.ListBySearchServiceOperationOptions{})
	if err == nil {
		if model := queryKeysResp.Model; model != nil {
			d.Set("query_keys", flattenSearchQueryKeys(*model))
		}
	}

	return nil
}

func resourceSearchServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Search.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id, services.DeleteOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenSearchQueryKeys(input []querykeys.QueryKey) []interface{} {
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

func expandSearchServiceIPRules(input []interface{}) *[]services.IPRule {
	output := make([]services.IPRule, 0)
	if input == nil {
		return &output
	}

	for _, rule := range input {
		if rule != nil {
			output = append(output, services.IPRule{
				Value: utils.String(rule.(string)),
			})
		}
	}

	return &output
}

func flattenSearchServiceIPRules(input *services.NetworkRuleSet) []interface{} {
	if input == nil || *input.IPRules == nil || len(*input.IPRules) == 0 {
		return nil
	}
	result := make([]interface{}, 0)
	for _, rule := range *input.IPRules {
		result = append(result, rule.Value)
	}
	return result
}
