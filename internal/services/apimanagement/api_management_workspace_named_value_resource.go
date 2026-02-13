package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceNamedValueResource struct{}

var (
	_ sdk.ResourceWithUpdate        = ApiManagementWorkspaceNamedValueResource{}
	_ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceNamedValueResource{}
)

func (r ApiManagementWorkspaceNamedValueResource) ResourceType() string {
	return "azurerm_api_management_workspace_named_value"
}

func (r ApiManagementWorkspaceNamedValueResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceNamedValueModel{}
}

func (r ApiManagementWorkspaceNamedValueResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return namedvalue.ValidateWorkspaceNamedValueID
}

type ApiManagementWorkspaceNamedValueModel struct {
	ApiManagementWorkspaceId string     `tfschema:"api_management_workspace_id"`
	DisplayName              string     `tfschema:"display_name"`
	Name                     string     `tfschema:"name"`
	Secret                   bool       `tfschema:"secret"`
	Tags                     []string   `tfschema:"tags"`
	Value                    string     `tfschema:"value"`
	ValueFromKeyVault        []KeyVault `tfschema:"value_from_key_vault"`
}

type KeyVault struct {
	IdentityClientId string `tfschema:"identity_client_id"`
	SecretId         string `tfschema:"secret_id"`
}

func (r ApiManagementWorkspaceNamedValueResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z](?:[a-zA-Z0-9-]{0,254}[a-zA-Z0-9])?$`),
				"`name` must be 1-256 characters, starting with a letter, and using only letters, numbers, or hyphens"),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z0-9-._]{1,256}$`),
				"`display_name` must be 1-256 characters and can only contain letters, numbers, hyphens, periods, and underscores"),
		},

		"secret": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// NOTE: this is not the common `tags` attribute used by most Azure resources, this field accepts a list of strings rather than a map of key-value pairs
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
					"secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
					},
					"identity_client_id": {
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

func (r ApiManagementWorkspaceNamedValueResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.ValueFromKeyVault) > 0 && !model.Secret {
				return errors.New("`secret` must be set to `true` when `value_from_key_vault` is specified")
			}

			return nil
		},
	}
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
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: model.DisplayName,
					KeyVault:    expandApiManagementWorkspaceNamedValueKeyVault(model.ValueFromKeyVault),
					Secret:      pointer.To(model.Secret),
					Tags:        pointer.To(model.Tags),
				},
			}

			if model.Value != "" {
				parameters.Properties.Value = pointer.To(model.Value)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, id, parameters, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

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

			var config ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceNamedValueModel{
				Name:                     id.NamedValueId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					model.DisplayName = props.DisplayName
					model.Secret = pointer.From(props.Secret)
					model.Tags = pointer.From(props.Tags)
					if model.Secret {
						// `value` is not retrievable when `secret` is `true`, so we use the config value
						model.Value = config.Value
					} else {
						model.Value = pointer.From(props.Value)
					}
					model.ValueFromKeyVault = flattenApiManagementWorkspaceNamedValueKeyVault(props.KeyVault)
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ApiManagementWorkspaceNamedValueResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: existing.Model.Properties.DisplayName,
					Secret:      existing.Model.Properties.Secret,
					Tags:        existing.Model.Properties.Tags,
					Value:       existing.Model.Properties.Value,
				},
			}

			if existing.Model.Properties.KeyVault != nil {
				parameters.Properties.KeyVault = &namedvalue.KeyVaultContractCreateProperties{
					SecretIdentifier: existing.Model.Properties.KeyVault.SecretIdentifier,
					IdentityClientId: existing.Model.Properties.KeyVault.IdentityClientId,
				}
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("secret") {
				parameters.Properties.Secret = pointer.To(model.Secret)
			}

			if metadata.ResourceData.HasChange("value") {
				if model.Value != "" {
					parameters.Properties.Value = pointer.To(model.Value)
				} else {
					parameters.Properties.Value = nil
				}
			}

			if metadata.ResourceData.HasChange("value_from_key_vault") {
				parameters.Properties.KeyVault = expandApiManagementWorkspaceNamedValueKeyVault(model.ValueFromKeyVault)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Properties.Tags = pointer.To(model.Tags)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, *id, parameters, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
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

func expandApiManagementWorkspaceNamedValueKeyVault(inputs []KeyVault) *namedvalue.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0]

	result := namedvalue.KeyVaultContractCreateProperties{
		SecretIdentifier: pointer.To(input.SecretId),
	}

	if input.IdentityClientId != "" {
		result.IdentityClientId = pointer.To(input.IdentityClientId)
	}

	return &result
}

func flattenApiManagementWorkspaceNamedValueKeyVault(input *namedvalue.KeyVaultContractProperties) []KeyVault {
	if input == nil {
		return []KeyVault{}
	}

	return []KeyVault{
		{
			SecretId:         pointer.From(input.SecretIdentifier),
			IdentityClientId: pointer.From(input.IdentityClientId),
		},
	}
}
