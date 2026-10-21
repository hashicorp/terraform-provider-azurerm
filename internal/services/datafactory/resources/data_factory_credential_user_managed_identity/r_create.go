// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_user_managed_identity

import (
	"context"
	"fmt"
	"time"

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

			var data DataFactoryCredentialUserAssignedManagedIdentityResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			dataFactoryId, err := credentials.ParseFactoryID(data.DataFactoryId)
			if err != nil {
				return err
			}

			id := credentials.NewCredentialID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, data.Name)
			existing, err := client.CredentialOperationsGet(ctx, id, credentials.DefaultCredentialOperationsGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_credential_user_managed_identity", id.ID())
			}

			props := credentials.ManagedIdentityCredential{
				TypeProperties: &credentials.ManagedIdentityTypeProperties{
					ResourceId: &data.IdentityId,
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
