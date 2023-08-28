// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStoreBlobStorage struct{}

type MachineLearningDataStoreBlobStorageModel struct {
	Name                    string            `tfschema:"name"`
	WorkSpaceID             string            `tfschema:"workspace_id"`
	StorageContainerID      string            `tfschema:"storage_container_id"`
	Description             string            `tfschema:"description"`
	IsDefault               bool              `tfschema:"is_default"`
	ServiceDataAuthIdentity string            `tfschema:"service_data_auth_identity"`
	AccountKey              string            `tfschema:"account_key"`
	SharedAccessSignature   string            `tfschema:"shared_access_signature"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (r MachineLearningDataStoreBlobStorage) Attributes() map[string]*schema.Schema {
	return nil
}

func (r MachineLearningDataStoreBlobStorage) ModelObject() interface{} {
	return &MachineLearningDataStoreBlobStorageModel{}
}

func (r MachineLearningDataStoreBlobStorage) ResourceType() string {
	return "azurerm_machine_learning_datastore_blobstorage"
}

func (r MachineLearningDataStoreBlobStorage) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastore.ValidateDataStoreID
}

var _ sdk.ResourceWithUpdate = MachineLearningDataStoreBlobStorage{}

func (r MachineLearningDataStoreBlobStorage) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataStoreName,
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"storage_container_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"is_default": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_data_auth_identity": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(datastore.ServiceDataAccessAuthIdentityNone),
				string(datastore.ServiceDataAccessAuthIdentityWorkspaceSystemAssignedIdentity),
				string(datastore.ServiceDataAccessAuthIdentityWorkspaceUserAssignedIdentity),
			},
				false),
			Default: string(datastore.ServiceDataAccessAuthIdentityNone),
		},

		"account_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"account_key", "shared_access_signature"},
		},

		"shared_access_signature": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"account_key", "shared_access_signature"},
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (r MachineLearningDataStoreBlobStorage) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningDataStoreBlobStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkSpaceID)
			if err != nil {
				return err
			}

			id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_datastore_blobstorage", id.ID())
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(model.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
			}

			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			props := &datastore.AzureBlobDatastore{
				AccountName:                   utils.String(containerId.StorageAccountName),
				Endpoint:                      storageDomainSuffix,
				ContainerName:                 utils.String(containerId.ContainerName),
				Description:                   utils.String(model.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(model.ServiceDataAuthIdentity)),
				IsDefault:                     utils.Bool(model.IsDefault),
				Tags:                          utils.ToPtr(model.Tags),
			}

			accountKey := model.AccountKey
			if accountKey != "" {
				props.Credentials = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeAccountKey),
					"secrets": map[string]interface{}{
						"secretsType": "AccountKey",
						"key":         accountKey,
					},
				}
			}

			sasToken := model.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeSas),
					"secrets": map[string]interface{}{
						"secretsType": "Sas",
						"sasToken":    sasToken,
					},
				}
			}
			datastoreRaw.Properties = props

			_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningDataStoreBlobStorage) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreBlobStorageModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			containerId, err := commonids.ParseStorageContainerID(state.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(id.DataStoreName),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
			}

			props := &datastore.AzureBlobDatastore{
				AccountName:                   utils.String(containerId.StorageAccountName),
				ContainerName:                 utils.String(containerId.ContainerName),
				Description:                   utils.String(state.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(state.ServiceDataAuthIdentity)),
				IsDefault:                     utils.Bool(state.IsDefault),
				Tags:                          utils.ToPtr(state.Tags),
			}

			accountKey := state.AccountKey
			if accountKey != "" {
				props.Credentials = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeAccountKey),
					"secrets": map[string]interface{}{
						"secretsType": "AccountKey",
						"key":         accountKey,
					},
				}
			}

			sasToken := state.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeSas),
					"secrets": map[string]interface{}{
						"secretsType": "Sas",
						"sasToken":    sasToken,
					},
				}
			}
			datastoreRaw.Properties = props

			_, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningDataStoreBlobStorage) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
			model := MachineLearningDataStoreBlobStorageModel{
				Name:        *resp.Model.Name,
				WorkSpaceID: workspaceId.ID(),
			}

			data := resp.Model.Properties.(datastore.AzureBlobDatastore)
			serviceDataAuth := ""
			if v := data.ServiceDataAccessAuthIdentity; v != nil {
				serviceDataAuth = string(*v)
			}
			model.ServiceDataAuthIdentity = serviceDataAuth

			containerId := commonids.NewStorageContainerID(subscriptionId, workspaceId.ResourceGroupName, *data.AccountName, *data.ContainerName)
			model.StorageContainerID = containerId.ID()
			model.IsDefault = *data.IsDefault

			if v, ok := metadata.ResourceData.GetOk("account_key"); ok {
				if v.(string) != "" {
					model.AccountKey = v.(string)
				}
			}

			if v, ok := metadata.ResourceData.GetOk("shared_access_signature"); ok {
				if v.(string) != "" {
					model.SharedAccessSignature = v.(string)
				}
			}

			desc := ""
			if v := data.Description; v != nil {
				desc = *v
			}
			model.Description = desc

			if data.Tags != nil {
				model.Tags = *data.Tags
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MachineLearningDataStoreBlobStorage) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
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
