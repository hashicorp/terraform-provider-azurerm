// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterCreate,
		Read:   resourceKustoClusterRead,
		Update: resourceKustoClusterUpdate,
		Delete: resourceKustoClusterDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoAttachedClusterV0ToV1{},
			1: migration.KustoAttachedClusterV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseKustoClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForAzureSkuName(),
								false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
						},
					},
				},
			},

			"allowed_fqdns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"allowed_ip_ranges": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"trusted_external_tenants": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty, validation.StringInSlice([]string{"*"}, false)),
				},
			},

			"optimized_auto_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"minimum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"maximum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},

			"uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_ingestion_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"language_extension": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionName(), false),
						},
						"image": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionImageName(), false),
						},
					},
				},
			},

			"public_ip_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(clusters.PublicIPTypeIPvFour),
				ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPublicIPType(), false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"outbound_network_access_restricted": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"double_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"auto_stop_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"disk_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"streaming_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"purge_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceKustoClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKustoClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster", id.ID())
		}
	}

	locks.ByName(id.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.KustoClusterName, "azurerm_kusto_cluster")

	sku, err := expandKustoClusterSku(d.Get("sku").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `sku`: %+v", err)
	}

	optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]interface{}))

	if optimizedAutoScale != nil && optimizedAutoScale.IsEnabled {
		if sku.Capacity == nil {
			return fmt.Errorf("sku.capacity could not be empty")
		}
		// Ensure that requested Capcity is always between min and max to support updating to not overlapping autoscale ranges
		if *sku.Capacity < optimizedAutoScale.Minimum {
			sku.Capacity = pointer.To(optimizedAutoScale.Minimum)
		}
		if *sku.Capacity > optimizedAutoScale.Maximum {
			sku.Capacity = pointer.To(optimizedAutoScale.Maximum)
		}
		if optimizedAutoScale.Minimum > optimizedAutoScale.Maximum {
			return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
		}
	}

	publicNetworkAccess := clusters.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = clusters.PublicNetworkAccessDisabled
	}

	publicIPType := clusters.PublicIPType(d.Get("public_ip_type").(string))

	clusterProperties := clusters.ClusterProperties{
		OptimizedAutoscale:     optimizedAutoScale,
		EnableAutoStop:         pointer.To(d.Get("auto_stop_enabled").(bool)),
		EnableDiskEncryption:   pointer.To(d.Get("disk_encryption_enabled").(bool)),
		EnableDoubleEncryption: pointer.To(d.Get("double_encryption_enabled").(bool)),
		EnableStreamingIngest:  pointer.To(d.Get("streaming_ingestion_enabled").(bool)),
		EnablePurge:            pointer.To(d.Get("purge_enabled").(bool)),
		PublicNetworkAccess:    &publicNetworkAccess,
		PublicIPType:           &publicIPType,
		TrustedExternalTenants: expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]interface{})),
	}

	if v, ok := d.GetOk("allowed_fqdns"); ok {
		clusterProperties.AllowedFqdnList = expandKustoListString(v.([]interface{}))
	}

	if v, ok := d.GetOk("allowed_ip_ranges"); ok {
		clusterProperties.AllowedIPRangeList = expandKustoListString(v.([]interface{}))
	}

	restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
	if v, ok := d.GetOk("outbound_network_access_restricted"); ok {
		if v.(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
	}
	clusterProperties.RestrictOutboundNetworkAccess = &restrictOutboundNetworkAccess

	if v, ok := d.GetOk("language_extension"); ok {
		clusterProperties.LanguageExtensions = expandKustoClusterLanguageExtensionList(v.([]any))
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	kustoCluster := clusters.Cluster{
		Location:   location.Normalize(d.Get("location").(string)),
		Identity:   expandedIdentity,
		Sku:        *sku,
		Properties: &clusterProperties,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		kustoCluster.Zones = &zones
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, kustoCluster, clusters.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoClusterRead(d, meta)
}

func resourceKustoClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.KustoClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.KustoClusterName, "azurerm_kusto_cluster")

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	model := existing.Model
	props := model.Properties

	if d.HasChanges("sku", "optimized_auto_scale") {
		sku, err := expandKustoClusterSku(d.Get("sku").([]interface{}))
		if err != nil {
			return err
		}

		optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]interface{}))

		if optimizedAutoScale != nil && optimizedAutoScale.IsEnabled {
			if sku.Capacity == nil {
				return fmt.Errorf("sku.capacity cannot be empty")
			}
			// Ensure that requested Capcity is always between min and max to support updating to not overlapping autoscale ranges
			if *sku.Capacity < optimizedAutoScale.Minimum {
				sku.Capacity = pointer.To(optimizedAutoScale.Minimum)
			}
			if *sku.Capacity > optimizedAutoScale.Maximum {
				sku.Capacity = pointer.To(optimizedAutoScale.Maximum)
			}
			if optimizedAutoScale.Minimum > optimizedAutoScale.Maximum {
				return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
			}
		}
		model.Sku = *sku
		// optimized_auto_scale can't be updated when the sku is also being updated, so we'll update sku first and then optimized_auto_scale
		if d.HasChange("optimized_auto_scale") && d.HasChange("sku") {
			model.Properties.OptimizedAutoscale = nil
			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model, clusters.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
		props.OptimizedAutoscale = optimizedAutoScale
	}

	if d.HasChange("location") {
		model.Location = location.Normalize(d.Get("location").(string))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		model.Identity = expandedIdentity
	}

	if d.HasChange("allowed_fqdns") {
		props.AllowedFqdnList = expandKustoListString(d.Get("allowed_fqdns").([]interface{}))
	}

	if d.HasChange("allowed_ip_ranges") {
		props.AllowedIPRangeList = expandKustoListString(d.Get("allowed_ip_ranges").([]interface{}))
	}

	if d.HasChange("auto_stop_enabled") {
		props.EnableAutoStop = pointer.To(d.Get("auto_stop_enabled").(bool))
	}

	if d.HasChange("disk_encryption_enabled") {
		props.EnableDiskEncryption = pointer.To(d.Get("disk_encryption_enabled").(bool))
	}

	if d.HasChange("double_encryption_enabled") {
		props.EnableDoubleEncryption = pointer.To(d.Get("double_encryption_enabled").(bool))
	}

	if d.HasChange("language_extension") {
		props.LanguageExtensions = expandKustoClusterLanguageExtensionList(d.Get("language_extension").([]any))
	}

	if d.HasChange("outbound_network_access_restricted") {
		restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
		if d.Get("outbound_network_access_restricted").(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
		props.RestrictOutboundNetworkAccess = pointer.To(restrictOutboundNetworkAccess)
	}

	if d.HasChange("purge_enabled") {
		props.EnablePurge = pointer.To(d.Get("purge_enabled").(bool))
	}

	if d.HasChange("public_ip_type") {
		publicIPType := clusters.PublicIPType(d.Get("public_ip_type").(string))
		props.PublicIPType = pointer.To(publicIPType)
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := clusters.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = clusters.PublicNetworkAccessDisabled
		}
		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("streaming_ingestion_enabled") {
		props.EnableStreamingIngest = pointer.To(d.Get("streaming_ingestion_enabled").(bool))
	}

	if d.HasChange("trusted_external_tenants") {
		props.TrustedExternalTenants = expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]interface{}))
	}

	if d.HasChange("zones") {
		zoneList := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		model.Zones = pointer.To(zoneList)
		if len(zoneList) > 0 {
			model.Zones = &zoneList
		} else {
			model.Zones = nil
		}
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	model.Properties = props

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *model, clusters.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterRead(d, meta)
}

func resourceKustoClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
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

	d.Set("name", id.KustoClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := d.Set("sku", flattenKustoClusterSku(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if props.PublicNetworkAccess != nil {
				d.Set("public_network_access_enabled", *props.PublicNetworkAccess == clusters.PublicNetworkAccessEnabled)
			}

			if props.RestrictOutboundNetworkAccess != nil {
				d.Set("outbound_network_access_restricted", *props.RestrictOutboundNetworkAccess == clusters.ClusterNetworkAccessFlagEnabled)
			}

			if err := d.Set("optimized_auto_scale", flattenOptimizedAutoScale(props.OptimizedAutoscale)); err != nil {
				return fmt.Errorf("setting `optimized_auto_scale`: %+v", err)
			}
			d.Set("allowed_fqdns", props.AllowedFqdnList)
			d.Set("allowed_ip_ranges", props.AllowedIPRangeList)
			d.Set("double_encryption_enabled", props.EnableDoubleEncryption)
			d.Set("trusted_external_tenants", flattenTrustedExternalTenants(props.TrustedExternalTenants))
			d.Set("auto_stop_enabled", props.EnableAutoStop)
			d.Set("disk_encryption_enabled", props.EnableDiskEncryption)
			d.Set("streaming_ingestion_enabled", props.EnableStreamingIngest)
			d.Set("purge_enabled", props.EnablePurge)
			d.Set("uri", props.Uri)
			d.Set("data_ingestion_uri", props.DataIngestionUri)
			d.Set("public_ip_type", string(pointer.From(props.PublicIPType)))

			d.Set("language_extension", flattenKustoClusterLanguageExtensionList(props.LanguageExtensions))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceKustoClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoClusterID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandOptimizedAutoScale(input []interface{}) *clusters.OptimizedAutoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	optimizedAutoScale := &clusters.OptimizedAutoscale{
		Version:   1,
		IsEnabled: true,
		Minimum:   int64(config["minimum_instances"].(int)),
		Maximum:   int64(config["maximum_instances"].(int)),
	}

	return optimizedAutoScale
}

func flattenOptimizedAutoScale(optimizedAutoScale *clusters.OptimizedAutoscale) []interface{} {
	if optimizedAutoScale == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"maximum_instances": int(optimizedAutoScale.Maximum),
			"minimum_instances": int(optimizedAutoScale.Minimum),
		},
	}
}

