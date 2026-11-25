// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_dns_alias

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverdnsaliases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.ServerDNSAliasClient

			var alias ServerDNSAliasModel
			if err := metadata.Decode(&alias); err != nil {
				return err
			}

			// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
			// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
			serverID, err := commonids.ParseSqlServerIDInsensitively(alias.MsSQLServerId)
			if err != nil {
				return err
			}

			id := serverdnsaliases.NewDnsAliasID(serverID.SubscriptionId, serverID.ResourceGroupName, serverID.ServerName, alias.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
