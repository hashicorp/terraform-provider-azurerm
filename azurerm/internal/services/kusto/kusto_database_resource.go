package kusto

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Database creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Database %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, clusterName, err)
			}
		}

		if existing.Value != nil {
			database, ok := existing.Value.AsReadWriteDatabase()
			if !ok {
				return fmt.Errorf("Exisiting Resource is not a Kusto Read/Write Database %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
			}

			if database.ID != nil && *database.ID != "" {
				return tf.ImportAsExistsError("azurerm_kusto_database", *database.ID)
			}
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	databaseProperties := expandKustoDatabaseProperties(d)

	readWriteDatabase := kusto.ReadWriteDatabase{
		Name:                        &name,
		Location:                    &location,
		ReadWriteDatabaseProperties: databaseProperties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, name, readWriteDatabase)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}
	if resp.Value == nil {
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): Invalid resource response", name, resourceGroup, clusterName)
	}

	database, ok := resp.Value.AsReadWriteDatabase()
	if !ok {
		return fmt.Errorf("Resource is not a Read/Write Database %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	}
	if database.ID == nil || *database.ID == "" {
		return fmt.Errorf("Cannot read ID for Kusto Database %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	}

	d.SetId(*database.ID)

	return resourceArmKustoDatabaseRead(d, meta)
}

func resourceArmKustoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, err)
	}

	if resp.Value == nil {
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): Invalid resource response", id.Name, id.ResourceGroup, id.ClusterName)
	}

	database, ok := resp.Value.AsReadWriteDatabase()
	if !ok {
		return fmt.Errorf("Existing resource is not a Read/Write Database (Resource Group %q, Cluster %q): %q", id.ResourceGroup, id.ClusterName, id.Name)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)

	if location := database.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := database.ReadWriteDatabaseProperties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}

func resourceArmKustoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

func expandKustoDatabaseProperties(d *schema.ResourceData) *kusto.ReadWriteDatabaseProperties {
	databaseProperties := &kusto.ReadWriteDatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = utils.String(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = utils.String(hotCachePeriod.(string))
	}

	return databaseProperties
}
