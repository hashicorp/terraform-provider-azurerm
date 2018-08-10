package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var servicePrincipalResourceName = "azurerm_service_principal"

func resourceArmActiveDirectoryServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActiveDirectoryServicePrincipalCreate,
		Read:   resourceArmActiveDirectoryServicePrincipalRead,
		Delete: resourceArmActiveDirectoryServicePrincipalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmActiveDirectoryServicePrincipalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	applicationId := d.Get("application_id").(string)

	apps, err := client.ListComplete(waitCtx, "")
	if err != nil {
		return fmt.Errorf("Error checking for existence of Service Principal %q: %+v", applicationId, err)
	}

	for apps.NotDone() {
		a := apps.Value()
		if a.AppID == nil || a.ObjectID == nil {
			continue
		}

		if *a.AppID == applicationId {
			return tf.ImportAsExistsError("azurerm_azuread_service_principal", *a.ObjectID)
		}

		e := apps.Next()
		if e != nil {
			return e
		}
	}

	properties := graphrbac.ServicePrincipalCreateParameters{
		AppID: utils.String(applicationId),
		// there's no way of retrieving this, and there's no way of changing it
		// given there's no way to change it - we'll just default this to true
		AccountEnabled: utils.Bool(true),
	}

	app, err := client.Create(waitCtx, properties)
	if err != nil {
		return fmt.Errorf("Error creating Service Principal %q: %+v", applicationId, err)
	}

	objectId := *app.ObjectID
	resp, err := client.Get(ctx, objectId)
	if err != nil {
		return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
	}

	d.SetId(*resp.ObjectID)

	return resourceArmActiveDirectoryServicePrincipalRead(d, meta)
}

func resourceArmActiveDirectoryServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	objectId := d.Id()
	app, err := client.Get(ctx, objectId)
	if err != nil {
		if utils.ResponseWasNotFound(app.Response) {
			log.Printf("[DEBUG] Service Principal with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
	}

	d.Set("application_id", app.AppID)
	d.Set("display_name", app.DisplayName)

	return nil
}

func resourceArmActiveDirectoryServicePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	applicationId := d.Id()
	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	app, err := client.Delete(waitCtx, applicationId)
	if err != nil {
		if !response.WasNotFound(app.Response) {
			return fmt.Errorf("Error deleting Service Principal ID %q: %+v", applicationId, err)
		}
	}

	return nil
}
