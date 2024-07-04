// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataFactoryCredentialServicePrincipalResource struct{}

var _ sdk.Resource = DataFactoryCredentialServicePrincipalResource{}
var _ sdk.ResourceWithUpdate = DataFactoryCredentialServicePrincipalResource{}

func (DataFactoryCredentialServicePrincipalResource) ResourceType() string {
	return "azurerm_data_factory_credential_service_principal"
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

func (DataFactoryCredentialServicePrincipalResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The desired name of the credential resource",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_factory_id": {
			Description:  "The resource ID of the parent Data Factory",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: credentials.ValidateFactoryID,
		},

		"tenant_id": {
			Description:  "The Tenant ID of the Service Principal",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"service_principal_id": {
			Description:  "The Client ID of the Service Principal",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"service_principal_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"description": {
			Description:  "(Optional) Short text description",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"annotations": { // this property is not visible in the azure portal
			Description: "(Optional) List of string annotations.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (DataFactoryCredentialServicePrincipalResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DataFactoryCredentialServicePrincipalResource) ModelObject() interface{} {
	return &DataFactoryCredentialServicePrincipalResourceSchema{}
}

func (DataFactoryCredentialServicePrincipalResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return credentials.ValidateCredentialID
}

func (r DataFactoryCredentialServicePrincipalResource) Create() sdk.ResourceFunc {
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
			existing, err := client.CredentialOperationsGet(ctx, id, credentials.DefaultCredentialOperationsGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_credential_service_principal", id.ID())
			}

			props := credentials.ServicePrincipalCredential{
				TypeProperties: credentials.ServicePrincipalCredentialTypeProperties{
					ServicePrincipalId:  pointer.To(data.ServicePrincipalId),
					ServicePrincipalKey: expandDataFactoryCredentialKeyVaultSecretReference(data.ServicePrincipalKey),
					Tenant:              pointer.To(data.TenantId),
				},
			}
			if len(data.Annotations) > 0 {
				annotations := make([]interface{}, len(data.Annotations))
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

func (DataFactoryCredentialServicePrincipalResource) Read() sdk.ResourceFunc {
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
				state.Annotations = flattenDataFactoryAnnotations(props.Annotations)

				state.Description = pointer.From(props.Description)
				state.TenantId = pointer.From(props.TypeProperties.Tenant)
				state.ServicePrincipalId = pointer.From(props.TypeProperties.ServicePrincipalId)
				state.ServicePrincipalKey = flattenDataFactoryCredentialKeyVaultSecretReference(props.TypeProperties.ServicePrincipalKey)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataFactoryCredentialServicePrincipalResource) Update() sdk.ResourceFunc {
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
					annotations := make([]interface{}, len(data.Annotations))
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

func (DataFactoryCredentialServicePrincipalResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Credentials

			id, err := credentials.ParseCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.CredentialOperationsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

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
