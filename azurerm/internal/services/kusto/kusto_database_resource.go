package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	kustoValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabaseCreateUpdate,
		Read:   resourceKustoDatabaseRead,
		Update: resourceKustoDatabaseCreateUpdate,
		Delete: resourceKustoDatabaseDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: kustoValidate.DatabaseName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: kustoValidate.ClusterName,
			},

			"soft_delete_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"hot_cache_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"size": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceKustoDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	return resourceKustoDatabaseRead(d, meta)
}

func resourceKustoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceKustoDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

func expandKustoDatabaseProperties(d *pluginsdk.ResourceData) *kusto.ReadWriteDatabaseProperties {
	databaseProperties := &kusto.ReadWriteDatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = utils.String(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = utils.String(hotCachePeriod.(string))
	}

	return databaseProperties
}
