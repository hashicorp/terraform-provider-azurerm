package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
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

			"manual_integration_enabled": {
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
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for App Service Source Control creation.")

	appServiceId := d.Get("app_service_id").(string)

	id, err := azure.ParseAzureResourceID(appServiceId)
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]

	repoUrl := d.Get("repo_url").(string)
	branch := d.Get("branch").(string)
	deploymentRollbackEnabled := d.Get("deployment_rollback_enabled").(bool)
	manualIntegrationEnabled := d.Get("manual_integration_enabled").(bool)

	locks.ByName(appServiceName, appServiceSourceControlResourceName)
	defer locks.UnlockByName(appServiceName, appServiceSourceControlResourceName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetSourceControl(ctx, resourceGroup, appServiceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Source Control (App Service %q / Resource Group %q): %s", appServiceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError(appServiceSourceControlResourceName, *existing.ID)
		}
	}

	sourceControl := web.SiteSourceControl{
		SiteSourceControlProperties: &web.SiteSourceControlProperties{
			RepoURL:                   utils.String(repoUrl),
			Branch:                    utils.String(branch),
			DeploymentRollbackEnabled: utils.Bool(deploymentRollbackEnabled),
			IsManualIntegration:       utils.Bool(manualIntegrationEnabled),
		},
	}

	if _, err := client.CreateOrUpdateSourceControl(ctx, resourceGroup, appServiceName, sourceControl); err != nil {
		return fmt.Errorf("Error updating App Service Source Control (App Service %q / Resource Group %q): %s", appServiceName, resourceGroup, err)
	}

	read, err := client.GetSourceControl(ctx, resourceGroup, appServiceName)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Source Control (App Service %q / Resource Group %q): %s", appServiceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Source Control (App Service %q / Resource Group %q) ID", appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceSourceControlRead(d, meta)
}

func resourceArmAppServiceSourceControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]

	resp, err := client.GetSourceControl(ctx, resourceGroup, appServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Source Control (App Service %q / Resource Group %q) was not found - removing from state", appServiceName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Source Control (App Service %q / Resource Group %q): %+v", appServiceName, resourceGroup, err)
	}

	if props := resp.SiteSourceControlProperties; props != nil {
		d.Set("repo_url", props.RepoURL)
		d.Set("branch", props.Branch)
		d.Set("deployment_rollback_enabled", props.DeploymentRollbackEnabled)
		d.Set("manual_integration_enabled", props.IsManualIntegration)
		d.Set("is_mercurial", props.IsMercurial)
	}

	return nil
}

func resourceArmAppServiceSourceControlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.AppServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]

	locks.ByName(appServiceName, appServiceSourceControlResourceName)
	defer locks.UnlockByName(appServiceName, appServiceSourceControlResourceName)

	log.Printf("[DEBUG] Deleting App Service Source Control (App Service %q / Resource Group %q)", appServiceName, resourceGroup)

	resp, err := client.DeleteSourceControl(ctx, resourceGroup, appServiceName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Source Control (App Service %q / Resource Group %q): %s)", appServiceName, resourceGroup, err)
		}
	}

	return nil
}
