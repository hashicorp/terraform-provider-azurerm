package machinelearning

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMachineLearningDataStoreDataLakeGen1() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMachineLearningDataStoreDataLakeGen1Create,
		Read:   resourceMachineLearningDataStoreDataLakeGen1Read,
		Update: resourceMachineLearningDataStoreDataLakeGen1Update,
		Delete: resourceMachineLearningDataStoreDataLakeGen1Delete,

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
		},
	}
	return resource
}

func resourceMachineLearningDataStoreDataLakeGen1Create(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_machine_learning_datastore_datalake_gen1", id.ID())
		}
	}

	datastoreRaw := datastore.DatastoreResource{
		Name: utils.String(d.Get("name").(string)),
		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
	}

	storeProps := &datastore.AzureDataLakeGen1Datastore{
		StoreName:   d.Get("store_name").(string),
		Description: utils.String(d.Get("description").(string)),
		IsDefault:   utils.Bool(d.Get("is_default").(bool)),
		Tags:        utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	if v, ok := d.GetOk("service_data_auth_identity"); ok {
		storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(v.(string)))
	}

	creds := map[string]interface{}{
		"credentialsType": "None",
	}

	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	if len(tenantId) != 0 && len(clientId) != 0 && len(clientSecret) != 0 {
		creds = map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeServicePrincipal),
			"authorityUrl":    d.Get("authority_url").(string),
			"resourceUrl":     "https://datalake.azure.net/",
			"tenantId":        tenantId,
			"clientId":        clientId,
			"secrets": map[string]interface{}{
				"secretsType":  "ServicePrincipal",
				"clientSecret": clientSecret,
			},
		}
	}
	storeProps.Credentials = creds
	datastoreRaw.Properties = storeProps

	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningDataStoreDataLakeGen1Read(d, meta)
}

func resourceMachineLearningDataStoreDataLakeGen1Update(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.DatastoreClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	datastoreRaw := datastore.DatastoreResource{
		Name: utils.String(d.Get("name").(string)),
		Type: utils.ToPtr(string(datastore.DatastoreTypeAzureBlob)),
	}

	storeProps := &datastore.AzureDataLakeGen1Datastore{
		StoreName:   d.Get("store_name").(string),
		Description: utils.String(d.Get("description").(string)),
		IsDefault:   utils.Bool(d.Get("is_default").(bool)),
		Tags:        utils.ToPtr(expandTags(d.Get("tags").(map[string]interface{}))),
	}

	if v, ok := d.GetOk("service_data_auth_identity"); ok {
		storeProps.ServiceDataAccessAuthIdentity = utils.ToPtr(datastore.ServiceDataAccessAuthIdentity(v.(string)))
	}

	creds := map[string]interface{}{
		"credentialsType": "None",
	}

	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	if len(tenantId) != 0 && len(clientId) != 0 && len(clientSecret) != 0 {
		creds = map[string]interface{}{
			"credentialsType": string(datastore.CredentialsTypeServicePrincipal),
			"authorityUrl":    d.Get("authority_url").(string),
			"resourceUrl":     "https://datalake.azure.net/",
			"tenantId":        tenantId,
			"clientId":        clientId,
			"secrets": map[string]interface{}{
				"secretsType":  "ServicePrincipal",
				"clientSecret": clientSecret,
			},
		}
	}
	storeProps.Credentials = creds
	datastoreRaw.Properties = storeProps

	_, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.DefaultCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf(" updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningDataStoreDataLakeGen1Read(d, meta)
}

func resourceMachineLearningDataStoreDataLakeGen1Read(d *pluginsdk.ResourceData, meta interface{}) error {
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

	data := resp.Model.Properties.(datastore.AzureDataLakeGen1Datastore)

	serviceDataAuth := ""
	if v := data.ServiceDataAccessAuthIdentity; v != nil {
		serviceDataAuth = string(*v)
	}
	d.Set("service_data_auth_identity", serviceDataAuth)
	d.Set("store_name", data.StoreName)

	if creds, ok := data.Credentials.(datastore.ServicePrincipalDatastoreCredentials); ok {
		if !strings.EqualFold(creds.TenantId, "00000000-0000-0000-0000-000000000000") && !strings.EqualFold(creds.ClientId, "00000000-0000-0000-0000-000000000000") {
			d.Set("tenant_id", creds.TenantId)
			d.Set("client_id", creds.ClientId)
		}
	}

	desc := ""
	if v := data.Description; v != nil {
		d.Set("description", desc)
	}

	d.Set("is_default", data.IsDefault)
	return flattenAndSetTags(d, *data.Tags)
}

func resourceMachineLearningDataStoreDataLakeGen1Delete(d *pluginsdk.ResourceData, meta interface{}) error {
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
