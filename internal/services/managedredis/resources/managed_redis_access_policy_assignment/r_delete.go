// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_redis_access_policy_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	managedredisClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := managedredisClient.From(metadata.Client).DatabaseClient
			id, err := databases.ParseAccessPolicyAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.AccessPolicyAssignmentDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
