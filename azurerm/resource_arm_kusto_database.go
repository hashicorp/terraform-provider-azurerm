package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-01-21/kusto"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoDatabaseCreateUpdate,
		Read:   resourceArmKustoDatabaseRead,
		Update: resourceArmKustoDatabaseCreateUpdate,
		Delete: resourceArmKustoDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoDatabaseName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"soft_delete_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"hot_cache_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceArmKustoDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Kusto Database creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		server, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Database %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, clusterName, err)
			}
		}

		if server.ID != nil && *server.ID != "" {
			return tf.ImportAsExistsError("azurerm_kusto_database", *server.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	databaseProperties := expandKustoDatabaseProperties(d)

	kustoDatabase := kusto.Database{
		Name:               &name,
		Location:           &location,
		DatabaseProperties: databaseProperties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, name, kustoDatabase)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto Cluster %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Cluster %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	resp, getDetailsErr := client.Get(ctx, resourceGroup, clusterName, name)
	if getDetailsErr != nil {
		return fmt.Errorf("Error retrieving Kusto Cluster %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Kusto Cluster %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	}

	d.SetId(*resp.ID)

	return resourceArmKustoDatabaseRead(d, meta)
}

func resourceArmKustoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	name := id.Path["Databases"]

	databaseResponse, err := client.Get(ctx, resourceGroup, clusterName, name)

	if err != nil {
		if utils.ResponseWasNotFound(databaseResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_name", clusterName)

	if location := databaseResponse.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := databaseResponse.DatabaseProperties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}

func resourceArmKustoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).kusto.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	name := id.Path["Databases"]

	future, err := client.Delete(ctx, resGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	}

	return nil
}

func validateAzureRMKustoDatabaseName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 260 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func expandKustoDatabaseProperties(d *schema.ResourceData) *kusto.DatabaseProperties {
	databaseProperties := &kusto.DatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = utils.String(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = utils.String(hotCachePeriod.(string))
	}

	return databaseProperties
}
