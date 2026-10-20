// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_user_managed_identity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.Credentials

			id, err := credentials.ParseCredentialID(d.Id())
			if err != nil {
				return err
			}

			existing, err := client.CredentialOperationsGet(ctx, *id, credentials.DefaultCredentialOperationsGetOperationOptions())
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DataFactoryCredentialUserAssignedManagedIdentityResourceSchema{
				Name:          id.CredentialName,
				DataFactoryId: factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID(),
			}

			if model := existing.Model; model != nil {
				props, ok := model.Properties.(credentials.ManagedIdentityCredential)
				if !ok {
					return fmt.Errorf("retrieving %s: expected `credentials.ManagedIdentityCredential` but got %T", id, model.Properties)
				}

				if props.Description != nil {
					state.Description = *props.Description
				}

				if props.TypeProperties.ResourceId != nil {
					state.IdentityId = *props.TypeProperties.ResourceId
				}

				state.Annotations = helpers.FlattenDataFactoryAnnotations(props.Annotations)
			}

			return metadata.Encode(&state)
		},
	}
}

type DataFactoryCredentialUserAssignedManagedIdentityResourceSchema struct {
	Name          string   `tfschema:"name"`
	DataFactoryId string   `tfschema:"data_factory_id"`
	IdentityId    string   `tfschema:"identity_id"`
	Description   string   `tfschema:"description"`
	Annotations   []string `tfschema:"annotations"`
}
