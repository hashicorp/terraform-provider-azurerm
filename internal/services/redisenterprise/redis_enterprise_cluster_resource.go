package redisenterprise

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redisenterprise/mgmt/2021-03-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRedisEnterpriseCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisEnterpriseClusterCreate,
		Read:   resourceRedisEnterpriseClusterRead,
		Update: resourceRedisEnterpriseClusterUpdate,
		Delete: resourceRedisEnterpriseClusterDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RedisEnterpriseClusterID(id)
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
				ValidateFunc: validate.RedisEnterpriseName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RedisEnterpriseClusterLocation,
				StateFunc:    azure.NormalizeLocation,
			},

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RedisEnterpriseClusterSkuName,
			},

			"zones": azure.SchemaMultipleZones(),

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.OneFullStopTwo),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.OneFullStopZero),
					string(redisenterprise.OneFullStopOne),
					string(redisenterprise.OneFullStopTwo),
				}, false),
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			// RP currently does not return this value, but will in the near future
			// https://github.com/Azure/azure-sdk-for-go/issues/14420
			"version": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Deprecated: "This field currently is not yet being returned from the service API, please see https://github.com/Azure/azure-sdk-for-go/issues/14420 for more information",
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceRedisEnterpriseClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.Client
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Redis Enterprise Cluster creation.")

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceId := parse.NewRedisEnterpriseClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.RedisEnterpriseName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", resourceId.RedisEnterpriseName, resourceId.ResourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_redis_enterprise_cluster", resourceId.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	sku := expandRedisEnterpriseClusterSku(d.Get("sku_name").(string))

	// If the sku type is flash check to make sure that the sku is supported in that region
	if strings.Contains(string(sku.Name), "Flash") {
		if err := validate.RedisEnterpriseClusterLocationFlashSkuSupport(location); err != nil {
			return fmt.Errorf("%s: %s", resourceId, err)
		}
	}

	parameters := redisenterprise.Cluster{
		Name:     utils.String(d.Get("name").(string)),
		Location: utils.String(location),
		Sku:      sku,
		ClusterProperties: &redisenterprise.ClusterProperties{
			MinimumTLSVersion: redisenterprise.TLSVersion(d.Get("minimum_tls_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("zones"); ok {
		// Zones are currently not supported in these regions
		if err := validate.RedisEnterpriseClusterLocationZoneSupport(location); err != nil {
			return fmt.Errorf("%s: %s", resourceId, err)
		}
		parameters.Zones = azure.ExpandZones(v.([]interface{}))
	}

	future, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.RedisEnterpriseName, parameters)
	if err != nil {
		return fmt.Errorf("waiting for creation of Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", resourceId.RedisEnterpriseName, resourceId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", resourceId.RedisEnterpriseName, resourceId.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Redis Enterprise Cluster (Name %q / Resource Group %q) to become available", resourceId.RedisEnterpriseName, resourceId.ResourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Enabling", "Deleting", "Disabling"},
		Target:     []string{"Running"},
		Refresh:    redisEnterpriseClusterStateRefreshFunc(ctx, client, resourceId),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Redis Enterprise Cluster (Name %q / Resource Group %q) to become available: %+v", resourceId.RedisEnterpriseName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceRedisEnterpriseClusterRead(d, meta)
}

func resourceRedisEnterpriseClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RedisEnterpriseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RedisEnterpriseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Redis Enterprise Cluster (Name %q / Resource Group %q) was not found - removing from state!", id.RedisEnterpriseName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
	}

	d.Set("name", id.RedisEnterpriseName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("sku_name", flattenRedisEnterpriseClusterSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku_name`: %+v", err)
	}

	if zones := resp.Zones; zones != nil {
		d.Set("zones", zones)
	}

	if props := resp.ClusterProperties; props != nil {
		d.Set("hostname", props.HostName)
		d.Set("version", props.RedisVersion)
		d.Set("minimum_tls_version", string(props.MinimumTLSVersion))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceRedisEnterpriseClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.Client
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache update.")

	id, err := parse.RedisEnterpriseClusterID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := redisenterprise.ClusterUpdate{
		Tags: expandedTags,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.RedisEnterpriseName, parameters); err != nil {
		return fmt.Errorf("updating Redis Cache %q (Resource Group %q): %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Redis Cache %q (Resource Group %q) to become available", id.RedisEnterpriseName, id.ResourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Enabling", "Deleting", "Disabling"},
		Target:     []string{"Running"},
		Refresh:    redisEnterpriseClusterStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Redis Cache %q (Resource Group %q) to become available: %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
	}

	return resourceRedisEnterpriseClusterRead(d, meta)
}

func resourceRedisEnterpriseClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RedisEnterpriseClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RedisEnterpriseName)
	if err != nil {
		return fmt.Errorf("deleting Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
	}

	return nil
}

func expandRedisEnterpriseClusterSku(v string) *redisenterprise.Sku {
	redisSku, _ := parse.RedisEnterpriseCacheSkuName(v)
	capacity, _ := strconv.ParseInt(redisSku.Capacity, 10, 32)

	sku := &redisenterprise.Sku{
		Name:     redisenterprise.SkuName(redisSku.Name),
		Capacity: utils.Int32(int32(capacity)),
	}

	return sku
}

func flattenRedisEnterpriseClusterSku(input *redisenterprise.Sku) *string {
	if input == nil {
		return nil
	}

	var name redisenterprise.SkuName
	var capacity int32

	if input.Name != "" {
		name = input.Name
	}

	if input.Capacity != nil {
		capacity = *input.Capacity
	}

	skuName := fmt.Sprintf("%s-%d", name, capacity)

	return &skuName
}

func redisEnterpriseClusterStateRefreshFunc(ctx context.Context, client *redisenterprise.Client, id parse.RedisEnterpriseClusterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.RedisEnterpriseName)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving status of Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", id.RedisEnterpriseName, id.ResourceGroup, err)
		}

		return res, string(res.ClusterProperties.ResourceState), nil
	}
}
