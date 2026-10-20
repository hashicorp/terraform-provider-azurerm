// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_service_principal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Credentials

			var data DataFactoryCredentialServicePrincipalResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := credentials.ParseFactoryID(data.DataFactoryId)
			if err != nil {
				return err
			}

			id := credentials.NewCredentialID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, data.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CredentialOperationsGet(ctx, id, credentials.DefaultCredentialOperationsGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_data_factory_credential_service_principal", id.ID())
				}
			}

			props := credentials.ServicePrincipalCredential{
				TypeProperties: credentials.ServicePrincipalCredentialTypeProperties{
					ServicePrincipalId:  pointer.To(data.ServicePrincipalId),
					ServicePrincipalKey: expandDataFactoryCredentialKeyVaultSecretReference(data.ServicePrincipalKey),
					Tenant:              pointer.To(data.TenantId),
				},
			}
			if len(data.Annotations) > 0 {
				annotations := make([]any, len(data.Annotations))
				for i, v := range data.Annotations {
					annotations[i] = v
				}
				props.Annotations = &annotations
			}
			if data.Description != "" {
				props.Description = &data.Description
			}

			payload := credentials.CredentialResource{
				Properties: props,
			}
			if _, err = client.CredentialOperationsCreateOrUpdate(ctx, id, payload, credentials.DefaultCredentialOperationsCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandDataFactoryCredentialKeyVaultSecretReference(input []ServicePrincipalKey) *credentials.AzureKeyVaultSecretReference {
	if len(input) == 0 {
		return nil
	}

	out := credentials.AzureKeyVaultSecretReference{
		SecretName: input[0].SecretName,
		Store: credentials.LinkedServiceReference{
			Type:          credentials.TypeLinkedServiceReference,
			ReferenceName: input[0].LinkedServiceName,
		},
	}

	if input[0].SecretVersion != "" {
		out.SecretVersion = pointer.To(input[0].SecretVersion)
	}

	return pointer.To(out)
}
