package machinelearning

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMachineLearningDataStore() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMachineLearningDataStoreCreateOrUpdate,
		Read:   resourceMachineLearningDataStoreRead,
		Update: resourceMachineLearningDataStoreCreateOrUpdate,
		Delete: resourceMachineLearningDataStoreDelete,

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

			"container_name": {
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

func resourceMachineLearningDataStoreCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_machine_learning_datastore_blobstorage", id.ID())
		}
	}

	datastoreRaw := datastore.DatastoreResource{
		Name: utils.String(d.Get("name").(string)),
		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
	}

	props := &datastore.AzureBlobDatastore{
		AccountName:                   utils.String(d.Get("storage_account_name").(string)),
		ContainerName:                 utils.String(d.Get("container_name").(string)),
		Description:                   utils.String(d.Get("description").(string)),
		ServiceDataAccessAuthIdentity: utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(d.Get("service_data_auth_identity").(string))),
		IsDefault:                     utils.Bool(d.Get("is_default").(bool)),
		Tags:                          utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	accountKey := d.Get("account_key").(string)
	if accountKey != "" {
		props.Credentials = map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeAccountKey),
			"secrets": map[string]interface{}{
				"secretsType": "AccountKey",
				"key":         accountKey,
			},
		}
	}

	sasToken := d.Get("shared_access_signature").(string)
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

	d.SetId(id.ID())
	return resourceMachineLearningDataStoreRead(d, meta)
}

func resourceMachineLearningDataStoreRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datastore.ParseDataStoreID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("name", resp.Model.Name)
	d.Set("workspace_id", workspaceId.ID())

	data := resp.Model.Properties.(datastore.AzureBlobDatastore)

	serviceDataAuth := ""
	if v := data.ServiceDataAccessAuthIdentity; v != nil {
		serviceDataAuth = string(*v)
	}
	d.Set("service_data_auth_identity", serviceDataAuth)

	storageAccountName := ""
	if v := data.AccountName; v != nil {
		storageAccountName = *v
	}
	d.Set("storage_account_name", storageAccountName)

	containerName := ""
	if v := data.ContainerName; v != nil {
		containerName = *v
	}
	d.Set("container_name", containerName)

	desc := ""
	if v := data.Description; v != nil {
		d.Set("description", desc)
	}

	d.Set("is_default", data.IsDefault)
	return flattenAndSetTags(d, *data.Tags)
}

func resourceMachineLearningDataStoreDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datastore.ParseDataStoreID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandTags(tagsMap map[string]interface{}) map[string]string {
	output := make(map[string]string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := tags.TagValueToString(v)
		output[i] = value
	}

	return output
}

func flattenAndSetTags(d *pluginsdk.ResourceData, tagMap map[string]string) error {
	output := make(map[string]interface{}, len(tagMap))
	for i, v := range tagMap {
		output[i] = v
	}

	if err := d.Set("tags", output); err != nil {
		return fmt.Errorf("setting `tags`: %s", err)
	}

	return nil
}
