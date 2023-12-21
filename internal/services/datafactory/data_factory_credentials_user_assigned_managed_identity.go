package datafactory

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryCredentialsUserAssignedManagedIdentity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryCredentialsUserAssignedManagedIdentityCreateUpdate,
		Read:   resourceDataFactoryCredentialsUserAssignedManagedIdentityRead,
		Update: resourceDataFactoryCredentialsUserAssignedManagedIdentityCreateUpdate,
		// Delete: resourceDataFactoryCustomDatasetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"data_factory_id": {
				Description:  "ID of the Data Factory",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},
			"description": {
				Description: "Text description of the credential",
				Type:        pluginsdk.TypeString,
				Optional:    true,
			},
			"identity_id": {
				Description: "Resource ID of a User-Assigned Managed Identity",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Credential Name",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ForceNew:    true, // TODO: figure out whats required
			},
		},
	}
}

// user managed identities only have one type
const IDENTITY_TYPE = "ManagedIdentity"

func resourceDataFactoryCredentialsUserAssignedManagedIdentityCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Credentials
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := credentials.CredentialId{
		SubscriptionId:    dataFactoryId.SubscriptionId,
		ResourceGroupName: dataFactoryId.ResourceGroupName,
		FactoryName:       dataFactoryId.FactoryName,
		CredentialName:    d.Get("name").(string),
	}

	if d.IsNewResource() {
		existing, err := client.CredentialOperationsGet(ctx, id, credentials.CredentialOperationsGetOperationOptions{})
		if err != nil {
			if existing.HttpResponse.Status == "404" {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.HttpResponse.Status == "404" {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_http", id.ID())
		}
	}

	credential := credentials.ManagedIdentityCredentialResource{
		Type:       utils.String(IDENTITY_TYPE),
		Properties: credentials.ManagedIdentityCredential{},
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		credential.Properties.Annotations = &annotations
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		credential.Properties.Description = &description
	}

	if v, ok := d.GetOk("identity_id"); ok {
		identityId := v.(string)
		credential.Properties.TypeProperties.ResourceId = &identityId
	}

	if _, err := client.CredentialOperationsCreateOrUpdate(ctx, id, credential, credentials.CredentialOperationsCreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryCredentialsUserAssignedManagedIdentityRead(d, meta)
}

func resourceDataFactoryCredentialsUserAssignedManagedIdentityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Credentials
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	credentialId, err := parse.CredentialID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.CredentialOperationsGet(ctx, *credentialId, credentials.CredentialOperationsGetOperationOptions{})
	if err != nil {
		if existing.HttpResponse.Status == "404" {
			return fmt.Errorf("checking for presence of existing %s: %+v", d.Id(), err)
		}
	}

	d.Set("name", credentialId.CredentialName)
	d.Set("data_factory_id", credentialId.FactoryName)
	d.Set("description", existing.Model.Properties.Description)
	if err := d.Set("annotations", flattenDataFactoryAnnotations(existing.Model.Properties.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	return nil
}
