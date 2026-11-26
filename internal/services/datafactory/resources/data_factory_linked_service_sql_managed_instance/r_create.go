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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.LinkedServicesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config LinkedServiceSqlManagedInstanceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := factories.ParseFactoryID(config.DataFactoryID)
			if err != nil {
				return err
			}

			id := linkedservices.NewLinkedServiceID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			sqlMILinkedService := &linkedservices.AzureSqlMILinkedService{
				Description: pointer.To(config.Description),
				Parameters:  expandLinkedServiceSqlManagedInstanceParameters(config.Parameters),
				Annotations: expandLinkedServiceSqlManagedInstanceAnnotations(config.Annotations),
				TypeProperties: linkedservices.AzureSqlMILinkedServiceTypeProperties{
					Password: expandLinkedServiceSqlManagedInstanceKeyVaultPassword(config.KeyVaultPassword),
				},
			}

			if config.ConnectionString != "" {
				sqlMILinkedService.TypeProperties.ConnectionString = pointer.To(any(config.ConnectionString))
			}

			if len(config.KeyVaultConnectionString) > 0 {
				sqlMILinkedService.TypeProperties.ConnectionString = pointer.To(expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(config.KeyVaultConnectionString))
			}

			if config.ServicePrincipalID != "" {
				sqlMILinkedService.TypeProperties.ServicePrincipalId = pointer.To(any(config.ServicePrincipalID))
			}

			if config.ServicePrincipalKey != "" {
				secureString := expandLinkedServiceSqlManagedInstanceServicePrincipalKey(config.ServicePrincipalKey)
				sqlMILinkedService.TypeProperties.ServicePrincipalKey = secureString
				sqlMILinkedService.TypeProperties.ServicePrincipalCredential = secureString
			}

			if config.Tenant != "" {
				sqlMILinkedService.TypeProperties.Tenant = pointer.To(any(config.Tenant))
			}

			sqlMILinkedService.ConnectVia = expandLinkedServiceSqlManagedInstanceIntegrationRuntimeName(config.IntegrationRuntimeName)

			linkedService := linkedservices.LinkedServiceResource{
				Properties: sqlMILinkedService,
			}

			if _, err := client.CreateOrUpdate(ctx, id, linkedService, linkedservices.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandLinkedServiceSqlManagedInstanceKeyVaultConnectionString(input []KeyVaultConnectionStringConfig) any {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &linkedservices.AzureKeyVaultSecretReference{
		SecretName: config.SecretName,
		Store: linkedservices.LinkedServiceReference{
			Type:          linkedservices.TypeLinkedServiceReference,
			ReferenceName: config.LinkedServiceName,
		},
	}
}

func expandLinkedServiceSqlManagedInstanceKeyVaultPassword(input []KeyVaultPasswordConfig) *linkedservices.AzureKeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &linkedservices.AzureKeyVaultSecretReference{
		SecretName: config.SecretName,
		Store: linkedservices.LinkedServiceReference{
			Type:          linkedservices.TypeLinkedServiceReference,
			ReferenceName: config.LinkedServiceName,
		},
	}
}

func expandLinkedServiceSqlManagedInstanceParameters(input map[string]any) *map[string]linkedservices.ParameterSpecification {
	if len(input) == 0 {
		return nil
	}

	parameterSpec := make(map[string]linkedservices.ParameterSpecification)
	for key, value := range input {
		parameterSpec[key] = linkedservices.ParameterSpecification{
			Type:         linkedservices.ParameterTypeString,
			DefaultValue: pointer.To(value),
		}
	}
	return &parameterSpec
}

func expandLinkedServiceSqlManagedInstanceAnnotations(input []string) *[]any {
	if len(input) == 0 {
		return nil
	}

	annotations := make([]any, len(input))
	for i, v := range input {
		annotations[i] = v
	}
	return &annotations
}

func expandLinkedServiceSqlManagedInstanceServicePrincipalKey(input string) *linkedservices.SecureString {
	if input == "" {
		return nil
	}

	return &linkedservices.SecureString{
		Value: input,
	}
}

func expandLinkedServiceSqlManagedInstanceIntegrationRuntimeName(input string) *linkedservices.IntegrationRuntimeReference {
	if input == "" {
		return nil
	}

	return &linkedservices.IntegrationRuntimeReference{
		Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
		ReferenceName: input,
	}
}
