package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	eventhubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoEventGridDataConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoEventGridDataConnectionCreateUpdate,
		Read:   resourceArmKustoEventGridDataConnectionRead,
		Update: resourceArmKustoEventGridDataConnectionCreateUpdate,
		Delete: resourceArmKustoEventGridDataConnectionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.KustoDataConnectionID(id)
			return err
		}, importDataConnection(kusto.KindEventGrid)),

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
				ValidateFunc: validate.KustoDataConnectionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KustoClusterName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KustoDatabaseName,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"eventhub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.EventHubID,
			},

			"consumer_group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubConsumerName(),
			},

			"table_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.KustoEntityName,
			},

			"mapping_rule_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.KustoEntityName,
			},

			"data_format": {
				Type:     schema.TypeString,
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

func resourceArmKustoEventGridDataConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Grid Data Connection creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)

	if d.IsNewResource() {
		connectionModel, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(connectionModel.Response) {
				return fmt.Errorf("checking for presence of existing Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %s", name, resourceGroup, clusterName, databaseName, err)
			}
		}

		if dataConnection, ok := connectionModel.Value.(kusto.EventGridDataConnection); ok {
			if dataConnection.ID != nil && *dataConnection.ID != "" {
				return tf.ImportAsExistsError("azurerm_kusto_eventgrid_data_connection", *dataConnection.ID)
			}
		}
	}

	dataConnection := kusto.EventGridDataConnection{
		Name:     &name,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		EventGridConnectionProperties: &kusto.EventGridConnectionProperties{
			StorageAccountResourceID: utils.String(d.Get("storage_account_id").(string)),
			EventHubResourceID:       utils.String(d.Get("eventhub_id").(string)),
			ConsumerGroup:            utils.String(d.Get("consumer_group").(string)),
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

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, databaseName, name, dataConnection)
	if err != nil {
		return fmt.Errorf("creating or updating Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}
	dataConnection, ok := resp.Value.(kusto.EventGridDataConnection)
	if !ok {
		return fmt.Errorf("expected Type of Kusto Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %q but actually not", name, resourceGroup, clusterName, databaseName, kusto.KindEventGrid)
	}
	if dataConnection.ID == nil || *dataConnection.ID == "" {
		return fmt.Errorf("cannot read Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database: %q) ID", name, resourceGroup, clusterName, databaseName)
	}

	d.SetId(*dataConnection.ID)

	return resourceArmKustoEventGridDataConnectionRead(d, meta)
}

func resourceArmKustoEventGridDataConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	connectionModel, err := client.Get(ctx, id.ResourceGroup, id.Cluster, id.Database, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(connectionModel.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.Cluster)
	d.Set("database_name", id.Database)

	if dataConnection, ok := connectionModel.Value.(kusto.EventGridDataConnection); ok {
		d.Set("location", location.NormalizeNilable(dataConnection.Location))
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

	id, err := parse.KustoDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Cluster, id.Database, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Kusto Event Grid Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	return nil
}
