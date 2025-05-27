// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentStorageResource struct{}

type ContainerAppEnvironmentStorageModel struct {
	Name                      string `tfschema:"name"`
	ContainerAppEnvironmentId string `tfschema:"container_app_environment_id"`
	AccountName               string `tfschema:"account_name"`
	AccessKey                 string `tfschema:"access_key"`
	ShareName                 string `tfschema:"share_name"`
	AccessMode                string `tfschema:"access_mode"`
	NfsServer                 string `tfschema:"nfs_server_url"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentStorageResource{}

func (r ContainerAppEnvironmentStorageResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentStorageModel{}
}

func (r ContainerAppEnvironmentStorageResource) ResourceType() string {
	return "azurerm_container_app_environment_storage"
}

func (r ContainerAppEnvironmentStorageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironmentsstorages.ValidateStorageID
}

func (r ContainerAppEnvironmentStorageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentStorageName,
			Description:  "The name for this Storage.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironmentsstorages.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to which this storage belongs.",
		},

		"account_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  storageValidate.StorageAccountName,
			RequiredWith:  []string{"access_key"},
			ConflictsWith: []string{"nfs_server_url"},
			Description:   "The Azure Storage Account in which the Share to be used is located.",
		},

		"access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"account_name"},
			Description:  "The Storage Account Access Key.",
		},

		"share_name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Azure Storage Share to use.",
		},

		"access_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedenvironmentsstorages.AccessModeReadOnly),
				string(managedenvironmentsstorages.AccessModeReadWrite),
			}, false),
			Description: "The access mode to connect this storage to the Container App. Possible values include `ReadOnly` and `ReadWrite`.",
		},

		"nfs_server_url": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"account_name"},
		},
	}
}

func (r ContainerAppEnvironmentStorageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerAppEnvironmentStorageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			var storage ContainerAppEnvironmentStorageModel

			if err := metadata.Decode(&storage); err != nil {
				return err
			}

			containerAppEnvironmentId, err := managedenvironmentsstorages.ParseManagedEnvironmentID(storage.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := managedenvironmentsstorages.NewStorageID(metadata.Client.Account.SubscriptionId, containerAppEnvironmentId.ResourceGroupName, containerAppEnvironmentId.ManagedEnvironmentName, storage.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			accessMode := managedenvironmentsstorages.AccessMode(storage.AccessMode)
			managedEnvironmentStorage := managedenvironmentsstorages.ManagedEnvironmentStorage{}

			if storage.NfsServer != "" {
				props := &managedenvironmentsstorages.ManagedEnvironmentStorageProperties{
					NfsAzureFile: &managedenvironmentsstorages.NfsAzureFileProperties{
						AccessMode: &accessMode,
						Server:     pointer.To(storage.NfsServer),
						ShareName:  pointer.To(storage.ShareName),
					},
				}
				managedEnvironmentStorage.Properties = props
			} else {
				props := &managedenvironmentsstorages.ManagedEnvironmentStorageProperties{
					AzureFile: &managedenvironmentsstorages.AzureFileProperties{
						AccessMode:  &accessMode,
						AccountKey:  pointer.To(storage.AccessKey),
						AccountName: pointer.To(storage.AccountName),
						ShareName:   pointer.To(storage.ShareName),
					},
				}
				managedEnvironmentStorage.Properties = props
			}

			if _, err := client.CreateOrUpdate(ctx, id, managedEnvironmentStorage); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppEnvironmentStorageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			id, err := managedenvironmentsstorages.ParseStorageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentStorageModel

			state.Name = id.StorageName
			state.ContainerAppEnvironmentId = managedenvironmentsstorages.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					if azureFile := props.AzureFile; azureFile != nil {
						state.AccountName = pointer.From(azureFile.AccountName)
						if azureFile.AccessMode != nil {
							state.AccessMode = string(*azureFile.AccessMode)
						}
						state.ShareName = pointer.From(azureFile.ShareName)
					} else if nfsAzureFile := props.NfsAzureFile; nfsAzureFile != nil {
						state.NfsServer = pointer.From(nfsAzureFile.Server)
						if nfsAzureFile.AccessMode != nil {
							state.AccessMode = string(*nfsAzureFile.AccessMode)
						}
						state.ShareName = pointer.From(nfsAzureFile.ShareName)
					}
				}
			}
			if keyFromConfig, ok := metadata.ResourceData.GetOk("access_key"); ok {
				state.AccessKey = keyFromConfig.(string)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentStorageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			id, err := managedenvironmentsstorages.ParseStorageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentStorageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			id, err := managedenvironmentsstorages.ParseStorageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var storage ContainerAppEnvironmentStorageModel
			if err := metadata.Decode(&storage); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s for update: %+v", *id, err)
			}

			if existing.Model.Properties == nil || existing.Model.Properties.AzureFile == nil {
				return fmt.Errorf("could not update %s: existing resource is missing `AzureFile` properties", *id)
			}

			// This *must* be sent, and is currently the only updatable property on the resource.
			existing.Model.Properties.AzureFile.AccountKey = pointer.To(metadata.ResourceData.Get("access_key").(string))

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
