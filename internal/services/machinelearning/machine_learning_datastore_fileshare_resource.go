package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMachineLearningFileShare() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMachineLearningFileShareCreateOrUpdate,
		Read:   resourceMachineLearningFileShareRead,
		Update: resourceMachineLearningFileShareCreateOrUpdate,
		Delete: resourceMachineLearningFileShareDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := datastore.ParseDataStoreID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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

			"storage_account_name": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"file_share_name": {
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
		},
	}
	return resource
}

func resourceMachineLearningFileShareCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
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

	datastoreRaw := datastore.DatastoreResource{
		Name: utils.String(d.Get("name").(string)),
		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureFile)),
	}

	prop := expandFileShare(d)
	datastoreRaw.Properties = prop

	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningFileShareRead(d, meta)
}

func resourceMachineLearningFileShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datastore.ParseDataStoreID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Data Store ID `%q`: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Machine Learning Data Store File Share %q (Resource Group %q): %+v", id.Name, id.ResourceGroupName, err)
	}

	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("name", resp.Model.Name)
	d.Set("workspace_id", workspaceId.ID())
	return flattenFileShare(d, resp.Model.Properties.(datastore.AzureFileDatastore))
}

func resourceMachineLearningFileShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datastore.ParseDataStoreID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace Date Store ID `%q`: %+v", d.Id(), err)
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace Date Strore  File Share %q (Resource Group %q): %+v", id.Name, id.ResourceGroupName, err)
	}

	return nil
}

func expandFileShare(d *pluginsdk.ResourceData) *datastore.AzureFileDatastore {
	storeProps := &datastore.AzureFileDatastore{
		AccountName:                   d.Get("storage_account_name").(string),
		FileShareName:                 d.Get("file_share_name").(string),
		Description:                   utils.String(d.Get("description").(string)),
		ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(d.Get("service_data_auth_identity").(string))),
		IsDefault:                     utils.Bool(d.Get("is_default").(bool)),
		Tags:                          utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	accountKey := d.Get("account_key").(string)
	if accountKey != "" {
		storeProps.Credentials = map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeAccountKey),
			"secrets": map[string]interface{}{
				"secretsType": "AccountKey",
				"key":         accountKey,
			},
		}
	}

	sasToken := d.Get("shared_access_signature").(string)
	if sasToken != "" {
		storeProps.Credentials = map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeSas),
			"secrets": map[string]interface{}{
				"secretsType": "Sas",
				"sasToken":    sasToken,
			},
		}
	}

	return storeProps
}

func flattenFileShare(d *pluginsdk.ResourceData, data datastore.AzureFileDatastore) error {
	d.Set("description", data.Description)
	d.Set("is_default", data.IsDefault)
	d.Set("service_data_auth_identity", string(*data.ServiceDataAccessAuthIdentity))
	d.Set("storage_account_name", data.AccountName)
	d.Set("file_share_name", data.FileShareName)
	return flattenAndSetTags(d, *data.Tags)
}
