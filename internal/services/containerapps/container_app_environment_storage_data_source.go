// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentStorageDataSource struct{}

type ContainerAppEnvironmentStorageDataSourceModel struct {
	Name                      string `tfschema:"name"`
	ContainerAppEnvironmentId string `tfschema:"container_app_environment_id"`

	AccountName  string `tfschema:"account_name"`
	ShareName    string `tfschema:"share_name"`
	AccessMode   string `tfschema:"access_mode"`
	NfsServerUrl string `tfschema:"nfs_server_url"`
}

var _ sdk.DataSource = ContainerAppEnvironmentStorageDataSource{}

func (r ContainerAppEnvironmentStorageDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentStorageDataSourceModel{}
}

func (r ContainerAppEnvironmentStorageDataSource) ResourceType() string {
	return "azurerm_container_app_environment_storage"
}

func (r ContainerAppEnvironmentStorageDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentStorageName,
			Description:  "The name for this Container App Environment Storage.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironmentsstorages.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to which this storage belongs.",
		},
	}
}

func (r ContainerAppEnvironmentStorageDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Azure Storage Account in which the Share is located.",
		},

		"share_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Azure Storage Share.",
		},

		"access_mode": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The access mode to connect this storage to the Container App.",
		},

		"nfs_server_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The NFS server URL for the Azure File Share.",
		},
	}
}

func (r ContainerAppEnvironmentStorageDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			var state ContainerAppEnvironmentStorageDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			envId, err := managedenvironmentsstorages.ParseManagedEnvironmentID(state.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := managedenvironmentsstorages.NewStorageID(envId.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, state.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state.Name = id.StorageName
			state.ContainerAppEnvironmentId = envId.ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					if azureFile := props.AzureFile; azureFile != nil {
						state.AccountName = pointer.From(azureFile.AccountName)
						if azureFile.AccessMode != nil {
							state.AccessMode = string(*azureFile.AccessMode)
						}
						state.ShareName = pointer.From(azureFile.ShareName)
					} else if nfsAzureFile := props.NfsAzureFile; nfsAzureFile != nil {
						state.NfsServerUrl = pointer.From(nfsAzureFile.Server)
						if nfsAzureFile.AccessMode != nil {
							state.AccessMode = string(*nfsAzureFile.AccessMode)
						}
						state.ShareName = pointer.From(nfsAzureFile.ShareName)
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
