package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-01-21/kusto"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAzureRMKustoClusterSkuName(),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKustoClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Kusto Cluster creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
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

	clusterProperties := kusto.ClusterProperties{}

	t := d.Get("tags").(map[string]interface{})

	kustoCluster := kusto.Cluster{
		Name:              &name,
		Location:          &location,
		Sku:               sku,
		ClusterProperties: &clusterProperties,
		Tags:              tags.Expand(t),
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
	client := meta.(*ArmClient).kusto.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["Clusters"]

	clusterResponse, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(clusterResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := clusterResponse.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("sku", flattenKustoClusterSku(clusterResponse.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if clusterProperties := clusterResponse.ClusterProperties; clusterProperties != nil {
		d.Set("uri", clusterProperties.URI)
		d.Set("data_ingestion_uri", clusterProperties.DataIngestionURI)
	}

	return tags.FlattenAndSet(d, clusterResponse.Tags)
}

func resourceArmKustoClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["Clusters"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Cluster %q (Resource Group %q): %+v", name, resGroup, err)
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

func validateAzureRMKustoClusterSkuName() schema.SchemaValidateFunc {
	// using hard coded values because they're not like this in the sdk as constants
	// found them here: https://docs.microsoft.com/en-us/rest/api/azurerekusto/clusters/createorupdate#azureskuname
	possibleSkuNames := []string{
		"Dev(No SLA)_Standard_D11_v2",
		"Standard_D11_v2",
		"Standard_D12_v2",
		"Standard_D13_v2",
		"Standard_D14_v2",
		"Standard_DS13_v2+1TB_PS",
		"Standard_DS13_v2+2TB_PS",
		"Standard_DS14_v2+3TB_PS",
		"Standard_DS14_v2+4TB_PS",
		"Standard_L16s",
		"Standard_L4s",
		"Standard_L8s",
	}

	return validation.StringInSlice(possibleSkuNames, false)
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
		Tier:     &tier,
		Capacity: utils.Int32(int32(capacity)),
	}

	return azureSku, nil
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
