// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/backupsautomaticandondemand"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = PostgresqlFlexibleServerBackupResource{}

type PostgresqlFlexibleServerBackupResource struct{}

func (r PostgresqlFlexibleServerBackupResource) ModelObject() interface{} {
	return &PostgresqlFlexibleServerBackupResourceModel{}
}

type PostgresqlFlexibleServerBackupResourceModel struct {
	Name          string `tfschema:"name"`
	ServerId      string `tfschema:"server_id"`
	CompletedTime string `tfschema:"completed_time"`
}

func (r PostgresqlFlexibleServerBackupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupsautomaticandondemand.ValidateBackupID
}

func (r PostgresqlFlexibleServerBackupResource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_backup"
}

func (r PostgresqlFlexibleServerBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerBackupName,
		},

		"server_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupsautomaticandondemand.FlexibleServerId{}),
	}
}

func (r PostgresqlFlexibleServerBackupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"completed_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r PostgresqlFlexibleServerBackupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.BackupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model PostgresqlFlexibleServerBackupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := backupsautomaticandondemand.ParseFlexibleServerID(model.ServerId)
			if err != nil {
				return err
			}

			id := backupsautomaticandondemand.NewBackupID(subscriptionId, serverId.ResourceGroupName, serverId.FlexibleServerName, model.Name)

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if err := client.CreateThenPoll(ctx, id); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PostgresqlFlexibleServerBackupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.BackupsClient

			id, err := backupsautomaticandondemand.ParseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PostgresqlFlexibleServerBackupResourceModel{
				Name:     id.BackupName,
				ServerId: backupsautomaticandondemand.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.CompletedTime = pointer.From(props.CompletedTime)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PostgresqlFlexibleServerBackupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.BackupsClient

			id, err := backupsautomaticandondemand.ParseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
