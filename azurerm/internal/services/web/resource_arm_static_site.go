package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStaticSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStaticSiteCreateOrUpdate,
		Read:   resourceArmStaticSiteRead,
		Update: resourceArmStaticSiteCreateOrUpdate,
		Delete: resourceArmStaticSiteDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StaticSiteName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"github_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"github_repo_url": {
				Type:     schema.TypeString,
				Required: true,
			},

			"branch": {
				Type:     schema.TypeString,
				Required: true,
			},

			"app_directory": {
				Type:     schema.TypeString,
				Required: true,
			},

			"api_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"artifact_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"sku": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tier": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Free",
						},
						"size": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Free",
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmStaticSiteCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Static Site creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetStaticSite(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed checking for presence of existing Static Site %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_static_site", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	nameSku := d.Get("sku.0.name").(string)
	tierSku := d.Get("sku.0.tier").(string)

	staticSiteSkuDescription := &web.SkuDescription{Name: &nameSku, Tier: &tierSku}

	githubRepoURL := d.Get("github_repo_url").(string)
	branch := d.Get("branch").(string)
	githubToken := d.Get("github_token").(string)

	appDirectory := d.Get("app_directory").(string)
	apiDirectory := d.Get("api_directory").(string)
	artifactDirectory := d.Get("artifact_directory").(string)

	staticSiteSourceControl := &web.StaticSite{
		RepositoryURL:   &githubRepoURL,
		Branch:          &branch,
		RepositoryToken: &githubToken,
		BuildProperties: &web.StaticSiteBuildProperties{
			AppLocation:         &appDirectory,
			APILocation:         &apiDirectory,
			AppArtifactLocation: &artifactDirectory,
		},
	}

	siteEnvelope := web.StaticSiteARMResource{
		Sku:        staticSiteSkuDescription,
		StaticSite: staticSiteSourceControl,
		Location:   &location,
		Tags:       tags.Expand(t),
	}

	_, err := client.CreateOrUpdateStaticSite(ctx, resourceGroup, name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("failed creating Static Site %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	read, err := client.GetStaticSite(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("failed retrieving Static Site %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read Static Site %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmStaticSiteRead(d, meta)
}

func resourceArmStaticSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Static Site %q (resource group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed making Read request on AzureRM Static Site %q: %+v", id.Name, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	sc := flattenStaticSiteSourceControl(resp.StaticSite, d)
	if err := d.Set("github_configuration", sc); err != nil {
		return fmt.Errorf("failed setting `github_configuration`: %s", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmStaticSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Static Site %q (resource group %q)", id.Name, id.ResourceGroup)

	resp, err := client.DeleteStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}

func flattenStaticSiteSourceControl(input *web.StaticSite, d *schema.ResourceData) []interface{} {
	if input == nil {
		log.Printf("[DEBUG] SiteSourceControlProperties is nil")
		return []interface{}{}
	}

	githubRepoURL := ""
	if input.RepositoryURL != nil {
		githubRepoURL = *input.RepositoryURL
	}
	branch := ""
	if input.Branch != nil && *input.Branch != "" {
		branch = *input.Branch
	}

	githubToken := ""
	apiDirectory := ""
	appDirectory := ""
	artifactDirectory := ""
	if sc, ok := d.GetOk("github_configuration"); ok {
		var val []interface{}

		if v, ok := sc.([]interface{}); ok {
			val = v
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			githubToken = raw["github_token"].(string)
			apiDirectory = raw["api_directory"].(string)
			appDirectory = raw["app_directory"].(string)
			artifactDirectory = raw["artifact_directory"].(string)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"github_repo_url":    githubRepoURL,
			"branch":             branch,
			"github_token":       githubToken,
			"api_directory":      apiDirectory,
			"artifact_directory": artifactDirectory,
			"app_directory":      appDirectory,
		},
	}
}
