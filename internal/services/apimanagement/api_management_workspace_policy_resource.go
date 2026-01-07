// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name api_management_workspace_policy -service-package-name apimanagement -compare-values "resource_group_name:api_management_workspace_id,service_name:api_management_workspace_id,workspace_id:api_management_workspace_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basic

type ApiManagementWorkspacePolicyModel struct {
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	XmlContent               string `tfschema:"xml_content"`
	XmlLink                  string `tfschema:"xml_link"`
}

type ApiManagementWorkspacePolicyResource struct{}

var (
	_ sdk.ResourceWithUpdate   = ApiManagementWorkspacePolicyResource{}
	_ sdk.ResourceWithIdentity = ApiManagementWorkspacePolicyResource{}
)

func (r ApiManagementWorkspacePolicyResource) ResourceType() string {
	return "azurerm_api_management_workspace_policy"
}

func (r ApiManagementWorkspacePolicyResource) Identity() resourceids.ResourceId {
	return &workspacepolicy.WorkspaceId{}
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
			Computed:         true,
			ExactlyOneOf:     []string{"xml_link", "xml_content"},
			DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
		},

		"xml_link": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"xml_link", "xml_content"},
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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

			workspaceId, err := workspacepolicy.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *workspaceId, workspacepolicy.GetOperationOptions{})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", *workspaceId, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), *workspaceId)
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

			if _, err := client.CreateOrUpdate(ctx, *workspaceId, parameters, workspacepolicy.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", *workspaceId, err)
			}

			metadata.SetID(*workspaceId)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, workspaceId); err != nil {
				return err
			}

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

			resp, err := client.Get(ctx, *workspaceId, workspacepolicy.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *workspaceId, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *workspaceId)
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

			if _, err := client.CreateOrUpdate(ctx, *workspaceId, *payload, workspacepolicy.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *workspaceId, err)
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
					state.XmlLink = metadata.ResourceData.Get("xml_link").(string)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
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
