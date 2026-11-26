// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_snowflake

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryDatasetSnowflakeCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).DataFactory.DatasetClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_dataset_snowflake", id.ID())
			}
		}
	}

	snowflakeDatasetProperties := datafactory.SnowflakeDatasetTypeProperties{
		Table:  d.Get("table_name").(string),
		Schema: d.Get("schema_name").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          pointer.To("LinkedServiceReference"),
	}

	description := d.Get("description").(string)
	snowflakeTableset := datafactory.SnowflakeDataset{
		SnowflakeDatasetTypeProperties: &snowflakeDatasetProperties,
		LinkedServiceName:              linkedService,
		Description:                    &description,
		Schema:                         make([]any, 0),
	}

	if v, ok := d.GetOk("folder"); ok {
		snowflakeTableset.Folder = &datafactory.DatasetFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		snowflakeTableset.Parameters = helpers.ExpandDataSetParameters(v.(map[string]any))
	}

	annotations := d.Get("annotations").([]any)
	snowflakeTableset.Annotations = &annotations

	if v, ok := d.GetOk("additional_properties"); ok {
		snowflakeTableset.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("schema_column"); ok {
		snowflakeTableset.Schema = helpers.ExpandDataFactoryDatasetSnowflakeSchema(v.([]any))
	}

	dataset := datafactory.DatasetResource{
		Properties: &snowflakeTableset,
		Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeSnowflakeTable)),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryDatasetSnowflakeRead(d, meta)
}
