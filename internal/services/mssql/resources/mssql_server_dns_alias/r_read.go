// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_dns_alias

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverdnsaliases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := serverdnsaliases.ParseDnsAliasID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.MSSQL.ServerDNSAliasClient
			alias, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(alias.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return err
			}
			state := ServerDNSAliasModel{
				Name:          id.DnsAliasName,
				MsSQLServerId: commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName).ID(),
			}
			if alias.Model != nil {
				if prop := alias.Model.Properties; prop != nil {
					state.DNSRecord = pointer.From(prop.AzureDnsRecord)
				}
			}
			return metadata.Encode(&state)
		},
	}
}
