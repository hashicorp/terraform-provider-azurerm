// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspacePolicyModel struct {
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	XmlContent               string `tfschema:"xml_content"`
	XmlLink                  string `tfschema:"xml_link"`
}

type ApiManagementWorkspacePolicyResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspacePolicyResource{}

func (r ApiManagementWorkspacePolicyResource) ResourceType() string {
	return "azurerm_api_management_workspace_policy"
}

func (r ApiManagementWorkspacePolicyResource) ModelObject() interface{} {
	return &ApiManagementWorkspacePolicyModel{}
}

func (r ApiManagementWorkspacePolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspacepolicy.ValidateWorkspaceID
}

func (r ApiManagementWorkspacePolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"xml_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C when `xml_link` is provided, the API downloads it into `xml_content`.
			// O+C ensures Terraform state is populated correctly and prevents spurious diffs.
			Computed:         true,
			ExactlyOneOf:     []string{"xml_link", "xml_content"},
			DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"xml_link", "xml_content"},
			ValidateFunc: validation.StringIsNotEmpty,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				// Suppress spurious diffs during plan or import because the API does not return `xml_link` by design.
				// This prevents Terraform from treating the missing API value as a difference, so the user-specified
				// `xml_link` in the configuration is not marked for change.
				return old == "" && d.Id() != "" && new != ""
			},
		},
	}
}

func (r ApiManagementWorkspacePolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspacePolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspacePolicyClient

			var model ApiManagementWorkspacePolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := workspacepolicy.NewWorkspaceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId)
			existing, err := client.Get(ctx, id, workspacepolicy.GetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := workspacepolicy.PolicyContract{}

			if model.XmlLink != "" {
				parameters.Properties = &workspacepolicy.PolicyContractProperties{
					Format: pointer.To(workspacepolicy.PolicyContentFormatRawxmlNegativelink),
					Value:  model.XmlLink,
				}
			}

			if model.XmlContent != "" {
				parameters.Properties = &workspacepolicy.PolicyContractProperties{
					Format: pointer.To(workspacepolicy.PolicyContentFormatRawxml),
					Value:  model.XmlContent,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters, workspacepolicy.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspacePolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspacePolicyClient

			workspaceId, err := workspacepolicy.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspacePolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := workspacepolicy.NewWorkspaceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId)
			resp, err := client.Get(ctx, id, workspacepolicy.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := resp.Model
			if metadata.ResourceData.HasChange("xml_link") {
				if model.XmlLink != "" {
					payload.Properties = &workspacepolicy.PolicyContractProperties{
						Format: pointer.To(workspacepolicy.PolicyContentFormatRawxmlNegativelink),
						Value:  model.XmlLink,
					}
				}
			}
			if metadata.ResourceData.HasChange("xml_content") {
				if model.XmlContent != "" {
					payload.Properties = &workspacepolicy.PolicyContractProperties{
						Format: pointer.To(workspacepolicy.PolicyContentFormatRawxml),
						Value:  model.XmlContent,
					}
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, *payload, workspacepolicy.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspacePolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspacePolicyClient

			id, err := workspacepolicy.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, workspacepolicy.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspacePolicyModel{
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
					// as such there is no way to set `xml_link` and we'll let Terraform handle it
					state.XmlContent = html.UnescapeString(props.Value)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspacePolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspacePolicyClient

			id, err := workspacepolicy.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id, workspacepolicy.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
