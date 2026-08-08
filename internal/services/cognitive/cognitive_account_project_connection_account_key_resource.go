// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_project_connection_account_key -properties "name" -compare-values "subscription_id:cognitive_account_project_id,resource_group_name:cognitive_account_project_id,account_name:cognitive_account_project_id,project_name:cognitive_account_project_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate         = CognitiveAccountProjectConnectionAccountKeyResource{}
	_ sdk.ResourceWithIdentity       = CognitiveAccountProjectConnectionAccountKeyResource{}
	_ sdk.ResourceWithCustomImporter = CognitiveAccountProjectConnectionAccountKeyResource{}
	_ sdk.ResourceWithCustomizeDiff  = CognitiveAccountProjectConnectionAccountKeyResource{}
)

type CognitiveAccountProjectConnectionAccountKeyResource struct{}

func (r CognitiveAccountProjectConnectionAccountKeyResource) CustomImporter() sdk.ResourceRunFunc {
	return cognitiveAccountProjectConnectionImporter(projectconnectionresource.ConnectionAuthTypeAccountKey, r.ResourceType())
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Identity() resourceids.ResourceId {
	return new(projectconnectionresource.ProjectConnectionId)
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) ResourceType() string {
	return "azurerm_cognitive_account_project_connection_account_key"
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) ModelObject() interface{} {
	return &CognitiveAccountProjectConnectionAccountKeyModel{}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projectconnectionresource.ValidateProjectConnectionID
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(_ context.Context, metadata sdk.ResourceMetaData) error {
			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if !rawConfig.IsKnown() || rawConfig.IsNull() {
				return nil
			}

			metadataValue := rawConfig.GetAttr("metadata")
			if !metadataValue.IsKnown() || metadataValue.IsNull() {
				return nil
			}

			for key, value := range metadataValue.AsValueMap() {
				if !strings.EqualFold(key, "ResourceId") {
					continue
				}

				if !value.IsKnown() {
					return nil
				}

				if !value.IsNull() && value.AsString() != "" {
					return nil
				}

				break
			}

			return errors.New("`metadata` must include a non-empty `ResourceId` value")
		},
	}
}

type CognitiveAccountProjectConnectionAccountKeyModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AccountKey                string            `tfschema:"account_key"`
	Category                  string            `tfschema:"category"`
	Metadata                  map[string]string `tfschema:"metadata"`
	Target                    string            `tfschema:"target"`
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&cognitiveservicesprojects.ProjectId{}),

		"account_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"category": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{string(projectconnectionresource.ConnectionCategoryAzureStorageAccount)}, false),
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
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			var model CognitiveAccountProjectConnectionAccountKeyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			projectId, err := cognitiveservicesprojects.ParseProjectID(model.CognitiveAccountProjectId)
			if err != nil {
				return err
			}

			id := projectconnectionresource.NewProjectConnectionID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.AccountName, projectId.ProjectName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ProjectConnectionsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := projectconnectionresource.AccountKeyAuthTypeConnectionProperties{
				AuthType: projectconnectionresource.ConnectionAuthTypeAccountKey,
				Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
				Target:   pointer.To(model.Target),
				Credentials: &projectconnectionresource.ConnectionAccountKey{
					Key: pointer.To(model.AccountKey),
				},
			}
			if len(model.Metadata) > 0 {
				properties.Metadata = pointer.To(model.Metadata)
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.ProjectConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var currentState CognitiveAccountProjectConnectionAccountKeyModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			return r.flatten(metadata, id, resp.Model, currentState.Metadata, currentState.AccountKey)
		},
	}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectConnectionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var model CognitiveAccountProjectConnectionAccountKeyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props, ok := resp.Model.Properties.(projectconnectionresource.AccountKeyAuthTypeConnectionProperties)
			if !ok {
				return fmt.Errorf("unexpected properties type for %s", *id)
			}

			if metadata.ResourceData.HasChange("account_key") {
				props.Credentials = &projectconnectionresource.ConnectionAccountKey{
					Key: pointer.To(model.AccountKey),
				}
			}

			if metadata.ResourceData.HasChange("target") {
				props.Target = pointer.To(model.Target)
			}

			if metadata.ResourceData.HasChange("metadata") {
				props.Metadata = pointer.To(model.Metadata)
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: props,
			}

			if _, err := client.ProjectConnectionsCreate(ctx, *id, connection); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectConnectionAccountKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ProjectConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (CognitiveAccountProjectConnectionAccountKeyResource) flatten(metadata sdk.ResourceMetaData, id *projectconnectionresource.ProjectConnectionId, model *projectconnectionresource.ConnectionPropertiesV2BasicResource, priorMetadata map[string]string, priorAccountKey string) error {
	state := CognitiveAccountProjectConnectionAccountKeyModel{
		CognitiveAccountProjectId: cognitiveservicesprojects.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName).ID(),
		Name:                      id.ConnectionName,
		AccountKey:                priorAccountKey,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if model != nil && model.Properties != nil {
		base := model.Properties.ConnectionPropertiesV2()
		state.Category = pointer.FromEnum(base.Category)
		state.Target = pointer.From(base.Target)
		state.Metadata = flattenAccountConnectionMetadata(priorMetadata, base.Metadata)
	}

	return metadata.Encode(&state)
}
