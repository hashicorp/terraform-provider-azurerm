// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceNamedValueModel struct {
	Name                     string                             `tfschema:"name"`
	ApiManagementWorkspaceId string                             `tfschema:"api_management_workspace_id"`
	DisplayName              string                             `tfschema:"display_name"`
	Value                    string                             `tfschema:"value"`
	ValueFromKeyVault        []WorkspaceNamedValueKeyVaultModel `tfschema:"value_from_key_vault"`
	SecretEnabled            bool                               `tfschema:"secret_enabled"`
	Tags                     []string                           `tfschema:"tags"`
}

type WorkspaceNamedValueKeyVaultModel struct {
	KeyVaultSecretId             string `tfschema:"key_vault_secret_id"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`
}

type ApiManagementWorkspaceNamedValueResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceNamedValueResource{}

var _ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceNamedValueResource{}

func (r ApiManagementWorkspaceNamedValueResource) ResourceType() string {
	return "azurerm_api_management_workspace_named_value"
}

func (r ApiManagementWorkspaceNamedValueResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceNamedValueModel{}
}

func (r ApiManagementWorkspaceNamedValueResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namedvalue.ValidateWorkspaceNamedValueID
}

func (r ApiManagementWorkspaceNamedValueResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementNamedValueDisplayName,
		},

		"secret_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"value_from_key_vault": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"value", "value_from_key_vault"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					},
					"user_assigned_identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceNamedValueResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := namedvalue.NewWorkspaceNamedValueID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			existing, err := client.WorkspaceNamedValueGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: model.DisplayName,
					Secret:      pointer.To(model.SecretEnabled),
				},
			}

			if len(model.ValueFromKeyVault) > 0 {
				payload.Properties.KeyVault = expandWorkspaceNamedValueKeyVault(model.ValueFromKeyVault)
			}

			if model.Value != "" {
				payload.Properties.Value = pointer.To(model.Value)
			}

			if len(model.Tags) > 0 {
				payload.Properties.Tags = pointer.To(model.Tags)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, id, payload, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: resp.Model.Properties.DisplayName,
					Secret:      resp.Model.Properties.Secret,
					Tags:        resp.Model.Properties.Tags,
					Value:       resp.Model.Properties.Value,
				},
			}

			if keyVault := resp.Model.Properties.KeyVault; keyVault != nil {
				parameters.Properties.KeyVault = &namedvalue.KeyVaultContractCreateProperties{
					IdentityClientId: keyVault.IdentityClientId,
					SecretIdentifier: keyVault.SecretIdentifier,
				}
			}
			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("secret_enabled") {
				parameters.Properties.Secret = pointer.To(model.SecretEnabled)
			}

			if metadata.ResourceData.HasChange("value") {
				parameters.Properties.Value = pointer.To(model.Value)
			}

			if metadata.ResourceData.HasChange("value_from_key_vault") {
				parameters.Properties.KeyVault = expandWorkspaceNamedValueKeyVault(model.ValueFromKeyVault)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Properties.Tags = pointer.To(model.Tags)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, *id, parameters, namedvalue.WorkspaceNamedValueCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceNamedValueModel{
				Name:                     id.NamedValueId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.SecretEnabled = pointer.From(props.Secret)
					// The `value` is retrieved from state since, when secret_enabled = true, the API does not return it due to its sensitivity.
					value := pointer.From(props.Value)
					if state.SecretEnabled {
						value = metadata.ResourceData.Get("value").(string)
					}
					state.Value = value
					state.ValueFromKeyVault = flattenWorkspaceNamedValueKeyVault(props.KeyVault)
					state.Tags = pointer.From(props.Tags)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceNamedValueDelete(ctx, *id, namedvalue.DefaultWorkspaceNamedValueDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if _, ok := rd.GetOk("value_from_key_vault"); ok {
				if !rd.Get("secret_enabled").(bool) {
					return errors.New("`secret_enabled` must be set to `true` when `value_from_key_vault` is specified")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandWorkspaceNamedValueKeyVault(inputs []WorkspaceNamedValueKeyVaultModel) *namedvalue.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := &inputs[0]
	result := namedvalue.KeyVaultContractCreateProperties{
		SecretIdentifier: pointer.To(input.KeyVaultSecretId),
	}

	if input.UserAssignedIdentityClientId != "" {
		result.IdentityClientId = pointer.To(input.UserAssignedIdentityClientId)
	}

	return &result
}

func flattenWorkspaceNamedValueKeyVault(input *namedvalue.KeyVaultContractProperties) []WorkspaceNamedValueKeyVaultModel {
	outputList := make([]WorkspaceNamedValueKeyVaultModel, 0)
	if input == nil {
		return outputList
	}

	output := WorkspaceNamedValueKeyVaultModel{
		KeyVaultSecretId:             pointer.From(input.SecretIdentifier),
		UserAssignedIdentityClientId: pointer.From(input.IdentityClientId),
	}

	return append(outputList, output)
}
