package kusto

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoEventHubDataConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoEventHubDataConnectionCreateUpdate,
		Read:   resourceArmKustoEventHubDataConnectionRead,
		Update: resourceArmKustoEventHubDataConnectionCreateUpdate,
		Delete: resourceArmKustoEventHubDataConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			// TODO: confirm these
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
				Optional:     true,
				ValidateFunc: validateAzureRMKustoEntityName,
			},

			"mapping_rule_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAzureRMKustoEntityName,
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
		CustomizeDiff: func(d *schema.ResourceDiff, _ interface{}) error {
			_, hasTableName := d.GetOk("table_name")
			_, hasMappingRuleName := d.GetOk("mapping_rule_name")
			_, hasDataFormat := d.GetOk("data_format")

			if !(utils.AllEquals(hasTableName, hasMappingRuleName, hasDataFormat)) {
				return fmt.Errorf("if one of the target table properties `table_name`, `mapping_rule_name` or `data_format` are set, the other values must also be defined")
			}

			return nil
		},
	}
}

func resourceArmKustoEventHubDataConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Event Hub Data Connection creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)

	if d.IsNewResource() {
		connectionModel, err := client.Get(ctx, resourceGroup, clusterName, databaseName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(connectionModel.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %s", name, resourceGroup, clusterName, databaseName, err)
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
		Name:                         &name,
		Location:                     &location,
		EventHubConnectionProperties: eventHubDataConnectionProperties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, databaseName, name, dataConnection1)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	connectionModel, getDetailsErr := client.Get(ctx, resourceGroup, clusterName, databaseName, name)

	if getDetailsErr != nil {
		return fmt.Errorf("Error retrieving Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
	}

	if dataConnection, ok := connectionModel.Value.(kusto.EventHubDataConnection); ok {
		if dataConnection.ID == nil {
			return fmt.Errorf("Cannot read ID for Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database: %q): %+v", name, resourceGroup, clusterName, databaseName, err)
		}

		d.SetId(*dataConnection.ID)
	}

	return resourceArmKustoEventHubDataConnectionRead(d, meta)
}

func resourceArmKustoEventHubDataConnectionRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
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
		}
	}

	return nil
}

func resourceArmKustoEventHubDataConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.Name, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}

func validateAzureRMKustoDataConnectionName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters, digits, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 40 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 1 and 40 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func validateAzureRMKustoEntityName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters, digits, underscores, spaces, dashes and dots: %q", k, name))
	}

	if len(name) > 1024 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 1 and 1024 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func expandKustoEventHubDataConnectionProperties(d *schema.ResourceData) *kusto.EventHubConnectionProperties {
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

	return eventHubConnectionProperties
}