func expandKustoListString(input []interface{}) *[]string {
	result := make([]string, 0)

	for _, v := range input {
		result = append(result, v.(string))
	}

	return &result
}

func expandKustoClusterSku(input []interface{}) (*clusters.AzureSku, error) {
	sku := input[0].(map[string]interface{})
	name := sku["name"].(string)

	skuNamePrefixToTier := map[string]string{
		"Dev(No SLA)": "Basic",
		"Standard":    "Standard",
	}

	skuNamePrefix := strings.Split(sku["name"].(string), "_")[0]
	tier, ok := skuNamePrefixToTier[skuNamePrefix]
	if !ok {
		return nil, fmt.Errorf("sku name begins with invalid tier, possible are Dev(No SLA) and Standard but is: %q", skuNamePrefix)
	}
	capacity := sku["capacity"].(int)

	azureSku := clusters.AzureSku{
		Name:     clusters.AzureSkuName(name),
		Tier:     clusters.AzureSkuTier(tier),
		Capacity: pointer.To(int64(capacity)),
	}

	return &azureSku, nil
}

func expandKustoClusterLanguageExtensionList(input []interface{}) *clusters.LanguageExtensionsList {
	extensions := make([]clusters.LanguageExtension, 0)

	for _, ext := range input {
		extMap := ext.(map[string]interface{})
		extensions = append(extensions, clusters.LanguageExtension{
			LanguageExtensionName:      pointer.To(clusters.LanguageExtensionName(extMap["name"].(string))),
			LanguageExtensionImageName: pointer.To(clusters.LanguageExtensionImageName(extMap["image"].(string))),
		})
	}

	return &clusters.LanguageExtensionsList{
		Value: &extensions,
	}
}

func flattenKustoClusterSku(sku *clusters.AzureSku) []interface{} {
	if sku == nil {
		return []interface{}{}
	}

	s := map[string]interface{}{
		"name": string(sku.Name),
	}

	if sku.Capacity != nil {
		s["capacity"] = int(*sku.Capacity)
	}

	return []interface{}{s}
}

func flattenKustoClusterLanguageExtensionList(extensions *clusters.LanguageExtensionsList) []interface{} {
	if extensions == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	if extensions.Value != nil {
		for _, v := range *extensions.Value {
			output = append(
				output,
				map[string]interface{}{
					"name":  string(*v.LanguageExtensionName),
					"image": string(*v.LanguageExtensionImageName),
				},
			)
		}
	}

	return output
}
