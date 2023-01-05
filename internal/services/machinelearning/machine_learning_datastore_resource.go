package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datastore.DatastoreTypeAzureBlob),
					string(datastore.DatastoreTypeAzureDataLakeGenOne),
					string(datastore.DatastoreTypeAzureDataLakeGenTwo),
					string(datastore.DatastoreTypeAzureFile),
				}, false),
			},

			"storage_account_name": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"file_share_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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

			"credentials": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_key": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							Sensitive:     true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"credentials.0.shared_access_signature", "credentials.0.tenant_id", "credentials.0.client_id", "credentials.0.client_secret"},
						},
						"shared_access_signature": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							Sensitive:     true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"credentials.0.account_key", "credentials.0.tenant_id", "credentials.0.client_id", "credentials.0.client_secret"},
						},

						"tenant_id": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.IsUUID,
							ConflictsWith: []string{"credentials.0.account_key", "credentials.0.shared_access_signature"},
						},

						"client_id": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.IsUUID,
							ConflictsWith: []string{"credentials.0.account_key", "credentials.0.shared_access_signature"},
						},

						"client_secret": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"credentials.0.account_key", "credentials.0.shared_access_signature"},
						},

						"authority_url": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							Computed:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"credentials.0.account_key", "credentials.0.shared_access_signature"},
						},
					},
				},
			},

			"store_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

	id := datastore.NewDataStoreID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_datastore", id.ID())
		}
	}

	datastoreRaw := datastore.DatastoreResource{
		Name: utils.String(d.Get("name").(string)),
		Type: utils.String(d.Get("type").(string)),
	}

	var prop datastore.Datastore
	var err error

	switch *datastoreRaw.Type {
	case string(datastore.DatastoreTypeAzureBlob):
		prop, err = expandBlobStorage(d)
	case string(datastore.DatastoreTypeAzureFile):
		prop, err = expandFileShare(d)
	case string(datastore.DatastoreTypeAzureDataLakeGenOne):
		prop, err = expandDataLakeGen1(d)
	case string(datastore.DatastoreTypeAzureDataLakeGenTwo):
		prop, err = expandDataLakeGen2(d)
	}

	if err != nil {
		return fmt.Errorf("%+v", err)
	}
	datastoreRaw.Properties = prop

	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningDataStoreRead(d, meta)
}

func resourceMachineLearningDataStoreRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
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
		return fmt.Errorf("making Read request on Machine Learning Data Store %q (Resource Group %q): %+v", id.Name, id.ResourceGroupName, err)
	}

	d.Set("name", resp.Model.Name)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_name", id.WorkspaceName)

	if prop, ok := resp.Model.Properties.(datastore.AzureBlobDatastore); ok {
		err = flattenBlobStorage(d, prop)
	}
	if prop, ok := resp.Model.Properties.(datastore.AzureFileDatastore); ok {
		err = flattenFileShare(d, prop)
	}
	if prop, ok := resp.Model.Properties.(datastore.AzureDataLakeGen1Datastore); ok {
		err = flattenDataLakeGen1(d, prop)
	}
	if prop, ok := resp.Model.Properties.(datastore.AzureDataLakeGen2Datastore); ok {
		err = flattenDataLakeGen2(d, prop)
	}

	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	return nil
}

func resourceMachineLearningDataStoreDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datastore.ParseDataStoreID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace Date Store ID `%q`: %+v", d.Id(), err)
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace Date Strore %q (Resource Group %q): %+v", id.Name, id.ResourceGroupName, err)
	}

	return nil
}

func expandBlobStorage(d *pluginsdk.ResourceData) (*datastore.AzureBlobDatastore, error) {
	storeProps := &datastore.AzureBlobDatastore{
		AccountName: utils.String(d.Get("storage_account_name").(string)),
		Description: utils.String(d.Get("description").(string)),
		IsDefault:   utils.Bool(d.Get("is_default").(bool)),
		Tags:        utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	containerName := d.Get("container_name").(string)
	if len(containerName) == 0 {
		return nil, fmt.Errorf(" `container_name` needs to be set if datastore type is `AzureBlob` ")
	}
	storeProps.ContainerName = utils.String(containerName)

	storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(d.Get("service_data_auth_identity").(string)))

	credentialRaw := d.Get("credentials").([]interface{})
	if credentialRaw != nil || len(credentialRaw) != 0 || credentialRaw[0] != nil {
		creds := credentialRaw[0].(map[string]interface{})
		if len(creds["account_key"].(string)) == 0 && len(creds["shared_access_signature"].(string)) == 0 {
			return nil, fmt.Errorf(" `account_key` or `shared_access_signature` needs to be set if datastore type is `AzureBlob` ")
		}
	}
	storeProps.Credentials = expandCredentials(d.Get("credentials").([]interface{}))
	return storeProps, nil
}

