// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CognitiveAccountConnectionCustomKeysResource{}

type CognitiveAccountConnectionCustomKeysResource struct{}

func (r CognitiveAccountConnectionCustomKeysResource) ResourceType() string {
	return "azurerm_cognitive_account_connection_custom_keys"
}

func (r CognitiveAccountConnectionCustomKeysResource) ModelObject() interface{} {
	return &CognitiveAccountConnectionCustomKeysModel{}
}

func (r CognitiveAccountConnectionCustomKeysResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

type CognitiveAccountConnectionCustomKeysModel struct {
	Category           string            `tfschema:"category"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	CustomKeys         map[string]string `tfschema:"custom_keys"`
	Metadata           map[string]string `tfschema:"metadata"`
	Name               string            `tfschema:"name"`
	Target             string            `tfschema:"target"`
}

func (r CognitiveAccountConnectionCustomKeysResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}),

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				accountconnectionresource.PossibleValuesForConnectionCategory(),
				false,
			),
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"target": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"custom_keys": {
			Type:      pluginsdk.TypeMap,
			Required:  true,
			Sensitive: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r CognitiveAccountConnectionCustomKeysResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionCustomKeysResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			var model CognitiveAccountConnectionCustomKeysModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := accountconnectionresource.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			id := accountconnectionresource.NewConnectionID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.AccountConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := accountconnectionresource.CustomKeysConnectionProperties{
				AuthType: accountconnectionresource.ConnectionAuthTypeCustomKeys,
				Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
				Metadata: pointer.To(model.Metadata),
				Target:   pointer.To(model.Target),
				Credentials: &accountconnectionresource.CustomKeys{
					Keys: &model.CustomKeys,
				},
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.AccountConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveAccountConnectionCustomKeysResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var currentState CognitiveAccountConnectionCustomKeysModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountConnectionCustomKeysModel{
				CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
				Name:               id.ConnectionName,
				CustomKeys:         currentState.CustomKeys,
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				base := props.ConnectionPropertiesV2()
				state.Category = pointer.FromEnum(base.Category)
				state.Target = pointer.From(base.Target)
				state.Metadata = map[string]string{}

				// Only include metadata fields that were in the original config.
				// The API returns additional metadata fields beyond what was configured (e.g., `ApiVersion`,
				// `DeploymentApiVersion`), which would cause unwanted diffs.
				if len(currentState.Metadata) > 0 {
					apiMetadata := pointer.From(base.Metadata)

					for configKey := range currentState.Metadata {
						for apiKey, apiValue := range apiMetadata {
							if strings.EqualFold(configKey, apiKey) {
								state.Metadata[configKey] = apiValue
								break
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountConnectionCustomKeysResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountConnectionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var model CognitiveAccountConnectionCustomKeysModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props, ok := resp.Model.Properties.(accountconnectionresource.CustomKeysConnectionProperties)
			if !ok {
				return fmt.Errorf("unexpected properties type for %s", *id)
			}

			props.Credentials = &accountconnectionresource.CustomKeys{
				Keys: &model.CustomKeys,
			}

			if metadata.ResourceData.HasChange("target") {
				props.Target = pointer.To(model.Target)
			}

			if metadata.ResourceData.HasChange("metadata") {
				props.Metadata = pointer.To(model.Metadata)
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: props,
			}

			if _, err := client.AccountConnectionsCreate(ctx, *id, connection); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionCustomKeysResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.AccountConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
