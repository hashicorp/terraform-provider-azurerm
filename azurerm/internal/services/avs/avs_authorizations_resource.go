package avs

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

func resourceArmAvsAuthorization() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAvsAuthorizationCreate,
		Read:   resourceArmAvsAuthorizationRead,
		Delete: resourceArmAvsAuthorizationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AvsAuthorizationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"private_cloud_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"express_route_authorization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"express_route_authorization_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmAvsAuthorizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.AuthorizationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	pcId, _ := parse.AvsPrivateCloudID(d.Get("private_cloud_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, pcId.ResourceGroup, pcId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", name, pcId.ResourceGroup, pcId.Name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_avs_authorization", *existing.ID)
		}
	}

	future, err := client.CreateOrUpdate(ctx, pcId.ResourceGroup, pcId.Name, name, nil)
	if err != nil {
		return fmt.Errorf("creating/updating Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", name, pcId.ResourceGroup, pcId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", name, pcId.ResourceGroup, pcId.Name, err)
	}

	resp, err := client.Get(ctx, pcId.ResourceGroup, pcId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", name, pcId.ResourceGroup, pcId.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Avs Authorization %q (Resource Group %q / privateCloudName %q) ID", name, pcId.ResourceGroup, pcId.Name)
	}

	d.SetId(*resp.ID)
	return resourceArmAvsAuthorizationRead(d, meta)
}

func resourceArmAvsAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.AuthorizationClient
	pcClient := meta.(*clients.Client).Avs.PrivateCloudClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvsAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] avs %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	d.Set("name", id.Name)

	pcResp, err := pcClient.Get(ctx, id.ResourceGroup, id.PrivateCloudName)
	if err != nil {
		return fmt.Errorf("retrieving Avs PrivateCloud %q (Resource Group %q): %+v", id.PrivateCloudName, id.ResourceGroup, err)
	}

	if pcResp.ID == nil || *pcResp.ID == "" {
		return fmt.Errorf("avs PrivateCloud %q (Resource Group %q) ID is empty or nil", id.PrivateCloudName, id.ResourceGroup)
	}

	d.Set("private_cloud_id", pcResp.ID)
	if props := resp.ExpressRouteAuthorizationProperties; props != nil {
		d.Set("express_route_authorization_id", props.ExpressRouteAuthorizationID)
		d.Set("express_route_authorization_key", props.ExpressRouteAuthorizationKey)
	}
	return nil
}

func resourceArmAvsAuthorizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Avs.AuthorizationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AvsAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Avs Authorization %q (Resource Group %q / privateCloudName %q): %+v", id.Name, id.ResourceGroup, id.PrivateCloudName, err)
	}
	return nil
}
