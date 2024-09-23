// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverdnsaliases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServerDNSAliasModel struct {
	MsSQLServerId string `tfschema:"mssql_server_id"`
	Name          string `tfschema:"name"`
	DNSRecord     string `tfschema:"dns_record"`
}

type ServerDNSAliasResource struct{}

var _ sdk.Resource = (*ServerDNSAliasResource)(nil)

func (m ServerDNSAliasResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"mssql_server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDNSAliasName,
		},
	}
}

func (m ServerDNSAliasResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_record": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m ServerDNSAliasResource) ModelObject() interface{} {
	return &ServerDNSAliasModel{}
}

func (m ServerDNSAliasResource) ResourceType() string {
	return "azurerm_mssql_server_dns_alias"
}

func (m ServerDNSAliasResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.ServerDNSAliasClient

			var alias ServerDNSAliasModel
			if err := metadata.Decode(&alias); err != nil {
				return err
			}

			serverID, err := parse.ServerID(alias.MsSQLServerId)
			if err != nil {
				return err
			}

			id := serverdnsaliases.NewDnsAliasID(serverID.SubscriptionId, serverID.ResourceGroup, serverID.Name, alias.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return metadata.ResourceRequiresImport(m.ResourceType(), id)
			}

			err = client.CreateOrUpdateThenPoll(ctx, id)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (m ServerDNSAliasResource) Read() sdk.ResourceFunc {
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
				MsSQLServerId: parse.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName).ID(),
			}
			if alias.Model != nil {
				if prop := alias.Model.Properties; prop != nil {
					state.DNSRecord = utils.NormalizeNilableString(prop.AzureDnsRecord)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (m ServerDNSAliasResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := serverdnsaliases.ParseDnsAliasID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			metadata.Logger.Infof("deleting %s", id)
			client := metadata.Client.MSSQL.ServerDNSAliasClient
			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m ServerDNSAliasResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ServerDNSAliasID
}
