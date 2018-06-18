package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServicePrincipalCreate,
		Read:   resourceArmServicePrincipalRead,
		Delete: resourceArmServicePrincipalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateUUID,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_principal_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmServicePrincipalCreate(d *schema.ResourceData, meta interface{}) error {
	servicePrincipalsClient := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	appId := d.Get("app_id").(string)

	properties := graphrbac.ServicePrincipalCreateParameters{
		AppID:          &appId,
		AccountEnabled: utils.Bool(true),
	}

	spn, err := servicePrincipalsClient.Create(ctx, properties)
	if err != nil {
		return err
	}

	d.SetId(*spn.ObjectID)
	return resourceArmServicePrincipalRead(d, meta)
}

func resourceArmServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Service Principal ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Service Principal %q: %+v", d.Id(), err)
	}

	d.Set("app_id", resp.AppID)
	d.Set("object_id", resp.ObjectID)
	d.Set("display_name", resp.DisplayName)
	d.Set("service_principal_names", resp.ServicePrincipalNames)

	return nil
}

func resourceArmServicePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
