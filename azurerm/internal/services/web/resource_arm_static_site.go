package web

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
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

			"sku_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Free",
				ValidateFunc: validation.StringInSlice([]string{"Free"}, false),
			},

			"sku_size": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Free",
				ValidateFunc: validation.StringInSlice([]string{"Free"}, false),
			},

			"default_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
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

	skuName := d.Get("sku_size").(string)
	skuTier := d.Get("sku_tier").(string)

	staticSiteSkuDescription := &web.SkuDescription{Name: &skuName, Tier: &skuTier}

	siteEnvelope := web.StaticSiteARMResource{
		Sku:        staticSiteSkuDescription,
		StaticSite: &web.StaticSite{},
		Location:   &location,
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

	if prop := resp.StaticSite; prop != nil {
		defaultHostname := ""
		if prop.DefaultHostname != nil {
			defaultHostname = *prop.DefaultHostname
		}
		d.Set("default_host_name", defaultHostname)
	}
	if sku := resp.Sku; sku != nil {
		skuName := ""
		if v := sku.Name; v != nil {
			skuName = *v
		}
		d.Set("sku_size", skuName)

		skuTier := ""
		if v := sku.Tier; v != nil {
			skuTier = *v
		}
		d.Set("sku_tier", skuTier)
	}

	secretResp, err := client.ListStaticSiteSecrets(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("listing secretes for %s: %v", id, err)
	}

	apiKey := ""
	if pkey := secretResp.Properties["apiKey"]; pkey != nil {
		apiKey = *pkey
	}
	d.Set("api_key", apiKey)

	return nil
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
