package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryDatasetSnowflake() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetSnowflakeCreateUpdate,
		Read:   resourceDataFactoryDatasetSnowflakeRead,
		Update: resourceDataFactoryDatasetSnowflakeCreateUpdate,
		Delete: resourceDataFactoryDatasetSnowflakeDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryID,
			},

			"linked_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"table_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"schema_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"schema_column": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"NUMBER",
								"DECIMAL",
								"NUMERIC",
								"INT",
								"INTEGER",
								"BIGINT",
								"SMALLINT",
								"FLOAT",
								"FLOAT4",
								"FLOAT8",
								"DOUBLE",
								"DOUBLE PRECISION",
								"REAL",
								"VARCHAR",
								"CHAR",
								"CHARACTER",
								"STRING",
								"TEXT",
								"BINARY",
								"VARBINARY",
								"BOOLEAN",
								"DATE",
								"DATETIME",
								"TIME",
								"TIMESTAMP",
								"TIMESTAMP_LTZ",
								"TIMESTAMP_NTZ",
								"TIMESTAMP_TZ",
								"VARIANT",
								"OBJECT",
								"ARRAY",
								"GEOGRAPHY",
							}, false),
						},
						"precision": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"scale": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},
		},
	}
}

func resourceDataFactoryDatasetSnowflakeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).DataFactory.DatasetClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := parse.DataFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_snowflake", id.ID())
		}
	}

	snowflakeDatasetProperties := datafactory.SnowflakeDatasetTypeProperties{
		Table:  d.Get("table_name").(string),
		Schema: d.Get("schema_name").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	snowflakeTableset := datafactory.SnowflakeDataset{
		SnowflakeDatasetTypeProperties: &snowflakeDatasetProperties,
		LinkedServiceName:              linkedService,
		Description:                    &description,
		Schema:                         make([]interface{}, 0),
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		snowflakeTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		snowflakeTableset.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	annotations := d.Get("annotations").([]interface{})
	snowflakeTableset.Annotations = &annotations

	if v, ok := d.GetOk("additional_properties"); ok {
		snowflakeTableset.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("schema_column"); ok {
		snowflakeTableset.Schema = expandDataFactoryDatasetSnowflakeSchema(v.([]interface{}))
	}

	datasetType := string(datafactory.TypeBasicDatasetTypeSnowflakeTable)
	dataset := datafactory.DatasetResource{
		Properties: &snowflakeTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryDatasetSnowflakeRead(d, meta)
}

func resourceDataFactoryDatasetSnowflakeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := parse.NewDataFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Data Factory Dataset Snowflake %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	snowflakeTable, ok := resp.Properties.AsSnowflakeDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset Snowflake %s: Expected: %q Received: %q", *id, datafactory.TypeBasicDatasetTypeSnowflakeTable, *resp.Type)
	}

	d.Set("additional_properties", snowflakeTable.AdditionalProperties)

	if snowflakeTable.Description != nil {
		d.Set("description", snowflakeTable.Description)
	}

	parameters := flattenDataFactoryParameters(snowflakeTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(snowflakeTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := snowflakeTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := snowflakeTable.SnowflakeDatasetTypeProperties; properties != nil {
		tableName, ok := properties.Table.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `table_name` since it's not a string")
		} else {
			d.Set("table_name", tableName)
		}
		schemaName, ok := properties.Schema.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `schema_name` since it's not a string")
		} else {
			d.Set("schema_name", schemaName)
		}
	}

	if folder := snowflakeTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	schemaColumns := flattenDataFactorySnowflakeSchemaColumns(snowflakeTable.Schema)
	if err := d.Set("schema_column", schemaColumns); err != nil {
		return fmt.Errorf("Error setting `schema_column`: %+v", err)
	}

	return nil
}

func resourceDataFactoryDatasetSnowflakeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
