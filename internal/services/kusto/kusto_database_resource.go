package kusto

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	kustoValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabaseCreateUpdate,
		Read:   resourceKustoDatabaseRead,
		Update: resourceKustoDatabaseCreateUpdate,
		Delete: resourceKustoDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatabaseID(id)
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_kusto_database", id.ID())
		}
	}

	databaseProperties := expandKustoDatabaseProperties(d)

	readWriteDatabase := kusto.ReadWriteDatabase{
		Name:                        utils.String(id.Name),
		Location:                    utils.String(location.Normalize(d.Get("location").(string))),
		ReadWriteDatabaseProperties: databaseProperties,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ClusterName, id.Name, readWriteDatabase)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Value == nil {
		return fmt.Errorf("retrieving %s: response was nil", *id)
	}

	database, ok := resp.Value.AsReadWriteDatabase()
	if !ok {
		return fmt.Errorf("%s was not a Read/Write Database", *id)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)

	d.Set("location", location.NormalizeNilable(database.Location))

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
