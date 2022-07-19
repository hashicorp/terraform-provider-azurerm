package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
		}, importDataConnection(kusto.KindBasicDataConnectionKindEventGrid)),

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
				ValidateFunc: eventhubs.ValidateEventhubID,
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
				Default:  string(kusto.BlobStorageEventTypeMicrosoftStorageBlobCreated),
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.BlobStorageEventTypeMicrosoftStorageBlobCreated),
					string(kusto.BlobStorageEventTypeMicrosoftStorageBlobRenamed),
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
					string(kusto.EventGridDataFormatAPACHEAVRO),
					string(kusto.EventGridDataFormatAVRO),
					string(kusto.EventGridDataFormatCSV),
					string(kusto.EventGridDataFormatJSON),
					string(kusto.EventGridDataFormatMULTIJSON),
					string(kusto.EventGridDataFormatORC),
					string(kusto.EventGridDataFormatPARQUET),
					string(kusto.EventGridDataFormatPSV),
					string(kusto.EventGridDataFormatRAW),
					string(kusto.EventGridDataFormatSCSV),
					string(kusto.EventGridDataFormatSINGLEJSON),
					string(kusto.EventGridDataFormatSOHSV),
					string(kusto.EventGridDataFormatTSV),
					string(kusto.EventGridDataFormatTSVE),
					string(kusto.EventGridDataFormatTXT),
					string(kusto.EventGridDataFormatW3CLOGFILE),
				}, false),
			},

			"database_routing_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(kusto.DatabaseRoutingSingle),
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.DatabaseRoutingSingle),
					string(kusto.DatabaseRoutingMulti),
				}, false),
			},

			"eventgrid_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"managed_identity_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validate.ClusterID,
					commonids.ValidateUserAssignedIdentityID,
				),
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

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dataConnection.DatabaseRouting = kusto.DatabaseRouting(databaseRouting.(string))
	}

	if eventGridRID, ok := d.GetOk("eventgrid_resource_id"); ok {
		dataConnection.EventGridConnectionProperties.EventGridResourceID = utils.String(eventGridRID.(string))
	}

	if managedIdentityRID, ok := d.GetOk("managed_identity_resource_id"); ok {
		dataConnection.EventGridConnectionProperties.ManagedIdentityResourceID = utils.String(managedIdentityRID.(string))
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
			d.Set("database_routing_type", props.DatabaseRouting)
			d.Set("eventgrid_resource_id", props.EventGridResourceID)
			d.Set("managed_identity_resource_id", props.ManagedIdentityResourceID)
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
