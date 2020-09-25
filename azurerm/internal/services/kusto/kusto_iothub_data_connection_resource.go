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
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoIotHubDataConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoIotHubDataConnectionCreateUpdate,
		Read:   resourceArmKustoIotHubDataConnectionRead,
		Update: resourceArmKustoIotHubDataConnectionCreateUpdate,
		Delete: resourceArmKustoIotHubDataConnectionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.KustoDataConnectionID(id)
			return err
		}, importDataConnection(kusto.KindIotHub)),

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

			"iothub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubID,
			},

			"consumer_group": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubConsumerGroupName,
			},

			"shared_access_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubSharedAccessPolicyName,
			},

			"event_system_properties": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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

func resourceArmKustoIotHubDataConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Iot Hub Data Connection creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)

	if d.IsNewResource() {
		connectionModel, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(connectionModel.Response) {
				return fmt.Errorf("checking for presence of existing Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %s", name, resourceGroup, clusterName, databaseName, err)
			}
		}

		if dataConnection, ok := connectionModel.Value.(kusto.EventHubDataConnection); ok {
			if dataConnection.ID != nil && *dataConnection.ID != "" {
				return tf.ImportAsExistsError("azurerm_kusto_iothub_data_connection", *dataConnection.ID)
			}
		}
	}

	dataConnection := kusto.IotHubDataConnection{
		Name:     &name,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		IotHubConnectionProperties: &kusto.IotHubConnectionProperties{
			IotHubResourceID:       utils.String(d.Get("iothub_id").(string)),
			ConsumerGroup:          utils.String(d.Get("consumer_group").(string)),
			SharedAccessPolicyName: utils.String(d.Get("shared_access_policy_name").(string)),
		},
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		dataConnection.IotHubConnectionProperties.EventSystemProperties = utils.ExpandStringSlice(eventSystemProperties.(*schema.Set).List())
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		dataConnection.IotHubConnectionProperties.TableName = utils.String(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		dataConnection.IotHubConnectionProperties.MappingRuleName = utils.String(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataConnection.IotHubConnectionProperties.DataFormat = kusto.IotHubDataFormat(df.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, databaseName, name, dataConnection)
	if err != nil {
		return fmt.Errorf("creating or updating Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}
	dataConnection, ok := resp.Value.(kusto.IotHubDataConnection)
	if !ok {
		return fmt.Errorf("expected Type of Kusto Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %q but actually not", name, resourceGroup, clusterName, databaseName, kusto.KindIotHub)
	}
	if dataConnection.ID == nil || *dataConnection.ID == "" {
		return fmt.Errorf("cannot read Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q) ID", name, resourceGroup, clusterName, databaseName)
	}

	d.SetId(*dataConnection.ID)

	return resourceArmKustoIotHubDataConnectionRead(d, meta)
}

func resourceArmKustoIotHubDataConnectionRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.Cluster)
	d.Set("database_name", id.Database)

	if dataConnection, ok := connectionModel.Value.(kusto.IotHubDataConnection); ok {
		if location := dataConnection.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := dataConnection.IotHubConnectionProperties; props != nil {
			d.Set("iothub_id", props.IotHubResourceID)
			d.Set("consumer_group", props.ConsumerGroup)
			d.Set("shared_access_policy_name", props.SharedAccessPolicyName)
			d.Set("event_system_properties", utils.FlattenStringSlice(props.EventSystemProperties))
			d.Set("table_name", props.TableName)
			d.Set("mapping_rule_name", props.MappingRuleName)
			d.Set("data_format", props.DataFormat)
		}
	}

	return nil
}

func resourceArmKustoIotHubDataConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KustoDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Cluster, id.Database, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Kusto Iot Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.Cluster, id.Database, err)
	}

	return nil
}
