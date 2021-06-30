package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	eventhubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKustoEventGridDataConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoEventGridDataConnectionCreateUpdate,
		Update: resourceKustoEventGridDataConnectionCreateUpdate,
		Read:   resourceKustoEventGridDataConnectionRead,
		Delete: resourceKustoEventGridDataConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataConnectionID(id)
			return err
		}, importDataConnection(kusto.KindEventGrid)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"eventhub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.EventHubID,
			},

			"eventhub_consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
			},

			"blob_storage_event_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(kusto.MicrosoftStorageBlobCreated),
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.MicrosoftStorageBlobCreated),
					string(kusto.MicrosoftStorageBlobRenamed),
				}, false),
			},

			"skip_first_record": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"mapping_rule_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"data_format": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.AVRO),
					string(kusto.CSV),
					string(kusto.JSON),
					string(kusto.MULTIJSON),
					string(kusto.PSV),
					string(kusto.RAW),
					string(kusto.SCSV),
					string(kusto.SINGLEJSON),
					string(kusto.SOHSV),
					string(kusto.TSV),
					string(kusto.TXT),
				}, false),
			},
		},
	}
}

func resourceKustoEventGridDataConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Grid Data Connection creation.")

	id := parse.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_kusto_eventgrid_data_connection", id.ID())
		}
	}

	dataConnection := kusto.EventGridDataConnection{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		EventGridConnectionProperties: &kusto.EventGridConnectionProperties{
			StorageAccountResourceID: utils.String(d.Get("storage_account_id").(string)),
			EventHubResourceID:       utils.String(d.Get("eventhub_id").(string)),
			ConsumerGroup:            utils.String(d.Get("eventhub_consumer_group_name").(string)),
			IgnoreFirstRecord:        utils.Bool(d.Get("skip_first_record").(bool)),
			BlobStorageEventType:     kusto.BlobStorageEventType(d.Get("blob_storage_event_type").(string)),
		},
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		dataConnection.EventGridConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		dataConnection.EventGridConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataConnection.EventGridConnectionProperties.DataFormat = kusto.EventGridDataFormat(df.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name, dataConnection)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoEventGridDataConnectionRead(d, meta)
}

func resourceKustoEventGridDataConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if dataConnection, ok := resp.Value.(kusto.EventGridDataConnection); ok {
		d.Set("location", location.NormalizeNilable(dataConnection.Location))
		if props := dataConnection.EventGridConnectionProperties; props != nil {
			d.Set("storage_account_id", props.StorageAccountResourceID)
			d.Set("eventhub_id", props.EventHubResourceID)
			d.Set("eventhub_consumer_group_name", props.ConsumerGroup)
			d.Set("skip_first_record", props.IgnoreFirstRecord)
			d.Set("blob_storage_event_type", props.BlobStorageEventType)
			d.Set("table_name", props.TableName)
			d.Set("mapping_rule_name", props.MappingRuleName)
			d.Set("data_format", props.DataFormat)
		}
	}

	return nil
}

func resourceKustoEventGridDataConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
