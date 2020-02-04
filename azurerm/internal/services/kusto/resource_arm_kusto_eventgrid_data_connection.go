package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-05-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoEventGridDataConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoEventGridDataConnectionCreateUpdate,
		Read:   resourceArmKustoEventGridDataConnectionRead,
		Update: resourceArmKustoEventGridDataConnectionCreateUpdate,
		Delete: resourceArmKustoEventGridDataConnectionDelete,

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
				ValidateFunc: validateAzureRMKustoDataConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoDatabaseName,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"consumer_group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubConsumerName(),
			},

			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMKustoEntityName,
			},

			"mapping_rule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMKustoEntityName,
			},

			"data_format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.MULTIJSON),
					string(kusto.JSON),
					string(kusto.CSV),
					string(kusto.TSV),
					string(kusto.SCSV),
					string(kusto.SOHSV),
					string(kusto.PSV),
					string(kusto.TXT),
					string(kusto.RAW),
					string(kusto.SINGLEJSON),
					string(kusto.AVRO),
					string(kusto.TSVE),
				}, false),
			},
		},
	}
}

func resourceArmKustoEventGridDataConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto EventGrid Data Connection creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		connectionModel, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(connectionModel.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %s", name, resourceGroup, clusterName, databaseName, err)
			}
		}

		if dataConnection, ok := connectionModel.Value.(kusto.EventGridDataConnection); ok {
			if dataConnection.ID != nil && *dataConnection.ID != "" {
				return tf.ImportAsExistsError("azurerm_kusto_eventgrid_data_connection", *dataConnection.ID)
			}
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	eventGridDataConnectionProperties := expandKustoEventGridDataConnectionProperties(d)

	dataConnection1 := kusto.EventGridDataConnection{
		Name:                          &name,
		Location:                      &location,
		EventGridConnectionProperties: eventGridDataConnectionProperties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, databaseName, name, dataConnection1)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	connectionModel, getDetailsErr := client.Get(ctx, resourceGroup, clusterName, databaseName, name)

	if getDetailsErr != nil {
		return fmt.Errorf("Error retrieving Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if dataConnection, ok := connectionModel.Value.(kusto.EventGridDataConnection); ok {
		if dataConnection.ID == nil {
			return fmt.Errorf("Cannot read ID for Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
		}

		d.SetId(*dataConnection.ID)
	}

	return resourceArmKustoEventGridDataConnectionRead(d, meta)
}

func resourceArmKustoEventGridDataConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	name := id.Path["DataConnections"]

	connectionModel, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)

	if err != nil {
		if utils.ResponseWasNotFound(connectionModel.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_name", clusterName)
	d.Set("database_name", databaseName)

	if dataConnection, ok := connectionModel.Value.(kusto.EventGridDataConnection); ok {
		if location := dataConnection.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := dataConnection.EventGridConnectionProperties; props != nil {
			d.Set("storage_account_id", props.StorageAccountResourceID)
			d.Set("eventhub_id", props.EventHubResourceID)
			d.Set("consumer_group", props.ConsumerGroup)
			d.Set("table_name", props.TableName)
			d.Set("mapping_rule_name", props.MappingRuleName)
			d.Set("data_format", props.DataFormat)
		}
	}

	return nil
}

func resourceArmKustoEventGridDataConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	name := id.Path["DataConnections"]

	future, err := client.Delete(ctx, resourceGroup, clusterName, databaseName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto EventGrid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	return nil
}


func expandKustoEventGridDataConnectionProperties(d *schema.ResourceData) *kusto.EventGridConnectionProperties {
	eventGridConnectionProperties := &kusto.EventGridConnectionProperties{}

	if sotrageAccountResourceId, ok := d.GetOk("storage_account_id"); ok {
		eventGridConnectionProperties.StorageAccountResourceID = utils.String(sotrageAccountResourceId.(string))
	}

	if eventhubResourceId, ok := d.GetOk("eventhub_id"); ok {
		eventGridConnectionProperties.EventHubResourceID = utils.String(eventhubResourceId.(string))
	}

	if consumerGroup, ok := d.GetOk("consumer_group"); ok {
		eventGridConnectionProperties.ConsumerGroup = utils.String(consumerGroup.(string))
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		eventGridConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		eventGridConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		eventGridConnectionProperties.DataFormat = kusto.DataFormat(df.(string))
	}

	return eventGridConnectionProperties
}
