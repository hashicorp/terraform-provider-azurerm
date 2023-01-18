package machinelearning

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	storageparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
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
			ValidateFunc: validate.WorkspaceID,
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
			client := metadata.Client.MachineLearning.DatastoreClient
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

			containerId, err := storageparse.StorageContainerResourceManagerID(model.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(model.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
			}

			props := &datastore.AzureBlobDatastore{
				AccountName:                   utils.String(containerId.StorageAccountName),
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
			client := metadata.Client.MachineLearning.DatastoreClient

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreBlobStorageModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			containerId, err := storageparse.StorageContainerResourceManagerID(state.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(id.Name),
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
			client := metadata.Client.MachineLearning.DatastoreClient
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

			containerId := storageparse.NewStorageContainerResourceManagerID(subscriptionId, workspaceId.ResourceGroupName, *data.AccountName, "default", *data.ContainerName)
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
			client := metadata.Client.MachineLearning.DatastoreClient

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

//
//func resourceMachineLearningDataStoreBlobStorage() *pluginsdk.Resource {
//	resource := &pluginsdk.Resource{
//		Create: resourceMachineLearningDataStoreBlobStorageCreate,
//		Read:   resourceMachineLearningDataStoreBlobStorageRead,
//		Update: resourceMachineLearningDataStoreBlobStorageUpdate,
//		Delete: resourceMachineLearningDataStoreBlobStorageDelete,
//
//		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
//			_, err := datastore.ParseDataStoreID(id)
//			return err
//		}),
//
//		Timeouts: &pluginsdk.ResourceTimeout{
//			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
//			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
//			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
//			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
//		},
//
//		Schema: map[string]*pluginsdk.Schema{
//			"name": {
//				Type:         pluginsdk.TypeString,
//				Required:     true,
//				ForceNew:     true,
//				ValidateFunc: validate.DataStoreName,
//			},
//
//			"workspace_id": {
//				Type:         pluginsdk.TypeString,
//				Required:     true,
//				ForceNew:     true,
//				ValidateFunc: validate.WorkspaceID,
//			},
//
//			"storage_container_id": {
//				Type:         pluginsdk.TypeString,
//				Required:     true,
//				ForceNew:     true,
//				ValidateFunc: validation.StringIsNotEmpty,
//			},
//
//			"description": {
//				Type:     pluginsdk.TypeString,
//				Optional: true,
//				ForceNew: true,
//			},
//
//			"is_default": {
//				Type:     pluginsdk.TypeBool,
//				Optional: true,
//				Default:  false,
//			},
//
//			"service_data_auth_identity": {
//				Type:     pluginsdk.TypeString,
//				Optional: true,
//				ValidateFunc: validation.StringInSlice([]string{
//					string(datastore.ServiceDataAccessAuthIdentityNone),
//					string(datastore.ServiceDataAccessAuthIdentityWorkspaceSystemAssignedIdentity),
//					string(datastore.ServiceDataAccessAuthIdentityWorkspaceUserAssignedIdentity),
//				},
//					false),
//				Default: string(datastore.ServiceDataAccessAuthIdentityNone),
//			},
//
//			"account_key": {
//				Type:         pluginsdk.TypeString,
//				Optional:     true,
//				Sensitive:    true,
//				ValidateFunc: validation.StringIsNotEmpty,
//				ExactlyOneOf: []string{"account_key", "shared_access_signature"},
//			},
//
//			"shared_access_signature": {
//				Type:         pluginsdk.TypeString,
//				Optional:     true,
//				Sensitive:    true,
//				ValidateFunc: validation.StringIsNotEmpty,
//				AtLeastOneOf: []string{"account_key", "shared_access_signature"},
//			},
//
//			"tags": commonschema.TagsForceNew(),
//		},
//	}
//	return resource
//}
//
//func resourceMachineLearningDataStoreBlobStorageCreate(d *pluginsdk.ResourceData, meta interface{}) error {
//	client := meta.(*clients.Client).MachineLearning.DatastoreClient
//	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
//	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
//	defer cancel()
//
//	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
//	if err != nil {
//		return err
//	}
//
//	id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
//	if d.IsNewResource() {
//		existing, err := client.Get(ctx, id)
//		if err != nil {
//			if !response.WasNotFound(existing.HttpResponse) {
//				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
//			}
//		}
//		if !response.WasNotFound(existing.HttpResponse) {
//			return tf.ImportAsExistsError("azurerm_machine_learning_datastore_blobstorage", id.ID())
//		}
//	}
//
//	containerId, err := storageparse.StorageContainerResourceManagerID(d.Get("storage_container_id").(string))
//	if err != nil {
//		return err
//	}
//
//	datastoreRaw := datastore.DatastoreResource{
//		Name: utils.String(d.Get("name").(string)),
//		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
//	}
//
//	props := &datastore.AzureBlobDatastore{
//		AccountName:                   utils.String(containerId.StorageAccountName),
//		ContainerName:                 utils.String(containerId.ContainerName),
//		Description:                   utils.String(d.Get("description").(string)),
//		ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(d.Get("service_data_auth_identity").(string))),
//		IsDefault:                     utils.Bool(d.Get("is_default").(bool)),
//		Tags:                          utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
//	}
//
//	accountKey := d.Get("account_key").(string)
//	if accountKey != "" {
//		props.Credentials = map[string]interface{}{
//			"credentialsType": string(datastore.CredentialsTypeAccountKey),
//			"secrets": map[string]interface{}{
//				"secretsType": "AccountKey",
//				"key":         accountKey,
//			},
//		}
//	}
//
//	sasToken := d.Get("shared_access_signature").(string)
//	if sasToken != "" {
//		props.Credentials = map[string]interface{}{
//			"credentialsType": string(datastore.CredentialsTypeSas),
//			"secrets": map[string]interface{}{
//				"secretsType": "Sas",
//				"sasToken":    sasToken,
//			},
//		}
//	}
//	datastoreRaw.Properties = props
//
//	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
//	if err != nil {
//		return fmt.Errorf("creating/updating %s: %+v", id, err)
//	}
//
//	d.SetId(id.ID())
//	return resourceMachineLearningDataStoreBlobStorageRead(d, meta)
//}
//
//func resourceMachineLearningDataStoreBlobStorageUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
//	client := meta.(*clients.Client).MachineLearning.DatastoreClient
//	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
//	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
//	defer cancel()
//
//	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
//	if err != nil {
//		return err
//	}
//
//	id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
//
//	containerId, err := storageparse.StorageContainerResourceManagerID(d.Get("storage_container_id").(string))
//	if err != nil {
//		return err
//	}
//
//	datastoreRaw := datastore.DatastoreResource{
//		Name: utils.String(id.Name),
//		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
//	}
//
//	props := &datastore.AzureBlobDatastore{
//		AccountName:                   utils.String(containerId.StorageAccountName),
//		ContainerName:                 utils.String(containerId.ContainerName),
//		Description:                   utils.String(d.Get("description").(string)),
//		ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(d.Get("service_data_auth_identity").(string))),
//		IsDefault:                     utils.Bool(d.Get("is_default").(bool)),
//		Tags:                          utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
//	}
//
//	accountKey := d.Get("account_key").(string)
//	if accountKey != "" {
//		props.Credentials = map[string]interface{}{
//			"credentialsType": string(datastore.CredentialsTypeAccountKey),
//			"secrets": map[string]interface{}{
//				"secretsType": "AccountKey",
//				"key":         accountKey,
//			},
//		}
//	}
//
//	sasToken := d.Get("shared_access_signature").(string)
//	if sasToken != "" {
//		props.Credentials = map[string]interface{}{
//			"credentialsType": string(datastore.CredentialsTypeSas),
//			"secrets": map[string]interface{}{
//				"secretsType": "Sas",
//				"sasToken":    sasToken,
//			},
//		}
//	}
//	datastoreRaw.Properties = props
//
//	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
//	if err != nil {
//		return fmt.Errorf("creating/updating %s: %+v", id, err)
//	}
//
//	d.SetId(id.ID())
//	return resourceMachineLearningDataStoreBlobStorageRead(d, meta)
//}
//
//func resourceMachineLearningDataStoreBlobStorageRead(d *pluginsdk.ResourceData, meta interface{}) error {
//	client := meta.(*clients.Client).MachineLearning.DatastoreClient
//	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
//	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
//	defer cancel()
//
//	id, err := datastore.ParseDataStoreID(d.Id())
//	if err != nil {
//		return err
//	}
//
//	resp, err := client.Get(ctx, *id)
//	if err != nil {
//		if response.WasNotFound(resp.HttpResponse) {
//			d.SetId("")
//			return nil
//		}
//		return fmt.Errorf("reading %s: %+v", *id, err)
//	}
//
//	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
//	d.Set("name", resp.Model.Name)
//	d.Set("workspace_id", workspaceId.ID())
//
//	data := resp.Model.Properties.(datastore.AzureBlobDatastore)
//
//	serviceDataAuth := ""
//	if v := data.ServiceDataAccessAuthIdentity; v != nil {
//		serviceDataAuth = string(*v)
//	}
//	d.Set("service_data_auth_identity", serviceDataAuth)
//
//	containerId := storageparse.NewStorageContainerResourceManagerID(subscriptionId, workspaceId.ResourceGroupName, *data.AccountName, "default", *data.ContainerName)
//	d.Set("storage_container_id", containerId.ID())
//
//	desc := ""
//	if v := data.Description; v != nil {
//		d.Set("description", desc)
//	}
//
//	d.Set("is_default", data.IsDefault)
//	return flattenAndSetTags(d, *data.Tags)
//}
//
//func resourceMachineLearningDataStoreBlobStorageDelete(d *pluginsdk.ResourceData, meta interface{}) error {
//	client := meta.(*clients.Client).MachineLearning.DatastoreClient
//	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
//	defer cancel()
//
//	id, err := datastore.ParseDataStoreID(d.Id())
//	if err != nil {
//		return err
//	}
//
//	if _, err := client.Delete(ctx, *id); err != nil {
//		return fmt.Errorf("deleting %s: %+v", *id, err)
//	}
//
//	return nil
//}
//
//func expandTags(tagsMap map[string]interface{}) map[string]string {
//	output := make(map[string]string, len(tagsMap))
//
//	for i, v := range tagsMap {
//		// Validate should have ignored this error already
//		value, _ := tags.TagValueToString(v)
//		output[i] = value
//	}
//
//	return output
//}
//
//func flattenAndSetTags(d *pluginsdk.ResourceData, tagMap map[string]string) error {
//	output := make(map[string]interface{}, len(tagMap))
//	for i, v := range tagMap {
//		output[i] = v
//	}
//
//	if err := d.Set("tags", output); err != nil {
//		return fmt.Errorf("setting `tags`: %s", err)
//	}
//
//	return nil
//}
