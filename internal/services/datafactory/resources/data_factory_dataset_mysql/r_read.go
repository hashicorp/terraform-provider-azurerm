// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_mysql

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryDatasetMySQLRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	mysqlTable, ok := resp.Properties.AsRelationalTableDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset MySQL %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeRelationalTable, resp.Properties)
	}

	d.Set("additional_properties", mysqlTable.AdditionalProperties)

	if mysqlTable.Description != nil {
		d.Set("description", mysqlTable.Description)
	}

	if err := d.Set("parameters", helpers.FlattenDataSetParameters(mysqlTable.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(mysqlTable.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := mysqlTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := mysqlTable.RelationalTableDatasetTypeProperties; properties != nil {
		val, ok := properties.TableName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `table_name` since it's not a string")
		} else {
			d.Set("table_name", val)
		}
	}

	if folder := mysqlTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	if err := d.Set("schema_column", helpers.FlattenDataFactoryStructureColumns(mysqlTable.Structure)); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}
