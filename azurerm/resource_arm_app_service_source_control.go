package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
		Create: resourceArmAppServiceSourceControlCreateUpdate,
		Read:   resourceArmAppServiceSourceControlRead,
		Update: resourceArmAppServiceSourceControlCreateUpdate,
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

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.ScmTypeBitbucketGit),
					// Bitbucket Mercurial support will be removed June 2020
					// https://bitbucket.org/blog/sunsetting-mercurial-support-in-bitbucket
					string(web.ScmTypeBitbucketHg),
					// CodePlex was shut down December 2017
					// https://devblogs.microsoft.com/bharry/shutting-down-codeplex/
					// string(web.ScmTypeCodePlexGit),
					// string(web.ScmTypeCodePlexHg),
					string(web.ScmTypeDropbox),
					string(web.ScmTypeExternalGit),
					string(web.ScmTypeExternalHg),
					string(web.ScmTypeGitHub),
					string(web.ScmTypeLocalGit),
					string(web.ScmTypeOneDrive),
					string(web.ScmTypeTfs),
					string(web.ScmTypeVSO),
					// Not in the specs, but is set by Azure Pipelines
					// https://github.com/Microsoft/azure-pipelines-tasks/blob/master/Tasks/AzureRmWebAppDeploymentV4/operations/AzureAppServiceUtility.ts#L19
					// upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/5345
					"VSTSRM",
				}, false),
			},

			"repo_url": {
				Type:         schema.TypeString,
				Optional:     true,
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
				Default:  false,
			},

			"is_manual_integration": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_mercurial": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"local_repo_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceSourceControlCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
	scmType := d.Get("type").(string)
	repoUrl := d.Get("repo_url").(string)
	branch := d.Get("branch").(string)
	deploymentRollbackEnabled := d.Get("deployment_rollback_enabled").(bool)

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

	if repoUrl != "" {
		sourceControl := web.SiteSourceControl{
			SiteSourceControlProperties: &web.SiteSourceControlProperties{
				RepoURL: utils.String(repoUrl),
				Branch:  utils.String(branch),
				DeploymentRollbackEnabled: utils.Bool(deploymentRollbackEnabled),
			},
		}

		if web.ScmType(scmType) == web.ScmTypeExternalGit || web.ScmType(scmType) == web.ScmTypeExternalHg {
			sourceControl.SiteSourceControlProperties.IsManualIntegration = utils.Bool(true)
		}

		if web.ScmType(scmType) == web.ScmTypeBitbucketHg || web.ScmType(scmType) == web.ScmTypeCodePlexHg || web.ScmType(scmType) == web.ScmTypeExternalHg {
			sourceControl.SiteSourceControlProperties.IsMercurial = utils.Bool(true)
		}

		if _, err := client.CreateOrUpdateSourceControl(ctx, resourceGroup, appServiceName, sourceControl); err != nil {
			return fmt.Errorf("Error updating App Service Source Control (App Service %q / Resource Group %q): %s", appServiceName, resourceGroup, err)
		}
	}

	sitePatchResource := web.SitePatchResource{
		SitePatchResourceProperties: &web.SitePatchResourceProperties{
			SiteConfig: &web.SiteConfig{
				ScmType: web.ScmType(scmType),
			},
		},
	}

	if _, err := client.Update(ctx, resourceGroup, appServiceName, sitePatchResource); err != nil {
		return fmt.Errorf("Error updating App Service Configuration (App Service %q / Resource Group %q): %+v", appServiceName, resourceGroup, err)
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

	siteConfigResp, err := client.GetConfiguration(ctx, resourceGroup, appServiceName)
	if err != nil {
		if !utils.ResponseWasNotFound(siteConfigResp.Response) {
			return fmt.Errorf("Error checking for presence of existing App Service Configuration (App Service %q / Resource Group %q): %s", appServiceName, resourceGroup, err)
		}
	}

	siteConfigProps := siteConfigResp.SiteConfig
	if siteConfigProps != nil {
		d.Set("type", siteConfigProps.ScmType)
	}

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
		if web.ScmType(siteConfigProps.ScmType) == web.ScmTypeLocalGit {
			d.Set("local_repo_url", props.RepoURL)
		} else {
			d.Set("repo_url", props.RepoURL)
		}
		d.Set("branch", props.Branch)
		d.Set("deployment_rollback_enabled", props.DeploymentRollbackEnabled)
		d.Set("is_manual_integration", props.IsManualIntegration)
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

	sitePatchResource := web.SitePatchResource{
		SitePatchResourceProperties: &web.SitePatchResourceProperties{
			SiteConfig: &web.SiteConfig{
				ScmType: web.ScmType(web.ScmTypeNone),
			},
		},
	}

	if _, err := client.Update(ctx, resourceGroup, appServiceName, sitePatchResource); err != nil {
		return fmt.Errorf("Error updating App Service Configuration (App Service %q / Resource Group %q): %+v", appServiceName, resourceGroup, err)
	}

	return nil
}
