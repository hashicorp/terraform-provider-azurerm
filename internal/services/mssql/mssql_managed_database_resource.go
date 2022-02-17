package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedDatabaseModel struct {
	Name              string `tfschema:"name"`
	ManagedInstanceId string `tfschema:"managed_instance_id"`
}

var _ sdk.Resource = MsSqlManagedDatabaseResource{}

type MsSqlManagedDatabaseResource struct{}

func (r MsSqlManagedDatabaseResource) ResourceType() string {
	return "azurerm_mssql_managed_database"
}

func (r MsSqlManagedDatabaseResource) ModelObject() interface{} {
	return &MsSqlManagedDatabaseModel{}
}

func (r MsSqlManagedDatabaseResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedDatabaseID
}

func (r MsSqlManagedDatabaseResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDatabaseName,
		},

		"managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},
	}
}

func (r MsSqlManagedDatabaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlManagedDatabaseResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.ManagedDatabasesClient
			instancesClient := metadata.Client.MSSQL.ManagedInstancesClient

			var model MsSqlManagedDatabaseModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := parse.ManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			id := parse.NewManagedDatabaseID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroup, managedInstanceId.Name, model.Name)

			managedInstance, err := instancesClient.Get(ctx, managedInstanceId.ResourceGroup, managedInstanceId.Name, "")
			if err != nil || managedInstance.Location == nil || *managedInstance.Location == "" {
				return fmt.Errorf("checking for existence and region of Managed Instance for %s: %+v", id, err)
			}

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := sql.ManagedDatabase{
				Location: managedInstance.Location,
			}

			metadata.Logger.Infof("Creating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MsSqlManagedDatabaseResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.ManagedDatabasesClient

			id, err := parse.ManagedDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedDatabaseModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			managedInstanceId := parse.NewManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			model := MsSqlManagedDatabaseModel{
				Name:              id.DatabaseName,
				ManagedInstanceId: managedInstanceId.ID(),
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedDatabaseResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.ManagedDatabasesClient

			id, err := parse.ManagedDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}
