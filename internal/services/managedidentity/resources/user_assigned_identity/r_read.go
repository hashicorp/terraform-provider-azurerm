// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package user_assigned_identity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20241130.Identities

			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.UserAssignedIdentitiesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

type UserAssignedIdentityResourceSchema struct {
	ClientId          string         `tfschema:"client_id"`
	IsolationScope    string         `tfschema:"isolation_scope"`
	Location          string         `tfschema:"location"`
	Name              string         `tfschema:"name"`
	PrincipalId       string         `tfschema:"principal_id"`
	ResourceGroupName string         `tfschema:"resource_group_name"`
	Tags              map[string]any `tfschema:"tags"`
	TenantId          string         `tfschema:"tenant_id"`
}
