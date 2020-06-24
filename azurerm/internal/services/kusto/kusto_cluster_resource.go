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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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

			"identity": azure.SchemaKustoIdentity(),

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
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
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

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	sku, err := expandKustoClusterSku(d)
	if err != nil {
		return err
	}

	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	clusterProperties := kusto.ClusterProperties{
		EnableDiskEncryption:  utils.Bool(d.Get("enable_disk_encryption").(bool)),
		EnableStreamingIngest: utils.Bool(d.Get("enable_streaming_ingest").(bool)),
		EnablePurge:           utils.Bool(d.Get("enable_purge").(bool)),
	}

	if v, ok := d.GetOk("virtual_network_configuration"); ok {
		vnet := expandKustoClusterVNET(v.([]interface{}))
		clusterProperties.VirtualNetworkConfiguration = vnet
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
		kustoIdentity := azure.ExpandKustoIdentity(kustoIdentityRaw)
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

	return resourceArmKustoClusterRead(d, meta)
}

func resourceArmKustoClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoClusterID(d.Id())
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

	if err := d.Set("identity", azure.FlattenKustoIdentity(clusterResponse.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %s", err)
	}

	if err := d.Set("sku", flattenKustoClusterSku(clusterResponse.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if err := d.Set("zones", azure.FlattenZones(clusterResponse.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	if clusterProperties := clusterResponse.ClusterProperties; clusterProperties != nil {
		d.Set("enable_disk_encryption", clusterProperties.EnableDiskEncryption)
		d.Set("enable_streaming_ingest", clusterProperties.EnableStreamingIngest)
		d.Set("enable_purge", clusterProperties.EnablePurge)
		d.Set("virtual_network_configuration", flatteKustoClusterVNET(clusterProperties.VirtualNetworkConfiguration))
		d.Set("uri", clusterProperties.URI)
		d.Set("data_ingestion_uri", clusterProperties.DataIngestionURI)
	}

	return tags.FlattenAndSet(d, clusterResponse.Tags)
}

func resourceArmKustoClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoClusterID(d.Id())
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

func expandKustoClusterSku(d *schema.ResourceData) (*kusto.AzureSku, error) {
	skuList := d.Get("sku").([]interface{})

	sku := skuList[0].(map[string]interface{})
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
