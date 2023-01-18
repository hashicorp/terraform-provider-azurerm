package machinelearning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStoreDataLakeGen1 struct{}

type MachineLearningDataStoreDataLakeGen1Model struct {
	Name                    string            `tfschema:"name"`
	WorkSpaceID             string            `tfschema:"workspace_id"`
	StoreName               string            `tfschema:"store_name"`
	TenantID                string            `tfschema:"tenant_id"`
	ClientID                string            `tfschema:"client_id"`
	ClientSecret            string            `tfschema:"client_secret"`
	AuthorityUrl            string            `tfschema:"authority_url"`
	Description             string            `tfschema:"description"`
	IsDefault               bool              `tfschema:"is_default"`
	ServiceDataAuthIdentity string            `tfschema:"service_data_auth_identity"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (r MachineLearningDataStoreDataLakeGen1) Attributes() map[string]*schema.Schema {
	return nil
}

func (r MachineLearningDataStoreDataLakeGen1) ModelObject() interface{} {
	return &MachineLearningDataStoreDataLakeGen1Model{}
}

func (r MachineLearningDataStoreDataLakeGen1) ResourceType() string {
	return "azurerm_machine_learning_datastore_datalake_gen1"
}

func (r MachineLearningDataStoreDataLakeGen1) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datastore.ValidateDataStoreID
}

var _ sdk.ResourceWithUpdate = MachineLearningDataStoreDataLakeGen1{}

func (r MachineLearningDataStoreDataLakeGen1) Arguments() map[string]*pluginsdk.Schema {
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

		"store_name": {
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

		"authority_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (r MachineLearningDataStoreDataLakeGen1) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.DatastoreClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningDataStoreDataLakeGen1Model
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
				return tf.ImportAsExistsError("azurerm_machine_learning_datastore_datalake_gen1", id.ID())
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(model.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
			}

			storeProps := &datastore.AzureDataLakeGen1Datastore{
				StoreName:   model.StoreName,
				Description: utils.String(model.Description),
				IsDefault:   utils.Bool(model.IsDefault),
				Tags:        utils.ToPtr(model.Tags),
			}

			if model.ServiceDataAuthIdentity != "" {
				storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(model.ServiceDataAuthIdentity))
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
			storeProps.Credentials = creds
			datastoreRaw.Properties = storeProps

			_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningDataStoreDataLakeGen1) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.DatastoreClient

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreDataLakeGen1Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: utils.String(state.Name),
				Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
			}

			storeProps := &datastore.AzureDataLakeGen1Datastore{
				StoreName:   state.StoreName,
				Description: utils.String(state.Description),
				IsDefault:   utils.Bool(state.IsDefault),
				Tags:        utils.ToPtr(state.Tags),
			}

			if state.ServiceDataAuthIdentity != "" {
				storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(state.ServiceDataAuthIdentity))
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
			storeProps.Credentials = creds
			datastoreRaw.Properties = storeProps

			_, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf(" updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningDataStoreDataLakeGen1) Read() sdk.ResourceFunc {
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
			model := MachineLearningDataStoreDataLakeGen1Model{
				Name:        *resp.Model.Name,
				WorkSpaceID: workspaceId.ID(),
			}

			data := resp.Model.Properties.(datastore.AzureDataLakeGen1Datastore)
			serviceDataAuth := ""
			if v := data.ServiceDataAccessAuthIdentity; v != nil {
				serviceDataAuth = string(*v)
			}
			model.ServiceDataAuthIdentity = serviceDataAuth

			model.StoreName = data.StoreName
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

func (r MachineLearningDataStoreDataLakeGen1) Delete() sdk.ResourceFunc {
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
