// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient

			id, err := linkedservices.ParseLinkedServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config LinkedServiceSqlManagedInstanceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id, linkedservices.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			props := existing.Model.Properties.(linkedservices.AzureSqlMILinkedService)
			typeProps := props.TypeProperties

			if metadata.ResourceData.HasChanges("connection_string", "key_vault_connection_string") {
				typeProps.ConnectionString = pointer.To(any(config.ConnectionString))

				if len(config.KeyVaultConnectionString) > 0 {
					typeProps.ConnectionString = pointer.To(expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(config.KeyVaultConnectionString))
				}
			}

			if metadata.ResourceData.HasChange("key_vault_password") {
				typeProps.Password = expandLinkedServiceSqlManagedInstanceKeyVaultPassword(config.KeyVaultPassword)
			}

			if metadata.ResourceData.HasChange("service_principal_id") {
				typeProps.ServicePrincipalId = pointer.To(any(config.ServicePrincipalID))
			}

			if metadata.ResourceData.HasChange("service_principal_key") {
				secureString := expandLinkedServiceSqlManagedInstanceServicePrincipalKey(config.ServicePrincipalKey)
				typeProps.ServicePrincipalKey = secureString
				typeProps.ServicePrincipalCredential = secureString
			} else if config.ServicePrincipalKey == "" {
				typeProps.ServicePrincipalKey = nil
				typeProps.ServicePrincipalCredential = nil
			}

			if metadata.ResourceData.HasChange("tenant") {
				typeProps.Tenant = pointer.To(any(config.Tenant))
			}

			if metadata.ResourceData.HasChange("parameters") {
				props.Parameters = expandLinkedServiceSqlManagedInstanceParameters(config.Parameters)
			}

			if metadata.ResourceData.HasChange("integration_runtime_name") {
				props.ConnectVia = expandLinkedServiceSqlManagedInstanceIntegrationRuntimeName(config.IntegrationRuntimeName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("annotations") {
				props.Annotations = expandLinkedServiceSqlManagedInstanceAnnotations(config.Annotations)
			}

			props.TypeProperties = typeProps
			payload := linkedservices.LinkedServiceResource{
				Properties: &props,
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
