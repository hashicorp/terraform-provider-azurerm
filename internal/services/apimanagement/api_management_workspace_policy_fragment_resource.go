// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name api_management_workspace_policy_fragment -service-package-name apimanagement -properties "name" -compare-values "resource_group_name:api_management_workspace_id,service_name:api_management_workspace_id,workspace_id:api_management_workspace_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basic

type ApiManagementWorkspacePolicyFragmentModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	Description              string `tfschema:"description"`
	XmlFormat                string `tfschema:"xml_format"`
	XmlContent               string `tfschema:"xml_content"`
}

type ApiManagementWorkspacePolicyFragmentResource struct{}

var (
	_ sdk.ResourceWithUpdate   = ApiManagementWorkspacePolicyFragmentResource{}
	_ sdk.ResourceWithIdentity = ApiManagementWorkspacePolicyFragmentResource{}
)

func (r ApiManagementWorkspacePolicyFragmentResource) ResourceType() string {
	return "azurerm_api_management_workspace_policy_fragment"
}

func (r ApiManagementWorkspacePolicyFragmentResource) Identity() resourceids.ResourceId {
	return &policyfragment.WorkspacePolicyFragmentId{}
}

func (r ApiManagementWorkspacePolicyFragmentResource) ModelObject() interface{} {
	return &ApiManagementWorkspacePolicyFragmentModel{}
}

func (r ApiManagementWorkspacePolicyFragmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policyfragment.ValidateWorkspacePolicyFragmentID
}

func (r ApiManagementWorkspacePolicyFragmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&policyfragment.WorkspaceId{}),

		"xml_content": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: XmlWhitespaceDiffSuppress,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"xml_format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(policyfragment.PolicyFragmentContentFormatXml),
			ValidateFunc: validation.StringInSlice(policyfragment.PossibleValuesForPolicyFragmentContentFormat(), false),
		},
	}
}

func (r ApiManagementWorkspacePolicyFragmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspacePolicyFragmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.PolicyFragmentClient_v2024_05_01

			var model ApiManagementWorkspacePolicyFragmentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := policyfragment.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := policyfragment.NewWorkspacePolicyFragmentID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)
			existing, err := client.WorkspacePolicyFragmentGet(ctx, id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormat(model.XmlFormat)),
			})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := policyfragment.PolicyFragmentContract{
				Properties: &policyfragment.PolicyFragmentContractProperties{
					Format: pointer.To(policyfragment.PolicyFragmentContentFormat(model.XmlFormat)),
					Value:  model.XmlContent,
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if err := client.WorkspacePolicyFragmentCreateOrUpdateThenPoll(ctx, id, parameters, policyfragment.WorkspacePolicyFragmentCreateOrUpdateOperationOptions{}); err != nil {
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

func (r ApiManagementWorkspacePolicyFragmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.PolicyFragmentClient_v2024_05_01

			var model ApiManagementWorkspacePolicyFragmentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := policyfragment.ParseWorkspacePolicyFragmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspacePolicyFragmentGet(ctx, *id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormat(model.XmlFormat)),
			})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := resp.Model
			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("xml_content") {
				payload.Properties.Value = model.XmlContent
			}

			if metadata.ResourceData.HasChange("xml_format") {
				payload.Properties.Format = pointer.To(policyfragment.PolicyFragmentContentFormat(model.XmlFormat))
			}

			if err := client.WorkspacePolicyFragmentCreateOrUpdateThenPoll(ctx, *id, *payload, policyfragment.WorkspacePolicyFragmentCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspacePolicyFragmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.PolicyFragmentClient_v2024_05_01

			id, err := policyfragment.ParseWorkspacePolicyFragmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			format := metadata.ResourceData.Get("xml_format").(string)
			if format == "" {
				format = string(policyfragment.PolicyFragmentContentFormatXml)
			}
			resp, err := client.WorkspacePolicyFragmentGet(ctx, *id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{
				Format: pointer.To(policyfragment.PolicyFragmentContentFormat(format)),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspacePolicyFragmentModel{
				Name:                     id.PolicyFragmentName,
				ApiManagementWorkspaceId: policyfragment.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.XmlContent = props.Value

					// The API only returns `xml_format` when set to "rawxml"; the default "xml" is intentionally never returned.
					format := policyfragment.PolicyFragmentContentFormatXml
					if props.Format != nil {
						format = pointer.From(props.Format)
					}
					state.XmlFormat = string(format)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspacePolicyFragmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.PolicyFragmentClient_v2024_05_01

			id, err := policyfragment.ParseWorkspacePolicyFragmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspacePolicyFragmentDelete(ctx, *id, policyfragment.WorkspacePolicyFragmentDeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