func flattenBlobStorage(d *pluginsdk.ResourceData, data datastore.AzureBlobDatastore) error {
	d.Set("description", data.Description)
	d.Set("is_default", data.IsDefault)
	d.Set("service_data_auth_identity", string(*data.ServiceDataAccessAuthIdentity))
	d.Set("type", datastore.DatastoreTypeAzureBlob)
	d.Set("storage_account_name", *data.AccountName)
	d.Set("container_name", *data.ContainerName)
	d.Set("credentials", flattenCredentials(d, data.Credentials))
	return flattenAndSetTags(d, *data.Tags)
}

func expandFileShare(d *pluginsdk.ResourceData) (*datastore.AzureFileDatastore, error) {
	accountName := d.Get("storage_account_name").(string)
	if len(accountName) == 0 {
		return nil, fmt.Errorf(" `storage_account_name` needs to be set if datastore type is `AzureFile` ")
	}

	fileShareName := d.Get("file_share_name").(string)
	if len(fileShareName) == 0 {
		return nil, fmt.Errorf(" `file_share_name` needs to be set if datastore type is `AzureFile` ")
	}

	storeProps := &datastore.AzureFileDatastore{
		AccountName:   accountName,
		FileShareName: fileShareName,
		Description:   utils.String(d.Get("description").(string)),
		IsDefault:     utils.Bool(d.Get("is_default").(bool)),
		Tags:          utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	if v, ok := d.GetOk("service_data_auth_identity"); ok {
		storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(v.(string)))
	}

	credentialRaw := d.Get("credentials").([]interface{})
	if credentialRaw != nil || len(credentialRaw) != 0 || credentialRaw[0] != nil {
		creds := credentialRaw[0].(map[string]interface{})
		if len(creds["account_key"].(string)) == 0 && len(creds["shared_access_signature"].(string)) == 0 {
			return nil, fmt.Errorf(" `account_key` or `shared_access_signature` needs to be set if datastore type is `AzureFile` ")
		}
	}

	storeProps.Credentials = expandCredentials(d.Get("credentials").([]interface{}))
	return storeProps, nil
}

func flattenFileShare(d *pluginsdk.ResourceData, data datastore.AzureFileDatastore) error {
	d.Set("description", data.Description)
	d.Set("is_default", data.IsDefault)
	d.Set("service_data_auth_identity", *data.ServiceDataAccessAuthIdentity)
	d.Set("type", datastore.DatastoreTypeAzureFile)
	d.Set("storage_account_name", data.AccountName)
	d.Set("file_share_name", data.FileShareName)
	d.Set("credentials", flattenCredentials(d, data.Credentials))
	return flattenAndSetTags(d, *data.Tags)
}

func expandDataLakeGen1(d *pluginsdk.ResourceData) (*datastore.AzureDataLakeGen1Datastore, error) {
	storeName := d.Get("store_name").(string)
	if len(storeName) == 0 {
		return nil, fmt.Errorf(" `store_name` needs to be set if datastore type is `AzureDataLakeGen1` ")
	}
	storeProps := &datastore.AzureDataLakeGen1Datastore{
		StoreName:   storeName,
		Description: utils.String(d.Get("description").(string)),
		IsDefault:   utils.Bool(d.Get("is_default").(bool)),
		Tags:        utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	if v, ok := d.GetOk("service_data_auth_identity"); ok {
		storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(v.(string)))
	}

	storeProps.Credentials = expandCredentials(d.Get("credentials").([]interface{}))
	return storeProps, nil
}

func flattenDataLakeGen1(d *pluginsdk.ResourceData, data datastore.AzureDataLakeGen1Datastore) error {
	d.Set("description", data.Description)
	d.Set("is_default", data.IsDefault)
	d.Set("service_data_auth_identity", data.ServiceDataAccessAuthIdentity)
	d.Set("type", datastore.DatastoreTypeAzureDataLakeGenOne)
	d.Set("store_name", data.StoreName)
	d.Set("credentials", flattenCredentials(d, data.Credentials))
	return flattenAndSetTags(d, *data.Tags)
}

