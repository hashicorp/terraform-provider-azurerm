// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_sql_server_table

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

func resourceDataFactoryDatasetSQLServerTableRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory Dataset SQL Server Table %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	sqlServerTable, ok := resp.Properties.AsSQLServerTableDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset SQL Server Table %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeSQLServerTable, resp.Properties)
	}

	d.Set("additional_properties", sqlServerTable.AdditionalProperties)

	if sqlServerTable.Description != nil {
		d.Set("description", sqlServerTable.Description)
	}

	if err := d.Set("parameters", helpers.FlattenDataSetParameters(sqlServerTable.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(sqlServerTable.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := sqlServerTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := sqlServerTable.SQLServerTableDatasetTypeProperties; properties != nil {
		val, ok := properties.TableName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `table_name` since it's not a string")
		} else {
			d.Set("table_name", val)
		}
	}

	if folder := sqlServerTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	if err := d.Set("schema_column", helpers.FlattenDataFactoryStructureColumns(sqlServerTable.Structure)); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}
