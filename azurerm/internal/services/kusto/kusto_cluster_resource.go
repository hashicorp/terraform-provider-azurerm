package kusto

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoClusterCreateUpdate,
		Read:   resourceArmKustoClusterRead,
		Update: resourceArmKustoClusterCreateUpdate,
		Delete: resourceArmKustoClusterDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"identity": schemaIdentity(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(kusto.DevNoSLAStandardD11V2),
								string(kusto.DevNoSLAStandardE2aV4),
								string(kusto.StandardD11V2),
								string(kusto.StandardD12V2),
								string(kusto.StandardD13V2),
								string(kusto.StandardD14V2),
								string(kusto.StandardDS13V21TBPS),
								string(kusto.StandardDS13V22TBPS),
								string(kusto.StandardDS14V23TBPS),
								string(kusto.StandardDS14V24TBPS),
								string(kusto.StandardE16asV43TBPS),
								string(kusto.StandardE16asV44TBPS),
								string(kusto.StandardE16aV4),
								string(kusto.StandardE2aV4),
								string(kusto.StandardE4aV4),
								string(kusto.StandardE8asV41TBPS),
								string(kusto.StandardE8asV42TBPS),
								string(kusto.StandardE8aV4),
								string(kusto.StandardL16s),
								string(kusto.StandardL4s),
								string(kusto.StandardL8s),
							}, false),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
						},
					},
				},
			},

			"trusted_external_tenants": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
				},
			},

			"optimized_auto_scale": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum_instances": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"maximum_instances": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},

			"enable_disk_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_streaming_ingest": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_purge": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"virtual_network_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"engine_public_ip_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"data_management_public_ip_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"language_extensions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(kusto.PYTHON),
						string(kusto.R),
					}, false),
				},
			},

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"data_ingestion_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaZones(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKustoClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Cluster creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		server, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Cluster %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if server.ID != nil && *server.ID != "" {
			return tf.ImportAsExistsError("azurerm_kusto_cluster", *server.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	sku, err := expandKustoClusterSku(d.Get("sku").([]interface{}))
	if err != nil {
		return err
	}

	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]interface{}))

	if optimizedAutoScale != nil && *optimizedAutoScale.IsEnabled {
		// if Capacity has not been set use min instances
		if *sku.Capacity == 0 {
			sku.Capacity = utils.Int32(*optimizedAutoScale.Minimum)
		}

		// Capacity must be set for the initial creation when using OptimizedAutoScaling but cannot be updated
		if d.HasChange("sku.0.capacity") && !d.IsNewResource() {
			return fmt.Errorf("cannot change `sku.capacity` when `optimized_auto_scaling.enabled` is set to `true`")
		}

		if *optimizedAutoScale.Minimum > *optimizedAutoScale.Maximum {
			return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
		}
	}

	clusterProperties := kusto.ClusterProperties{
		OptimizedAutoscale:    optimizedAutoScale,
		EnableDiskEncryption:  utils.Bool(d.Get("enable_disk_encryption").(bool)),
		EnableStreamingIngest: utils.Bool(d.Get("enable_streaming_ingest").(bool)),
		EnablePurge:           utils.Bool(d.Get("enable_purge").(bool)),
	}

	if v, ok := d.GetOk("virtual_network_configuration"); ok {
		vnet := expandKustoClusterVNET(v.([]interface{}))
		clusterProperties.VirtualNetworkConfiguration = vnet
	}

	if v, ok := d.GetOk("trusted_external_tenants"); ok {
		trustedExternalTenants := expandTrustedExternalTenants(v.([]interface{}))
		clusterProperties.TrustedExternalTenants = trustedExternalTenants
	}

	t := d.Get("tags").(map[string]interface{})

	kustoCluster := kusto.Cluster{
		Name:              &name,
		Location:          &location,
		Sku:               sku,
		Zones:             zones,
		ClusterProperties: &clusterProperties,
		Tags:              tags.Expand(t),
	}

	if _, ok := d.GetOk("identity"); ok {
		kustoIdentityRaw := d.Get("identity").([]interface{})
		kustoIdentity := expandIdentity(kustoIdentityRaw)
		kustoCluster.Identity = kustoIdentity
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, kustoCluster)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, getDetailsErr := client.Get(ctx, resourceGroup, name)
	if getDetailsErr != nil {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Kusto Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	if v, ok := d.GetOk("language_extensions"); ok {
		languageExtensions := expandKustoClusterLanguageExtensions(v.([]interface{}))

		currentLanguageExtensions, err := client.ListLanguageExtensions(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Error reading current added language extensions from Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		languageExtensionsToAdd := diffLanguageExtensions(*languageExtensions.Value, *currentLanguageExtensions.Value)
		if len(languageExtensionsToAdd) > 0 {
			languageExtensionsListToAdd := kusto.LanguageExtensionsList{
				Value: &languageExtensionsToAdd,
			}

			addLanguageExtensionsFuture, err := client.AddLanguageExtensions(ctx, resourceGroup, name, languageExtensionsListToAdd)
			if err != nil {
				return fmt.Errorf("Error adding language extensions to Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}

			if err = addLanguageExtensionsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for completion of adding language extensions to Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		languageExtensionsToRemove := diffLanguageExtensions(*currentLanguageExtensions.Value, *languageExtensions.Value)
		if len(languageExtensionsToRemove) > 0 {
			languageExtensionsListToRemove := kusto.LanguageExtensionsList{
				Value: &languageExtensionsToRemove,
			}

			removeLanguageExtensionsFuture, err := client.RemoveLanguageExtensions(ctx, resourceGroup, name, languageExtensionsListToRemove)
			if err != nil {
				return fmt.Errorf("Error removing language extensions from Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
			if err = removeLanguageExtensionsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for completion of removing language extensions from Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	return resourceArmKustoClusterRead(d, meta)
}

func resourceArmKustoClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	clusterResponse, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(clusterResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := clusterResponse.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenIdentity(clusterResponse.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	if err := d.Set("sku", flattenKustoClusterSku(clusterResponse.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if err := d.Set("zones", azure.FlattenZones(clusterResponse.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}
	if err := d.Set("optimized_auto_scale", flattenOptimizedAutoScale(clusterResponse.OptimizedAutoscale)); err != nil {
		return fmt.Errorf("Error setting `optimized_auto_scale`: %+v", err)
	}

	if clusterProperties := clusterResponse.ClusterProperties; clusterProperties != nil {
		d.Set("trusted_external_tenants", flattenTrustedExternalTenants(clusterProperties.TrustedExternalTenants))
		d.Set("enable_disk_encryption", clusterProperties.EnableDiskEncryption)
		d.Set("enable_streaming_ingest", clusterProperties.EnableStreamingIngest)
		d.Set("enable_purge", clusterProperties.EnablePurge)
		d.Set("virtual_network_configuration", flatteKustoClusterVNET(clusterProperties.VirtualNetworkConfiguration))
		d.Set("language_extensions", flattenKustoClusterLanguageExtensions(clusterProperties.LanguageExtensions))
		d.Set("uri", clusterProperties.URI)
		d.Set("data_ingestion_uri", clusterProperties.DataIngestionURI)
	}

	return tags.FlattenAndSet(d, clusterResponse.Tags)
}

func resourceArmKustoClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func validateAzureRMKustoClusterName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if !regexp.MustCompile(`^[a-z][a-z0-9]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must begin with a letter and may only contain alphanumeric characters: %q", k, name))
	}

	if len(name) < 4 || len(name) > 22 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func expandOptimizedAutoScale(input []interface{}) *kusto.OptimizedAutoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	optimizedAutoScale := &kusto.OptimizedAutoscale{
		Version:   utils.Int32(1),
		IsEnabled: utils.Bool(true),
		Minimum:   utils.Int32(int32(config["minimum_instances"].(int))),
		Maximum:   utils.Int32(int32(config["maximum_instances"].(int))),
	}

	return optimizedAutoScale
}

func flattenOptimizedAutoScale(optimizedAutoScale *kusto.OptimizedAutoscale) []interface{} {
	if optimizedAutoScale == nil {
		return []interface{}{}
	}

	maxInstances := 0
	if optimizedAutoScale.Maximum != nil {
		maxInstances = int(*optimizedAutoScale.Maximum)
	}

	minInstances := 0
	if optimizedAutoScale.Minimum != nil {
		minInstances = int(*optimizedAutoScale.Minimum)
	}

	return []interface{}{
		map[string]interface{}{
			"maximum_instances": maxInstances,
			"minimum_instances": minInstances,
		},
	}
}

func expandKustoClusterSku(input []interface{}) (*kusto.AzureSku, error) {
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

	azureSku := &kusto.AzureSku{
		Name:     kusto.AzureSkuName(name),
		Tier:     kusto.AzureSkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
	}

	return azureSku, nil
}

func expandKustoClusterVNET(input []interface{}) *kusto.VirtualNetworkConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vnet := input[0].(map[string]interface{})
	subnetID := vnet["subnet_id"].(string)
	enginePublicIPID := vnet["engine_public_ip_id"].(string)
	dataManagementPublicIPID := vnet["data_management_public_ip_id"].(string)

	return &kusto.VirtualNetworkConfiguration{
		SubnetID:                 &subnetID,
		EnginePublicIPID:         &enginePublicIPID,
		DataManagementPublicIPID: &dataManagementPublicIPID,
	}
}

func expandKustoClusterLanguageExtensions(input []interface{}) *kusto.LanguageExtensionsList {
	if len(input) == 0 {
		return nil
	}

	extensions := make([]kusto.LanguageExtension, 0)
	for _, language := range input {
		v := kusto.LanguageExtension{
			LanguageExtensionName: kusto.LanguageExtensionName(language.(string)),
		}
		extensions = append(extensions, v)
	}

	return &kusto.LanguageExtensionsList{
		Value: &extensions,
	}
}

func flattenKustoClusterSku(sku *kusto.AzureSku) []interface{} {
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

func flatteKustoClusterVNET(vnet *kusto.VirtualNetworkConfiguration) []interface{} {
	if vnet == nil {
		return []interface{}{}
	}

	subnetID := ""
	if vnet.SubnetID != nil {
		subnetID = *vnet.SubnetID
	}

	enginePublicIPID := ""
	if vnet.EnginePublicIPID != nil {
		enginePublicIPID = *vnet.EnginePublicIPID
	}

	dataManagementPublicIPID := ""
	if vnet.DataManagementPublicIPID != nil {
		dataManagementPublicIPID = *vnet.DataManagementPublicIPID
	}

	output := map[string]interface{}{
		"subnet_id":                    subnetID,
		"engine_public_ip_id":          enginePublicIPID,
		"data_management_public_ip_id": dataManagementPublicIPID,
	}

	return []interface{}{output}
}

func flattenKustoClusterLanguageExtensions(extensions *kusto.LanguageExtensionsList) []interface{} {
	if extensions == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *extensions.Value {
		output = append(output, v.LanguageExtensionName)
	}

	return output
}

func diffLanguageExtensions(a, b []kusto.LanguageExtension) []kusto.LanguageExtension {
	target := make(map[string]bool)
	for _, x := range b {
		target[string(x.LanguageExtensionName)] = true
	}

	diff := make([]kusto.LanguageExtension, 0)
	for _, x := range a {
		if _, ok := target[string(x.LanguageExtensionName)]; !ok {
			diff = append(diff, x)
		}
	}

	return diff
}
