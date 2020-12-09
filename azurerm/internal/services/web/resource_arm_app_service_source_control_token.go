package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var appServiceSourceControlTokenResourceName = "azurerm_app_service_source_control_token"

func resourceArmAppServiceSourceControlToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSourceControlTokenCreateUpdate,
		Read:   resourceArmAppServiceSourceControlTokenRead,
		Update: resourceArmAppServiceSourceControlTokenCreateUpdate,
		Delete: resourceArmAppServiceSourceControlTokenDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := validate.SourceControlTokenName()(id, "id")
			if len(err) > 0 {
				return fmt.Errorf("%s", err)
			}

			return nil
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SourceControlTokenName(),
			},

			"token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"token_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmAppServiceSourceControlTokenCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.BaseClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Source Control Token creation.")

	scmType := d.Get("type").(string)
	token := d.Get("token").(string)
	tokenSecret := d.Get("token_secret").(string)

	locks.ByName(scmType, appServiceSourceControlTokenResourceName)
	defer locks.UnlockByName(scmType, appServiceSourceControlTokenResourceName)

	properties := web.SourceControl{
		SourceControlProperties: &web.SourceControlProperties{
			Token:       utils.String(token),
			TokenSecret: utils.String(tokenSecret),
		},
	}

	if _, err := client.UpdateSourceControl(ctx, scmType, properties); err != nil {
		return fmt.Errorf("Error updating App Service Source Control Token (Type %q): %s", scmType, err)
	}

	read, err := client.GetSourceControl(ctx, scmType)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Source Control Token (Type %q): %s", scmType, err)
	}
	if read.Name == nil {
		return fmt.Errorf("Cannot read App Service Source Control Token (Type %q)", scmType)
	}

	d.SetId(*read.Name)

	return resourceArmAppServiceSourceControlTokenRead(d, meta)
}

func resourceArmAppServiceSourceControlTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.BaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	scmType := d.Id()

	resp, err := client.GetSourceControl(ctx, scmType)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Source Control Token (Type %q) was not found - removing from state", scmType)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Source Control Token (Type %q): %+v", scmType, err)
	}

	d.Set("type", resp.Name)

	if props := resp.SourceControlProperties; props != nil {
		d.Set("token", props.Token)
		d.Set("token_secret", props.TokenSecret)
	}

	return nil
}

func resourceArmAppServiceSourceControlTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.BaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	scmType := d.Id()

	// Delete cleans up existing tokens by setting their values to ""
	token := ""
	tokenSecret := ""

	locks.ByName(scmType, appServiceSourceControlTokenResourceName)
	defer locks.UnlockByName(scmType, appServiceSourceControlTokenResourceName)

	log.Printf("[DEBUG] Deleting App Service Source Control Token (Type %q)", scmType)

	properties := web.SourceControl{
		SourceControlProperties: &web.SourceControlProperties{
			Token:       utils.String(token),
			TokenSecret: utils.String(tokenSecret),
		},
	}

	if _, err := client.UpdateSourceControl(ctx, scmType, properties); err != nil {
		return fmt.Errorf("Error deleting App Service Source Control Token (Type %q): %s", scmType, err)
	}

	return nil
}
