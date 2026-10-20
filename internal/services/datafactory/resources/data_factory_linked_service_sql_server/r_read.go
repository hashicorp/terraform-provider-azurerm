// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_server

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

func resourceDataFactoryLinkedServiceSQLServerRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory SQL Server %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	sqlServer, ok := resp.Properties.AsSQLServerLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory SQL Server %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeSQLServer, *resp.Type)
	}

	d.Set("additional_properties", sqlServer.AdditionalProperties)
	d.Set("description", sqlServer.Description)

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(sqlServer.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(sqlServer.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := sqlServer.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	if properties := sqlServer.SQLServerLinkedServiceTypeProperties; properties != nil {
		d.Set("user_name", properties.UserName)
		if properties.ConnectionString != nil {
			if val, ok := properties.ConnectionString.(map[string]any); ok {
				if err := d.Set("key_vault_connection_string", helpers.FlattenAzureKeyVaultConnectionString(val)); err != nil {
					return fmt.Errorf("setting `key_vault_connection_string`: %+v", err)
				}
			} else if val, ok := properties.ConnectionString.(string); ok {
				d.Set("connection_string", val)
			} else {
				return fmt.Errorf("setting `connection_string`: %+v", err)
			}
		}

		if password := properties.Password; password != nil {
			if keyVaultPassword, ok := password.AsAzureKeyVaultSecretReference(); ok {
				if err := d.Set("key_vault_password", helpers.FlattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
					return fmt.Errorf("setting `key_vault_password`: %+v", err)
				}
			}
		}
	}

	return nil
}
