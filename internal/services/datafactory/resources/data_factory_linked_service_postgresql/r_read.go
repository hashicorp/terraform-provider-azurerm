// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_postgresql

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServicePostgreSQLRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
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

		return fmt.Errorf("retrieving Data Factory PostgreSQL %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	postgresql, ok := resp.Properties.AsPostgreSQLLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Linked Service PostgreSQL %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeBasicLinkedServiceTypePostgreSQL, *resp.Type)
	}

	d.Set("additional_properties", postgresql.AdditionalProperties)
	d.Set("description", postgresql.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(postgresql.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(postgresql.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := postgresql.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}
