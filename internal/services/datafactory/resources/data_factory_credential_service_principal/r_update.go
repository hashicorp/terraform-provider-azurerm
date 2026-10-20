// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_service_principal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Credentials
			id, err := credentials.ParseCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var data DataFactoryCredentialServicePrincipalResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.CredentialOperationsGet(ctx, *id, credentials.DefaultCredentialOperationsGetOperationOptions())
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			props, ok := existing.Model.Properties.(credentials.ServicePrincipalCredential)
			if !ok {
				return fmt.Errorf("retrieving %s: expected `credentials.ServicePrincipalCredential` but got %T", id, existing.Model.Properties)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = &data.Description
			}

			if metadata.ResourceData.HasChange("annotations") {
				if len(data.Annotations) > 0 {
					annotations := make([]any, len(data.Annotations))
					for i, v := range data.Annotations {
						annotations[i] = v
					}
					props.Annotations = &annotations
				} else {
					props.Annotations = nil
				}
			}

			if metadata.ResourceData.HasChange("service_principal_key") {
				props.TypeProperties.ServicePrincipalKey = expandDataFactoryCredentialKeyVaultSecretReference(data.ServicePrincipalKey)
			}

			if metadata.ResourceData.HasChange("service_principal_id") {
				props.TypeProperties.ServicePrincipalId = pointer.To(data.ServicePrincipalId)
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				props.TypeProperties.Tenant = pointer.To(data.TenantId)
			}

			payload := credentials.CredentialResource{
				Properties: props,
			}
			if _, err = client.CredentialOperationsCreateOrUpdate(ctx, *id, payload, credentials.DefaultCredentialOperationsCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
