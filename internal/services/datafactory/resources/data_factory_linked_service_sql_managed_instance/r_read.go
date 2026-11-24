// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			sqlMILinkedService, ok := resp.Model.Properties.(linkedservices.AzureSqlMILinkedService)
			if !ok {
				return fmt.Errorf("classifying %s: Expected: %q Received: %T", id, "AzureSqlMI", resp.Model.Properties)
			}

			if sqlMILinkedService.Type != "AzureSqlMI" {
				return fmt.Errorf("classifying %s: Expected: %q Received: %q", id, "AzureSqlMI", sqlMILinkedService.Type)
			}

			props := sqlMILinkedService.TypeProperties

			state := LinkedServiceSqlManagedInstanceModel{
				Name:                id.LinkedServiceName,
				DataFactoryID:       factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID(),
				Annotations:         flattenLinkedServiceSqlManagedInstanceAnnotations(sqlMILinkedService.Annotations),
				Description:         pointer.From(sqlMILinkedService.Description),
				KeyVaultPassword:    flattenLinkedServiceSqlManagedInstanceKeyVaultPassword(props.Password),
				Parameters:          flattenLinkedServiceSqlManagedInstanceParameters(sqlMILinkedService.Parameters),
				ServicePrincipalKey: metadata.ResourceData.Get("service_principal_key").(string),
			}

			if props.ConnectionString != nil {
				val := pointer.From(props.ConnectionString)
				switch v := val.(type) {
				case map[string]any:
					state.KeyVaultConnectionString = flattenLinkedServiceSqlManagedInstanceKeyVaultConnectionString(v)
				case string:
					state.ConnectionString = v
				default:
					return fmt.Errorf("classifying %s: expected a map or a string but received: %q", id, v)
				}
			}

			if str, ok := pointer.From(props.ServicePrincipalId).(string); ok {
				state.ServicePrincipalID = str
			}

			if str, ok := pointer.From(props.Tenant).(string); ok {
				state.Tenant = str
			}

			if connectVia := sqlMILinkedService.ConnectVia; connectVia != nil {
				state.IntegrationRuntimeName = connectVia.ReferenceName
			}

			return metadata.Encode(&state)
		},
	}
}

type KeyVaultConnectionStringConfig struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
}

type KeyVaultPasswordConfig struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
}
