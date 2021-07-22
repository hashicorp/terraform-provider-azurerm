package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVmwareExpressRouteAuthorization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwareExpressRouteAuthorizationCreate,
		Read:   resourceVmwareExpressRouteAuthorizationRead,
		Delete: resourceVmwareExpressRouteAuthorizationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRouteAuthorizationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_cloud_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateCloudID,
			},

			"express_route_authorization_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"express_route_authorization_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
func resourceVmwareExpressRouteAuthorizationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	privateCloudId, err := parse.PrivateCloudID(d.Get("private_cloud_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewExpressRouteAuthorizationID(subscriptionId, privateCloudId.ResourceGroup, privateCloudId.Name, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.AuthorizationName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing %q : %+v", id.ID(), err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_vmware_express_route_authorization", id.ID())
	}

	props := avs.ExpressRouteAuthorization{}

	future, err := client.CreateOrUpdate(ctx, privateCloudId.ResourceGroup, privateCloudId.Name, name, props)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of the %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwareExpressRouteAuthorizationRead(d, meta)
}

func resourceVmwareExpressRouteAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.AuthorizationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] avs %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("name", id.AuthorizationName)
	d.Set("private_cloud_id", parse.NewPrivateCloudID(subscriptionId, id.ResourceGroup, id.PrivateCloudName).ID())
	if props := resp.ExpressRouteAuthorizationProperties; props != nil {
		d.Set("express_route_authorization_id", props.ExpressRouteAuthorizationID)
		d.Set("express_route_authorization_key", props.ExpressRouteAuthorizationKey)
	}
	return nil
}

func resourceVmwareExpressRouteAuthorizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.AuthorizationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.PrivateCloudName, id.AuthorizationName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of the %q: %+v", id, err)
	}
	return nil
}
