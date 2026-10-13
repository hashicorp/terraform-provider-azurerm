// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
