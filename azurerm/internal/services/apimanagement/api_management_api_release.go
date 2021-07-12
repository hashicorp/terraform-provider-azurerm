package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiRelease() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiReleaseCreateUpdate,
		Read:   resourceApiManagementApiReleaseRead,
		Update: resourceApiManagementApiReleaseCreateUpdate,
		Delete: resourceApiManagementApiReleaseDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiReleaseID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementChildName,
			},

			"api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiID,
			},

			"notes": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
func resourceApiManagementApiReleaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	apiId, err := parse.ApiID(d.Get("api_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewApiReleaseID(subscriptionId, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api_release", id.ID())
		}
	}

	parameters := apimanagement.APIReleaseContract{
		APIReleaseContractProperties: &apimanagement.APIReleaseContractProperties{
			APIID: utils.String(d.Get("api_id").(string)),
			Notes: utils.String(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, apiId.ResourceGroup, apiId.ServiceName, apiId.Name, name, parameters, ""); err != nil {
		return fmt.Errorf("creating/ updating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementApiReleaseRead(d, meta)
}

func resourceApiManagementApiReleaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiReleaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.ReleaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] apimanagement %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.ReleaseName)
	if props := resp.APIReleaseContractProperties; props != nil {
		d.Set("api_id", props.APIID)
		d.Set("notes", props.Notes)
	}
	return nil
}

func resourceApiManagementApiReleaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiReleasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiReleaseID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.ReleaseName, "*"); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	return nil
}
