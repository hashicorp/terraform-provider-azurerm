// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy_fragment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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
				Format: pointer.ToEnum[policyfragment.PolicyFragmentContentFormat](model.XmlFormat),
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
				payload.Properties.Format = pointer.ToEnum[policyfragment.PolicyFragmentContentFormat](model.XmlFormat)
			}

			if err := client.WorkspacePolicyFragmentCreateOrUpdateThenPoll(ctx, *id, *payload, policyfragment.WorkspacePolicyFragmentCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
