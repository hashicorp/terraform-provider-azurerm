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

			"tags": tags.Schema(),

			"github_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo_token": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"repo_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"app_location": {
							Type:     schema.TypeString,
							Required: true,
						},
						"api_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"artifact_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmStaticSiteCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Static Site creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetStaticSite(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Static Site %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_static_site", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	nameSku := "Free"
	tierSku := "Free"
	staticSiteSkuDescription := &web.SkuDescription{Name: &nameSku, Tier: &tierSku}

	staticSiteType := "Microsoft.Web/staticSites"

	staticSiteSourceControlRaw := d.Get("github_configuration").([]interface{})
	staticSiteSourceControl := expandStaticSiteSourceControl(staticSiteSourceControlRaw)

	siteEnvelope := web.StaticSiteARMResource{
		Sku:        staticSiteSkuDescription,
		Type:       &staticSiteType,
		StaticSite: staticSiteSourceControl,
		Location:   &location,
		Tags:       tags.Expand(t),
	}

	_, err := client.CreateOrUpdateStaticSite(ctx, resGroup, name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("Error creating Static Site %q (Resource Group %q): %s", name, resGroup, err)
	}

	read, err := client.GetStaticSite(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Static Site %q (Resource Group %q): %s", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Static Site %q (resource group %q) ID", name, resGroup)
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
		return fmt.Errorf("Error making Read request on AzureRM Static Site %q: %+v", id.Name, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	sc := flattenStaticSiteSourceControl(resp.StaticSite, d)
	if err := d.Set("github_configuration", sc); err != nil {
		return fmt.Errorf("Error setting `github_configuration`: %s", err)
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

	repoURL := ""
	if input.RepositoryURL != nil {
		repoURL = *input.RepositoryURL
	}
	branch := ""
	if input.Branch != nil && *input.Branch != "" {
		branch = *input.Branch
	}

	repoToken := ""
	apiLocation := ""
	appLocation := ""
	artifactLocation := ""
	if sc, ok := d.GetOk("github_configuration"); ok {
		var val []interface{}

		if v, ok := sc.([]interface{}); ok {
			val = v
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			repoToken = raw["repo_token"].(string)
			apiLocation = raw["api_location"].(string)
			appLocation = raw["app_location"].(string)
			artifactLocation = raw["artifact_location"].(string)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"repo_url":          repoURL,
			"branch":            branch,
			"repo_token":        repoToken,
			"api_location":      apiLocation,
			"artifact_location": artifactLocation,
			"app_location":      appLocation,
		},
	}
}

func expandStaticSiteSourceControl(input []interface{}) *web.StaticSite {
	if len(input) == 0 {
		return nil
	}
	sourceControl := input[0].(map[string]interface{})
	repoURL := sourceControl["repo_url"].(string)
	branch := sourceControl["branch"].(string)
	repoToken := sourceControl["repo_token"].(string)

	appLocation := sourceControl["app_location"].(string)
	apiLocation := ""
	if v, ok := sourceControl["api_location"]; ok {
		apiLocation = v.(string)
	}
	artifactLocation := ""
	if v, ok := sourceControl["artifact_location"]; ok {
		artifactLocation = v.(string)
	}

	staticSite := &web.StaticSite{
		RepositoryURL:   &repoURL,
		Branch:          &branch,
		RepositoryToken: &repoToken,
		BuildProperties: &web.StaticSiteBuildProperties{
			AppLocation:         &appLocation,
			APILocation:         &apiLocation,
			AppArtifactLocation: &artifactLocation,
		},
	}

	return staticSite
}
