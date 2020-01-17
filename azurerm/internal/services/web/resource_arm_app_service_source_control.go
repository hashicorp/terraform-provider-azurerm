package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var appServiceSourceControlResourceName = "azurerm_app_service_source_control"

func resourceArmAppServiceSourceControl() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceSourceControlCreate,
		Read:   resourceArmAppServiceSourceControlRead,
		Delete: resourceArmAppServiceSourceControlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"repo_url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"branch": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"deployment_rollback_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"is_manual_integration": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"is_mercurial": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceSourceControlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Source Control creation.")

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	repoUrl := d.Get("repo_url").(string)
	branch := d.Get("branch").(string)
	deploymentRollbackEnabled := d.Get("deployment_rollback_enabled").(bool)
	isManualIntegration := d.Get("is_manual_integration").(bool)

	locks.ByName(name, appServiceSourceControlResourceName)
	defer locks.UnlockByName(name, appServiceSourceControlResourceName)

	siteSourceControl := web.SiteSourceControl{
		SiteSourceControlProperties: &web.SiteSourceControlProperties{
			RepoURL: utils.String(repoUrl),
			Branch:  utils.String(branch),
			DeploymentRollbackEnabled: utils.Bool(deploymentRollbackEnabled),
			IsManualIntegration:       utils.Bool(isManualIntegration),
		},
	}

	if _, err := client.CreateOrUpdateSourceControl(ctx, resourceGroup, name, siteSourceControl); err != nil {
		return fmt.Errorf("Error updating App Service Source Control (App Service %q / Resource Group %q): %s", name, resourceGroup, err)
	}

	read, err := client.GetSourceControl(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Source Control (App Service %q / Resource Group %q): %s", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Source Control (App Service %q / Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSourceControlRead(d, meta)
}

func resourceArmAppServiceSourceControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	resp, err := client.GetSourceControl(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Source Control (App Service %q / Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Source Control (App Service %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	if props := resp.SiteSourceControlProperties; props != nil {
		d.Set("repo_url", props.RepoURL)
		d.Set("branch", props.Branch)
		d.Set("deployment_rollback_enabled", props.DeploymentRollbackEnabled)
		d.Set("is_manual_integration", props.IsManualIntegration)
		d.Set("is_mercurial", props.IsMercurial)
	}

	return nil
}

func resourceArmAppServiceSourceControlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Get("app_service_id").(string))
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["sites"]

	locks.ByName(name, appServiceSourceControlResourceName)
	defer locks.UnlockByName(name, appServiceSourceControlResourceName)

	log.Printf("[DEBUG] Deleting App Service Source Control (App Service %q / Resource Group %q)", name, resourceGroup)

	resp, err := client.DeleteSourceControl(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Source Control (App Service %q / Resource Group %q): %s)", name, resourceGroup, err)
		}
	}

	return nil
}
