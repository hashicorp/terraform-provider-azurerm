package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var appServiceSourceControlTokenResourceName = "azurerm_app_service_source_control_token"

func resourceAppServiceSourceControlToken() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSourceControlTokenCreateUpdate,
		Read:   resourceAppServiceSourceControlTokenRead,
		Update: resourceAppServiceSourceControlTokenCreateUpdate,
		Delete: resourceAppServiceSourceControlTokenDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := validate.SourceControlTokenName()(id, "id")
			if len(err) > 0 {
				return fmt.Errorf("%s", err)
			}

			return nil
		}),

		DeprecationMessage: "The `azurerm_app_service_source_control_token` resource has been superseded by the `azurerm_source_control_token` resource. Whilst this resource will continue to be available in the 2.x and 3.x releases it is feature-frozen for compatibility purposes, will no longer receive any updates and will be removed in a future major release of the Azure Provider.",

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SourceControlTokenName(),
			},

			"token": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"token_secret": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceAppServiceSourceControlTokenCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.BaseClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Source Control Token creation.")

	token := d.Get("token").(string)
	tokenSecret := d.Get("token_secret").(string)
	id := parse.NewAppServiceSourceControlTokenID(d.Get("type").(string))

	locks.ByName(id.Type, appServiceSourceControlTokenResourceName)
	defer locks.UnlockByName(id.Type, appServiceSourceControlTokenResourceName)

	properties := web.SourceControl{
		SourceControlProperties: &web.SourceControlProperties{
			Token:       utils.String(token),
			TokenSecret: utils.String(tokenSecret),
		},
	}

	if _, err := client.UpdateSourceControl(ctx, id.Type, properties); err != nil {
		return fmt.Errorf("updating %s: %s", id, err)
	}

	d.SetId(id.Type)

	return resourceAppServiceSourceControlTokenRead(d, meta)
}

func resourceAppServiceSourceControlTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("making Read request on App Service Source Control Token (Type %q): %+v", scmType, err)
	}

	d.Set("type", resp.Name)

	if props := resp.SourceControlProperties; props != nil {
		d.Set("token", props.Token)
		d.Set("token_secret", props.TokenSecret)
	}

	return nil
}

func resourceAppServiceSourceControlTokenDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("deleting App Service Source Control Token (Type %q): %s", scmType, err)
	}

	return nil
}
