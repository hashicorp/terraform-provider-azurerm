// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_credential_service_principal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
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

			state := DataFactoryCredentialServicePrincipalResourceSchema{
				Name:          id.CredentialName,
				DataFactoryId: credentials.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID(),
			}

			if model := existing.Model; model != nil {
				props, ok := model.Properties.(credentials.ServicePrincipalCredential)
				if !ok {
					return fmt.Errorf("retrieving %s: expected `credentials.ServicePrincipalCredential` but got %T", id, model.Properties)
				}

				state.Description = pointer.From(props.Description)
				state.Annotations = helpers.FlattenDataFactoryAnnotations(props.Annotations)

				state.Description = pointer.From(props.Description)
				state.TenantId = pointer.From(props.TypeProperties.Tenant)
				state.ServicePrincipalId = pointer.From(props.TypeProperties.ServicePrincipalId)
				state.ServicePrincipalKey = flattenDataFactoryCredentialKeyVaultSecretReference(props.TypeProperties.ServicePrincipalKey)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenDataFactoryCredentialKeyVaultSecretReference(input *credentials.AzureKeyVaultSecretReference) []ServicePrincipalKey {
	if input == nil {
		return []ServicePrincipalKey{}
	}

	return []ServicePrincipalKey{
		{
			LinkedServiceName: input.Store.ReferenceName,
			SecretName:        input.SecretName,
			SecretVersion:     pointer.From(input.SecretVersion),
		},
	}
}

type DataFactoryCredentialServicePrincipalResourceSchema struct {
	Name                string                `tfschema:"name"`
	DataFactoryId       string                `tfschema:"data_factory_id"`
	TenantId            string                `tfschema:"tenant_id"`
	ServicePrincipalId  string                `tfschema:"service_principal_id"`
	ServicePrincipalKey []ServicePrincipalKey `tfschema:"service_principal_key"`
	Description         string                `tfschema:"description"`
	Annotations         []string              `tfschema:"annotations"`
}

type ServicePrincipalKey struct {
	LinkedServiceName string `tfschema:"linked_service_name"`
	SecretName        string `tfschema:"secret_name"`
	SecretVersion     string `tfschema:"secret_version"`
}
