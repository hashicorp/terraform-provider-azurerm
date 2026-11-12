// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_sql_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryLinkedServiceAzureSQLDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
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

		return fmt.Errorf("retrieving Data Factory Azure SQL Database %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	sql, ok := resp.Properties.AsAzureSQLDatabaseLinkedService()
	if !ok {
		return fmt.Errorf("classifying Data Factory Azure SQL Database %s: Expected: %q Received: %q", *id, datafactory.TypeBasicLinkedServiceTypeAzureSQLDatabase, *resp.Type)
	}

	if sql != nil {
		if sql.Tenant != nil {
			d.Set("tenant_id", sql.Tenant)
		}

		if sql.ServicePrincipalID != nil {
			d.Set("service_principal_id", sql.ServicePrincipalID)
			d.Set("use_managed_identity", false)
		} else {
			d.Set("use_managed_identity", true)
		}
	}

	if sql.ConnectionString != nil {
		if val, ok := sql.ConnectionString.(map[string]any); ok {
			if val["type"] != "SecureString" {
				if err := d.Set("key_vault_connection_string", helpers.FlattenAzureKeyVaultConnectionString(val)); err != nil {
					return fmt.Errorf("setting `key_vault_connection_string`: %+v", err)
				}
			}
		}
	}

	d.Set("additional_properties", sql.AdditionalProperties)
	d.Set("description", sql.Description)

	if password := sql.Password; password != nil {
		if keyVaultPassword, ok := password.AsAzureKeyVaultSecretReference(); ok {
			if err := d.Set("key_vault_password", helpers.FlattenAzureKeyVaultSecretReference(keyVaultPassword)); err != nil {
				return fmt.Errorf("setting `key_vault_password`: %+v", err)
			}
		}
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(sql.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(sql.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if connectVia := sql.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	if credential := sql.Credential; credential != nil {
		d.Set("credential_name", pointer.From(credential.ReferenceName))
	}

	return nil
}
