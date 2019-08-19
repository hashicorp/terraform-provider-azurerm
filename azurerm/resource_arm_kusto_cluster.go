package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-01-21/kusto"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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
							ValidateFunc: validateAzureRMKustoClusterSkuName,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"trusted_external_tenants": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tagsSchema(),
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

	sku := expandKustoClusterSku(d)

	clusterProperties := expandKustoClusterProperties(d)

	tags := d.Get("tags").(map[string]interface{})

	kustoCluster := kusto.Cluster{
		Name:              &name,
		Location:          &location,
		Sku:               sku,
		ClusterProperties: clusterProperties,
		Tags:              expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, kustoCluster)
	if err != nil {
		return fmt.Errorf("Error creating Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Analysis Services Server %q (Resource Group %q): %+v", name, resourceGroup, err)
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
	name := id.Path["clusters"]

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

	if sku := clusterResponse.Sku; sku != nil {
		d.Set("sku", flattenKustoClusterSku(clusterResponse.Sku))
	}

	if clusterProperties := clusterResponse.ClusterProperties; clusterProperties != nil {
		if clusterProperties.TrustedExternalTenants != nil {
			trustedTenantIds := make([]string, len(*clusterProperties.TrustedExternalTenants))
			for i, tenant := range *clusterProperties.TrustedExternalTenants {
				trustedTenantIds[i] = *tenant.Value
			}
		}
	}

	flattenAndSetTags(d, clusterResponse.Tags)

	return nil
}

func resourceArmKustoClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.ClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["clusters"]

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

func validateAzureRMKustoClusterSkuName(v interface{}, k string) (warnings []string, errors []error) {
	skuName := v.(string)
	skuParts := strings.Split(skuName, "_")

	if len(skuParts) != 2 {
		errors = append(errors, fmt.Errorf("%q must be in this format: TIER_NAMEs but is: %q", k, skuName))
		return warnings, errors
	}

	if skuParts[0] != "Standard" {
		errors = append(errors, fmt.Errorf("possible tiers in %q is \"Standard\" but is: %q", k, skuParts[0]))
	}

	skuIsValid := false
	for _, sku := range kusto.PossibleAzureSkuNameValues() {
		if string(sku)+"s" == skuParts[1] {
			skuIsValid = true
		}
	}

	if !skuIsValid {
		possibleSkuNames := make([]string, len(kusto.PossibleAzureSkuNameValues()))
		for i, sku := range kusto.PossibleAzureSkuNameValues() {
			possibleSkuNames[i] = string(sku)
		}
		errors = append(errors, fmt.Errorf("possible name in %q is one of \"%q\" but is: %q", k, strings.Join(possibleSkuNames, "\", \""), skuParts[1]))
	}

	return warnings, errors
}

func expandKustoClusterSku(d *schema.ResourceData) *kusto.AzureSku {
	skuList := d.Get("sku").([]interface{})

	sku := skuList[0].(map[string]interface{})
	name := sku["name"].(string)
	tier := strings.Split(sku["name"].(string), "_")[0]

	azureSku := &kusto.AzureSku{
		Name: kusto.AzureSkuName(name),
		Tier: &tier,
	}

	if capacity, ok := sku["capacity"]; ok {
		azureSku.Capacity = utils.Int32(int32(capacity.(int)))
	}

	return azureSku
}

func expandKustoClusterProperties(d *schema.ResourceData) *kusto.ClusterProperties {
	clusterProperties := kusto.ClusterProperties{}

	tenantIds := d.Get("trusted_external_tenants").(*schema.Set)
	if len(tenantIds.List()) > 0 {
		trustedTenants := make([]kusto.TrustedExternalTenant, 0)

		for _, tenantId := range tenantIds.List() {
			trustedTenants = append(trustedTenants, kusto.TrustedExternalTenant{Value: utils.String(tenantId.(string))})
		}

		clusterProperties.TrustedExternalTenants = &trustedTenants
	}

	return &clusterProperties
}

func flattenKustoClusterSku(sku *kusto.AzureSku) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":     string((*sku).Name),
			"tier":     *(*sku).Tier,
			"capacity": utils.Int(int(*(*sku).Capacity)),
		},
	}
}
