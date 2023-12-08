// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	storageparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStoreFileShare struct{}

type MachineLearningDataStoreFileShareModel struct {
	Name                  string            `tfschema:"name"`
	WorkSpaceID           string            `tfschema:"workspace_id"`
	StorageFileShareID    string            `tfschema:"storage_fileshare_id"`
	Description           string            `tfschema:"description"`
	IsDefault             bool              `tfschema:"is_default"`
	ServiceDataIdentity   string            `tfschema:"service_data_identity"`
	AccountKey            string            `tfschema:"account_key"`
	SharedAccessSignature string            `tfschema:"shared_access_signature"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (r MachineLearningDataStoreFileShare) ModelObject() interface{} {
	return &MachineLearningDataStoreFileShareModel{}
}

func (r MachineLearningDataStoreFileShare) ResourceType() string {
	return "azurerm_machine_learning_datastore_fileshare"
}

func (r MachineLearningDataStoreFileShare) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastore.ValidateDataStoreID
}

var _ sdk.ResourceWithUpdate = MachineLearningDataStoreFileShare{}

func (r MachineLearningDataStoreFileShare) Arguments() map[string]*pluginsdk.Schema {
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

		"storage_fileshare_id": {
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

		"service_data_identity": {
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

func (r MachineLearningDataStoreFileShare) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"is_default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r MachineLearningDataStoreFileShare) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningDataStoreFileShareModel
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
				return tf.ImportAsExistsError("azurerm_machine_learning_datastore_fileshare", id.ID())
			}

			fileShareId, err := storageparse.StorageShareResourceManagerID(model.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(model.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   fileShareId.StorageAccountName,
				FileShareName:                 fileShareId.FileshareName,
				Description:                   utils.String(model.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(model.ServiceDataIdentity)),
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
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningDataStoreFileShare) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreFileShareModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			fileShareId, err := storageparse.StorageShareResourceManagerID(state.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(id.DataStoreName),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   fileShareId.StorageAccountName,
				FileShareName:                 fileShareId.FileshareName,
				Description:                   utils.String(state.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(state.ServiceDataIdentity)),
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

func (r MachineLearningDataStoreFileShare) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			storageClient := metadata.Client.Storage
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
			model := MachineLearningDataStoreFileShareModel{
				Name:        *resp.Model.Name,
				WorkSpaceID: workspaceId.ID(),
			}

			data := resp.Model.Properties.(datastore.AzureFileDatastore)
			serviceDataIdentity := ""
			if v := data.ServiceDataAccessAuthIdentity; v != nil {
				serviceDataIdentity = string(*v)
			}
			model.ServiceDataIdentity = serviceDataIdentity

			storageAccount, err := storageClient.FindAccount(ctx, data.AccountName)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Share %q: %s", data.AccountName, data.FileShareName, err)
			}
			if storageAccount == nil {
				return fmt.Errorf("Unable to locate Storage Account %q!", data.AccountName)
			}
			fileShareId := storageparse.NewStorageShareResourceManagerID(subscriptionId, storageAccount.ResourceGroup, data.AccountName, "default", data.FileShareName)
			model.StorageFileShareID = fileShareId.ID()

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

func (r MachineLearningDataStoreFileShare) Delete() sdk.ResourceFunc {
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
