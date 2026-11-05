// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ScopedRoleDefinitionsClient

			stateId, err := parse.RoleDefinitionId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			id := roledefinitions.NewScopedRoleDefinitionID(stateId.Scope, stateId.RoleID)

			resp, err := client.Delete(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("deleting %s: %+v", stateId, err)
			}

			// Deletes are not instant and can take time to propagate
			poller := custompollers.NewEventualConsistencyPoller(20, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := client.Get(pollerCtx, id)
				return resp.HttpResponse, err
			}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", stateId, err)
			}

			return nil
		},
	}
}
