// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataFactoryCredentialUserAssignedManagedIdentityResource struct{}

// user managed identities only have one type
const IDENTITY_TYPE = "ManagedIdentity"

var _ sdk.Resource = DataFactoryCredentialUserAssignedManagedIdentityResource{}
var _ sdk.ResourceWithUpdate = DataFactoryCredentialUserAssignedManagedIdentityResource{}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) ResourceType() string {
	return "azurerm_data_factory_credential_user_managed_identity"
}

type DataFactoryCredentialUserAssignedManagedIdentityResourceSchema struct {
	Name          string   `tfschema:"name"`
	DataFactoryId string   `tfschema:"data_factory_id"`
	IdentityId    string   `tfschema:"identity_id"`
	Description   string   `tfschema:"description"`
	Annotations   []string `tfschema:"annotations"`
}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: factories.ValidateFactoryID,
		},
		"identity_id": {
			Description:  "The resource ID of the User Assigned Managed Identity",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
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

func (DataFactoryCredentialUserAssignedManagedIdentityResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) ModelObject() interface{} {
	return &DataFactoryCredentialUserAssignedManagedIdentityResourceSchema{}
}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return credentials.ValidateCredentialID
}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.Credentials

			credentialId, err := credentials.ParseCredentialID(d.Id())
			if err != nil {
				return err
			}

			var state DataFactoryCredentialUserAssignedManagedIdentityResourceSchema

			existing, err := client.CredentialOperationsGet(ctx, *credentialId, credentials.CredentialOperationsGetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", d.Id(), err)
			}

			state.Name = credentialId.CredentialName

			if existing.Model != nil {
				if existing.Model.Properties.Description != nil {
					state.Description = *existing.Model.Properties.Description
				}

				if existing.Model.Properties.TypeProperties != nil && existing.Model.Properties.TypeProperties.ResourceId != nil {
					state.IdentityId = *existing.Model.Properties.TypeProperties.ResourceId
				}
			}

			state.DataFactoryId = factories.NewFactoryID(credentialId.SubscriptionId, credentialId.ResourceGroupName, credentialId.FactoryName).ID()
			state.Annotations = flattenDataFactoryAnnotations(existing.Model.Properties.Annotations)

			return metadata.Encode(&state)
		},
	}
}

func (r DataFactoryCredentialUserAssignedManagedIdentityResource) Create() sdk.ResourceFunc {
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

			id := credentials.CredentialId{
				SubscriptionId:    dataFactoryId.SubscriptionId,
				ResourceGroupName: dataFactoryId.ResourceGroupName,
				FactoryName:       dataFactoryId.FactoryName,
				CredentialName:    data.Name,
			}

			existing, err := client.CredentialOperationsGet(ctx, id, credentials.CredentialOperationsGetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_dataset_http", id.ID())
			}

			credential := credentials.ManagedIdentityCredentialResource{
				Type: utils.String(IDENTITY_TYPE),
				Properties: credentials.ManagedIdentityCredential{
					TypeProperties: &credentials.ManagedIdentityTypeProperties{
						ResourceId: &data.IdentityId,
					},
				},
			}

			if len(data.Annotations) > 0 {
				annotations := make([]interface{}, len(data.Annotations))
				for i, v := range data.Annotations {
					annotations[i] = v
				}
				credential.Properties.Annotations = &annotations
			}

			if data.Description != "" {
				credential.Properties.Description = &data.Description
			}

			_, err = client.CredentialOperationsCreateOrUpdate(ctx, id, credential, credentials.CredentialOperationsCreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataFactoryCredentialUserAssignedManagedIdentityResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Credentials
			id, err := credentials.ParseCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var data DataFactoryCredentialUserAssignedManagedIdentityResourceSchema
			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.CredentialOperationsGet(ctx, *id, credentials.CredentialOperationsGetOperationOptions{})
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}

			if existing.Model == nil {
				return fmt.Errorf("model was nil for %s", id)
			}

			credential := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				credential.Properties.Description = &data.Description
			}

			if metadata.ResourceData.HasChange("annotations") {
				if len(data.Annotations) > 0 {
					annotations := make([]interface{}, len(data.Annotations))
					for i, v := range data.Annotations {
						annotations[i] = v
					}
					credential.Properties.Annotations = &annotations
				} else {
					credential.Properties.Annotations = nil
				}
			}

			_, err = client.CredentialOperationsCreateOrUpdate(ctx, *id, credential, credentials.CredentialOperationsCreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (DataFactoryCredentialUserAssignedManagedIdentityResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			d := metadata.ResourceData
			client := metadata.Client.DataFactory.Credentials

			credentialId, err := credentials.ParseCredentialID(d.Id())
			if err != nil {
				return err
			}

			_, err = client.CredentialOperationsDelete(ctx, *credentialId)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
