// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"strings"
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

type MachineLearningDataStoreDataLakeGen2 struct{}

type MachineLearningDataStoreDataLakeGen2Model struct {
	Name                string            `tfschema:"name"`
	WorkSpaceID         string            `tfschema:"workspace_id"`
	StorageContainerID  string            `tfschema:"storage_container_id"`
	TenantID            string            `tfschema:"tenant_id"`
	ClientID            string            `tfschema:"client_id"`
	ClientSecret        string            `tfschema:"client_secret"`
	AuthorityUrl        string            `tfschema:"authority_url"`
	Description         string            `tfschema:"description"`
	IsDefault           bool              `tfschema:"is_default"`
	ServiceDataIdentity string            `tfschema:"service_data_identity"`
	Tags                map[string]string `tfschema:"tags"`
}

func (r MachineLearningDataStoreDataLakeGen2) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"is_default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r MachineLearningDataStoreDataLakeGen2) ModelObject() interface{} {
	return &MachineLearningDataStoreDataLakeGen2Model{}
}

func (r MachineLearningDataStoreDataLakeGen2) ResourceType() string {
	return "azurerm_machine_learning_datastore_datalake_gen2"
}

func (r MachineLearningDataStoreDataLakeGen2) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastore.ValidateDataStoreID
}

var _ sdk.ResourceWithUpdate = MachineLearningDataStoreDataLakeGen2{}

func (r MachineLearningDataStoreDataLakeGen2) Arguments() map[string]*pluginsdk.Schema {
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

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"client_id", "client_secret"},
		},

		"client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"tenant_id", "client_secret"},
		},

		"client_secret": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"tenant_id", "client_id"},
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

		"authority_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (r MachineLearningDataStoreDataLakeGen2) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningDataStoreDataLakeGen2Model
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
				return tf.ImportAsExistsError("azurerm_machine_learning_datastore_datalake_gen2", id.ID())
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(model.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureDataLakeGenTwo)),
			}

			props := &datastore.AzureDataLakeGen2Datastore{
				AccountName:                   containerId.StorageAccountName,
				Filesystem:                    containerId.ContainerName,
				Description:                   utils.String(model.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(model.ServiceDataIdentity)),
				Tags:                          utils.ToPtr(model.Tags),
			}

			creds := map[string]interface{}{
				"credentialsType": "None",
			}

			if len(model.TenantID) != 0 && len(model.ClientID) != 0 && len(model.ClientSecret) != 0 {
				creds = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeServicePrincipal),
					"authorityUrl":    model.AuthorityUrl,
					"resourceUrl":     "https://datalake.azure.net/",
					"tenantId":        model.TenantID,
					"clientId":        model.ClientID,
					"secrets": map[string]interface{}{
						"secretsType":  "ServicePrincipal",
						"clientSecret": model.ClientSecret,
					},
				}
			}
			props.Credentials = creds
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

func (r MachineLearningDataStoreDataLakeGen2) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreDataLakeGen2Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}
			containerId, err := commonids.ParseStorageContainerID(state.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(id.DataStoreName),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureDataLakeGenTwo)),
			}

			props := &datastore.AzureDataLakeGen2Datastore{
				AccountName:                   containerId.StorageAccountName,
				Filesystem:                    containerId.ContainerName,
				Description:                   utils.String(state.Description),
				ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(state.ServiceDataIdentity)),
				Tags:                          utils.ToPtr(state.Tags),
			}

			creds := map[string]interface{}{
				"credentialsType": "None",
			}

			if len(state.TenantID) != 0 && len(state.ClientID) != 0 && len(state.ClientSecret) != 0 {
				creds = map[string]interface{}{
					"credentialsType": string(datastore.CredentialsTypeServicePrincipal),
					"authorityUrl":    state.AuthorityUrl,
					"resourceUrl":     "https://datalake.azure.net/",
					"tenantId":        state.TenantID,
					"clientId":        state.ClientID,
					"secrets": map[string]interface{}{
						"secretsType":  "ServicePrincipal",
						"clientSecret": state.ClientSecret,
					},
				}
			}
			props.Credentials = creds
			datastoreRaw.Properties = props

			_, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningDataStoreDataLakeGen2) Read() sdk.ResourceFunc {
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
			model := MachineLearningDataStoreDataLakeGen2Model{
				Name:        *resp.Model.Name,
				WorkSpaceID: workspaceId.ID(),
			}

			data := resp.Model.Properties.(datastore.AzureDataLakeGen2Datastore)
			serviceDataIdentity := ""
			if v := data.ServiceDataAccessAuthIdentity; v != nil {
				serviceDataIdentity = string(*v)
			}
			model.ServiceDataIdentity = serviceDataIdentity

			storageAccount, err := storageClient.FindAccount(ctx, data.AccountName)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Data Lake Gen2 File System %q: %s", data.AccountName, data.Filesystem, err)
			}
			if storageAccount == nil {
				return fmt.Errorf("Unable to locate Storage Account %q!", data.AccountName)
			}
			containerId := commonids.NewStorageContainerID(subscriptionId, storageAccount.ResourceGroup, data.AccountName, data.Filesystem)
			model.StorageContainerID = containerId.ID()

			model.IsDefault = *data.IsDefault

			if creds, ok := data.Credentials.(datastore.ServicePrincipalDatastoreCredentials); ok {
				if !strings.EqualFold(creds.TenantId, "00000000-0000-0000-0000-000000000000") && !strings.EqualFold(creds.ClientId, "00000000-0000-0000-0000-000000000000") {
					model.TenantID = creds.TenantId
					model.ClientID = creds.ClientId
					if v, ok := metadata.ResourceData.GetOk("client_secret"); ok {
						if v.(string) != "" {
							model.ClientSecret = v.(string)
						}
					}
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

func (r MachineLearningDataStoreDataLakeGen2) Delete() sdk.ResourceFunc {
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
