package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementApiTagDescription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiTagDescriptionCreateUpdate,
		Read:   resourceApiManagementApiTagDescriptionRead,
		Update: resourceApiManagementApiTagDescriptionCreateUpdate,
		Delete: resourceApiManagementApiTagDescriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiTagDescriptionsID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{

			"tag_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TagID,
			},

			"api_name": schemaz.SchemaApiManagementApiName(),

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"external_documentation_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"external_documentation_description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementApiTagDescriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tagId, err := parse.TagID(d.Get("tag_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	id := parse.NewApiTagDescriptionsID(tagId.SubscriptionId, tagId.ResourceGroup, tagId.ServiceName, d.Get("api_name").(string), tagId.Name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api_tag_description", id.ID())
		}
	}

	tagDescParameter := apimanagement.TagDescriptionCreateParameters{TagDescriptionBaseProperties: &apimanagement.TagDescriptionBaseProperties{}}
	if v, ok := d.GetOk("description"); ok {
		tagDescParameter.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_url"); ok {
		tagDescParameter.ExternalDocsURL = utils.String(v.(string))
	}

	if v, ok := d.GetOk("external_documentation_description"); ok {
		tagDescParameter.ExternalDocsDescription = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName, tagDescParameter, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiTagDescriptionRead(d, meta)
}

func resourceApiManagementApiTagDescriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagDescriptionsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	tagId := parse.NewTagID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.TagDescriptionName)

	d.Set("tag_id", tagId.ID())
	d.Set("api_name", id.ApiName)
	d.Set("description", resp.Description)
	d.Set("external_documentation_url", resp.ExternalDocsURL)
	d.Set("external_documentation_description", resp.ExternalDocsDescription)

	return nil
}

func resourceApiManagementApiTagDescriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiTagDescriptionsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