func expandDataLakeGen2(d *pluginsdk.ResourceData) (*datastore.AzureDataLakeGen2Datastore, error) {
	accountName := d.Get("storage_account_name").(string)
	if len(accountName) == 0 {
		return nil, fmt.Errorf(" `storage_account_name` needs to be set if datastore type is `AzureDataLakeGen2` ")
	}

	fileSystem := d.Get("container_name").(string)
	if len(fileSystem) == 0 {
		return nil, fmt.Errorf(" `container_name` needs to be set if datastore type is `AzureDataLakeGen2` ")
	}
	storeProps := &datastore.AzureDataLakeGen2Datastore{
		AccountName: accountName,
		Filesystem:  fileSystem,
		Description: utils.String(d.Get("description").(string)),
		IsDefault:   utils.Bool(d.Get("is_default").(bool)),
		Tags:        utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	if v, ok := d.GetOk("service_data_auth_identity"); ok {
		storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(v.(string)))
	}

	storeProps.Credentials = expandCredentials(d.Get("credentials").([]interface{}))
	return storeProps, nil
}

func flattenDataLakeGen2(d *pluginsdk.ResourceData, data datastore.AzureDataLakeGen2Datastore) error {
	d.Set("storage_account_name", data.AccountName)
	d.Set("service_data_auth_identity", data.ServiceDataAccessAuthIdentity)
	d.Set("is_default", data.IsDefault)
	d.Set("type", datastore.DatastoreTypeAzureDataLakeGenTwo)
	d.Set("container_name", data.Filesystem)
	d.Set("storage_account_name", data.AccountName)
	d.Set("credentials", flattenCredentials(d, data.Credentials))
	return flattenAndSetTags(d, *data.Tags)
}

func expandCredentials(input []interface{}) map[string]interface{} {
	if len(input) == 0 || input[0] == nil {
		return map[string]interface{}{
			"credentialsType": "None",
		}
	}

	creds := input[0].(map[string]interface{})

	accountKey := creds["account_key"].(string)
	if len(accountKey) != 0 {
		return map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeAccountKey),
			"secrets": map[string]interface{}{
				"secretsType": "AccountKey",
				"key":         accountKey,
			},
		}
	}

	sasToken := creds["shared_access_signature"].(string)
	if len(sasToken) != 0 {
		return map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeSas),
			"secrets": map[string]interface{}{
				"secretsType": "Sas",
				"sasToken":    sasToken,
			},
		}
	}

	tenantId := creds["tenant_id"].(string)
	clientId := creds["client_id"].(string)
	clientSecret := creds["client_secret"].(string)
	if len(tenantId) != 0 && len(clientId) != 0 && len(clientSecret) != 0 {
		return map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeServicePrincipal),
			"authorityUrl":    creds["authority_url"].(string),
			"resourceUrl":     "https://datalake.azure.net/",
			"tenantId":        tenantId,
			"clientId":        clientId,
			"secrets": map[string]interface{}{
				"secretsType":  "ServicePrincipal",
				"clientSecret": clientSecret,
			},
		}
	}

	return map[string]interface{}{
		"credentialsType": "None",
	}
}

func flattenCredentials(d *pluginsdk.ResourceData, cred datastore.DatastoreCredentials) *[]interface{} {
	if _, ok := cred.(datastore.AccountKeyDatastoreCredentials); ok {
		return &[]interface{}{
			map[string]interface{}{
				"account_key": d.Get("credentials.0.account_key").(string),
			},
		}
	}

	if _, ok := cred.(datastore.SasDatastoreCredentials); ok {
		return &[]interface{}{
			map[string]interface{}{
				"shared_access_signature": d.Get("credentials.0.shared_access_signature").(string),
			},
		}
	}

	if v, ok := cred.(datastore.ServicePrincipalDatastoreCredentials); ok {
		return &[]interface{}{
			map[string]interface{}{
				"tenant_id":     v.TenantId,
				"client_id":     v.ClientId,
				"client_secret": d.Get("credentials.0.client_secret").(string),
			},
		}
	}

	if _, ok := cred.(datastore.NoneDatastoreCredentials); ok {
		return nil
	}

	return &[]interface{}{}
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
