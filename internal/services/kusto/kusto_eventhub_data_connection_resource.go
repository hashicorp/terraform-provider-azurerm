package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/eventhubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoEventHubDataConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoEventHubDataConnectionCreateUpdate,
		Read:   resourceKustoEventHubDataConnectionRead,
		Update: resourceKustoEventHubDataConnectionCreateUpdate,
		Delete: resourceKustoEventHubDataConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataConnectionID(id)
			return err
		}, importDataConnection(kusto.KindBasicDataConnectionKindEventHub)),

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

			"compression": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  kusto.CompressionNone,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.CompressionGZip),
					string(kusto.CompressionNone),
				}, false),
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"eventhub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubs.ValidateEventhubID,
			},

			"event_system_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"consumer_group": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					eventhubValidate.ValidateEventHubConsumerName(),
					validation.StringInSlice([]string{"$Default"}, false)),
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.EntityName,
			},

			"identity_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					validate.ClusterID,
					commonids.ValidateUserAssignedIdentityID,
				),
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
					string(kusto.EventHubDataFormatAPACHEAVRO),
					string(kusto.EventHubDataFormatAVRO),
					string(kusto.EventHubDataFormatCSV),
					string(kusto.EventHubDataFormatJSON),
					string(kusto.EventHubDataFormatMULTIJSON),
					string(kusto.EventHubDataFormatPSV),
					string(kusto.EventHubDataFormatRAW),
					string(kusto.EventHubDataFormatSCSV),
					string(kusto.EventHubDataFormatSINGLEJSON),
					string(kusto.EventHubDataFormatSOHSV),
					string(kusto.EventHubDataFormatTSV),
					string(kusto.EventHubDataFormatTXT),
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
		},
	}
}

func resourceKustoEventHubDataConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Hub Data Connection creation.")

	id := parse.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		connectionModel, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(connectionModel.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if dataConnection, ok := connectionModel.Value.(kusto.EventHubDataConnection); ok {
			if dataConnection.ID != nil && *dataConnection.ID != "" {
				return tf.ImportAsExistsError("azurerm_kusto_eventhub_data_connection", *dataConnection.ID)
			}
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	eventHubDataConnectionProperties := expandKustoEventHubDataConnectionProperties(d)

	dataConnection1 := kusto.EventHubDataConnection{
		Name:                         &id.Name,
		Location:                     &location,
		EventHubConnectionProperties: eventHubDataConnectionProperties,
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dataConnection1.DatabaseRouting = kusto.DatabaseRouting(databaseRouting.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name, dataConnection1)
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceKustoEventHubDataConnectionRead(d, meta)
}

func resourceKustoEventHubDataConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectionID(d.Id())
	if err != nil {
		return err
	}

	connectionModel, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(connectionModel.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if dataConnection, ok := connectionModel.Value.(kusto.EventHubDataConnection); ok {
		if location := dataConnection.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := dataConnection.EventHubConnectionProperties; props != nil {
			d.Set("eventhub_id", props.EventHubResourceID)
			d.Set("consumer_group", props.ConsumerGroup)
			d.Set("table_name", props.TableName)
			d.Set("mapping_rule_name", props.MappingRuleName)
			d.Set("data_format", props.DataFormat)
			d.Set("database_routing_type", props.DatabaseRouting)
			d.Set("compression", props.Compression)
			d.Set("event_system_properties", props.EventSystemProperties)
			d.Set("identity_id", props.ManagedIdentityResourceID)
		}
	}

	return nil
}

func resourceKustoEventHubDataConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}

func expandKustoEventHubDataConnectionProperties(d *pluginsdk.ResourceData) *kusto.EventHubConnectionProperties {
	eventHubConnectionProperties := &kusto.EventHubConnectionProperties{}

	if eventhubResourceID, ok := d.GetOk("eventhub_id"); ok {
		eventHubConnectionProperties.EventHubResourceID = utils.String(eventhubResourceID.(string))
	}

	if consumerGroup, ok := d.GetOk("consumer_group"); ok {
		eventHubConnectionProperties.ConsumerGroup = utils.String(consumerGroup.(string))
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		eventHubConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		eventHubConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		eventHubConnectionProperties.DataFormat = kusto.EventHubDataFormat(df.(string))
	}

	if compression, ok := d.GetOk("compression"); ok {
		eventHubConnectionProperties.Compression = kusto.Compression(compression.(string))
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		props := make([]string, 0)
		for _, prop := range eventSystemProperties.([]interface{}) {
			props = append(props, prop.(string))
		}
		eventHubConnectionProperties.EventSystemProperties = &props
	}

	if identityId, ok := d.GetOk("identity_id"); ok {
		eventHubConnectionProperties.ManagedIdentityResourceID = utils.String(identityId.(string))
	}

	return eventHubConnectionProperties
}
