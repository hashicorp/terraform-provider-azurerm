// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	storageAccountHelper "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	storageparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storagevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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

var (
	_ sdk.ResourceWithUpdate         = MachineLearningDataStoreFileShare{}
	_ sdk.ResourceWithStateMigration = MachineLearningDataStoreFileShare{}
)

func (r MachineLearningDataStoreFileShare) StateUpgraders() sdk.StateUpgradeData {
	// rewrite `storage_fileshare_id` to account for azurerm_storage_share.resource_manager_id deprecation
	// in favour of azurerm_storage_share.id in 5.0
	if !features.FivePointOh() {
		return sdk.StateUpgradeData{}
	}

	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.MachineLearningDataStoreFileShareV0ToV1{},
		},
	}
}

func (r MachineLearningDataStoreFileShare) Arguments() map[string]*pluginsdk.Schema {
	res := map[string]*pluginsdk.Schema{
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
			ValidateFunc: fileshares.ValidateShareID,
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

	if !features.FivePointOh() {
		res["storage_fileshare_id"].ValidateFunc = storagevalidate.StorageShareResourceManagerID
	}

	return res
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_machine_learning_datastore_fileshare", id.ID())
				}
			}

			shareId, err := parseStorageFileShareID(model.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   shareId.StorageAccountName,
				Endpoint:                      pointer.To(metadata.Client.Storage.StorageDomainSuffix),
				FileShareName:                 shareId.ShareName,
				Description:                   pointer.To(model.Description),
				ServiceDataAccessAuthIdentity: pointer.To(datastore.ServiceDataAccessAuthIdentity(model.ServiceDataIdentity)),
				Tags:                          pointer.To(model.Tags),
			}

			accountKey := model.AccountKey
			if accountKey != "" {
				props.Credentials = datastore.AccountKeyDatastoreCredentials{
					Secrets: datastore.AccountKeyDatastoreSecrets{
						Key: pointer.To(accountKey),
					},
				}
			}

			sasToken := model.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = datastore.SasDatastoreCredentials{
					Secrets: datastore.SasDatastoreSecrets{
						SasToken: pointer.To(sasToken),
					},
				}
			}
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
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

			shareId, err := parseStorageFileShareID(state.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(id.DataStoreName),
				Type: pointer.To(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   shareId.StorageAccountName,
				FileShareName:                 shareId.ShareName,
				Description:                   pointer.To(state.Description),
				ServiceDataAccessAuthIdentity: pointer.To(datastore.ServiceDataAccessAuthIdentity(state.ServiceDataIdentity)),
				Tags:                          pointer.To(state.Tags),
			}

			accountKey := state.AccountKey
			if accountKey != "" {
				props.Credentials = datastore.AccountKeyDatastoreCredentials{
					Secrets: datastore.AccountKeyDatastoreSecrets{
						Key: pointer.To(accountKey),
					},
				}
			}

			sasToken := state.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = datastore.SasDatastoreCredentials{
					Secrets: datastore.SasDatastoreSecrets{
						SasToken: pointer.To(sasToken),
					},
				}
			}
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
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

			var storageAccount *storageAccountHelper.AccountDetails
			if fileShareIdStr := metadata.ResourceData.Get("storage_fileshare_id").(string); fileShareIdStr != "" {
				shareId, err := parseStorageFileShareID(fileShareIdStr)
				if err != nil {
					return err
				}
				storageAccount, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.StorageAccountName))
				if err != nil {
					return fmt.Errorf("retrieving Account %q for Share %q: %s", data.AccountName, data.FileShareName, err)
				}
			} else {
				// In the case of import, we cannot rely on having a value for `storage_container_id` so we need to fallback on listing the accounts to search.
				storageAccount, err = storageClient.FindAccount(ctx, subscriptionId, data.AccountName)
				if err != nil {
					return fmt.Errorf("retrieving Account %q for Share %q: %s", data.AccountName, data.FileShareName, err)
				}
			}

			if storageAccount == nil {
				return fmt.Errorf("unable to locate Storage Account %q", data.AccountName)
			}

			model.StorageFileShareID = fileshares.NewShareID(storageAccount.StorageAccountId.SubscriptionId, storageAccount.StorageAccountId.ResourceGroupName, data.AccountName, data.FileShareName).ID()
			if !features.FivePointOh() {
				model.StorageFileShareID = storageparse.NewStorageShareResourceManagerID(storageAccount.StorageAccountId.SubscriptionId, storageAccount.StorageAccountId.ResourceGroupName, data.AccountName, "default", data.FileShareName).ID()
			}

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

// parses the `storage_fileshare_id` into a fileshares.ShareId needed for 5.0 where the id changes
// from  "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/fileshares/%s"
// to "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
// since we only use id.StorageAccountName and id.FileshareName parsing into a fileshares.ShareId works for both 4.0 and 5.0
func parseStorageFileShareID(input string) (*fileshares.ShareId, error) {
	if !features.FivePointOh() {
		id, err := storageparse.StorageShareResourceManagerID(input)
		if err != nil {
			return nil, err
		}
		return pointer.To(fileshares.NewShareID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileshareName)), nil
	}
	return fileshares.ParseShareID(input)
}
