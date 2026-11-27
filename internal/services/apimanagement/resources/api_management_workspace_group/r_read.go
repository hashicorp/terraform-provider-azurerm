// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GroupClient_v2024_05_01

			id, err := group.ParseWorkspaceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceGroupGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceGroupModel{
				Name:                     id.GroupId,
				ApiManagementWorkspaceId: group.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.Description = pointer.From(props.Description)
					state.ExternalId = pointer.From(props.ExternalId)
					state.Type = pointer.FromEnum(props.Type)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
