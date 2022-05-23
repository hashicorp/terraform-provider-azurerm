package mssql

import (
	"context"
	"fmt"
	"time"

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

			id := parse.NewServerDNSAliasID(serverID.SubscriptionId, serverID.ResourceGroup, serverID.Name, alias.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.DnsAliaseName)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return metadata.ResourceRequiresImport(m.ResourceType(), id)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.DnsAliaseName)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %v", id, err)
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
			id, err := parse.ServerDNSAliasID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.MSSQL.ServerDNSAliasClient
			alias, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.DnsAliaseName)
			if err != nil {
				return err
			}
			state := ServerDNSAliasModel{
				Name:          id.DnsAliaseName,
				MsSQLServerId: parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName).ID(),
			}
			if prop := alias.ServerDNSAliasProperties; prop != nil {
				state.DNSRecord = utils.NormalizeNilableString(prop.AzureDNSRecord)
			}
			return metadata.Encode(&state)
		},
	}
}

func (m ServerDNSAliasResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ServerDNSAliasID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			metadata.Logger.Infof("deleting %s", id)
			client := metadata.Client.MSSQL.ServerDNSAliasClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.ServerName, id.DnsAliaseName); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m ServerDNSAliasResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ServerDNSAliasID
}
