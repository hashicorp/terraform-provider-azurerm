// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_account_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/helpers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			projectId := projectconnectionresource.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName)
			locks.ByID(projectId.ID())
			defer locks.UnlockByID(projectId.ID())

			resp, err := client.ProjectConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if !helpers.CognitiveAccountProjectConnectionHasExpectedAuthType(resp.Model, projectconnectionresource.ConnectionAuthTypeAccountKey) {
				return fmt.Errorf("connection %s does not have auth type `AccountKey` and cannot be deleted by `%s`", *id, r.ResourceType())
			}

			if resp, err := client.ProjectConnectionsDelete(ctx, *id); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
