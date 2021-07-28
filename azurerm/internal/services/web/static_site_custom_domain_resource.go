package web

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	txtValidationType   = "dns-txt-token"
	cnameValidationType = "cname-delegation"
)

func resourceStaticSiteCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStaticSiteCustomDomainCreateOrUpdate,
		Read:   resourceStaticSiteCustomDomainRead,
		Delete: resourceStaticSiteCustomDomainDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StaticSiteCustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"static_site_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"validation_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					txtValidationType,
					cnameValidationType,
				}, false),
			},

			"validation_token": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceStaticSiteCustomDomainCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Static Site custom domain creation.")

	staticSiteId, err := parse.StaticSiteID(d.Get("static_site_id").(string))

	if err != nil {
		return err
	}

	domainName := d.Get("name").(string)

	resourceGroup := staticSiteId.ResourceGroup
	staticSiteName := staticSiteId.Name

	_, err = client.GetStaticSite(ctx, resourceGroup, staticSiteName)
	if err != nil {
		return fmt.Errorf("finding Static Site (Name %q / Resource Group %q): %s", staticSiteName, resourceGroup, err)
	}

	if d.IsNewResource() {
		existing, err := client.GetStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, staticSiteName, domainName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed checking for presence of existing %s: %+v", staticSiteId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_static_site_custom_domain", staticSiteId.ID())
		}
	}

	validationMethod := d.Get("validation_type").(string)

	siteEnvelope := web.StaticSiteCustomDomainRequestPropertiesARMResource{
		StaticSiteCustomDomainRequestPropertiesARMResourceProperties: &web.StaticSiteCustomDomainRequestPropertiesARMResourceProperties{
			ValidationMethod: &validationMethod,
		},
	}

	future, err := client.CreateOrUpdateStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, staticSiteName, domainName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("failed creating %s: %+v", staticSiteId, err)
	}

	// we can't wait for the future to be complete for txt validation as we need to give the user the validation token
	if validationMethod == cnameValidationType {
		err := future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("error waiting for creation of Static Site Custom Domain %q (Static Site %q / Resource Group %q): %+v", domainName, staticSiteName, staticSiteId.ResourceGroup, err)
		}
	}

	// we need to poll till the validation token is ready
	if validationMethod == txtValidationType {
		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				string(web.CustomDomainStatusRetrievingValidationToken),
			},
			Target: []string{
				string(web.CustomDomainStatusValidating),
			},
			MinTimeout: 30 * time.Second,
			Timeout:    10 * time.Minute,
			Refresh: func() (interface{}, string, error) {
				domain, err := client.GetStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, staticSiteName, domainName)
				if err != nil {
					return domain, "Error", fmt.Errorf("error retrieving Static Site custom domain %q (Static Site %q / Resource Group %q): %+v", domainName, staticSiteName, staticSiteId.ResourceGroup, err)
				}

				status := domain.StaticSiteCustomDomainOverviewARMResourceProperties.Status

				return domain, string(status), nil
			},
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("error waiting for Static Site custom domain %q (Static Site %q / Resource Group %q): %+v", domainName, staticSiteName, staticSiteId.ResourceGroup, err)
		}
	}

	resourceId := parse.NewStaticSiteCustomDomainID(staticSiteId.SubscriptionId, staticSiteId.ResourceGroup, staticSiteName, domainName)
	d.SetId(resourceId.ID())

	return resourceStaticSiteCustomDomainRead(d, meta)
}

func resourceStaticSiteCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetStaticSiteCustomDomain(ctx, id.ResourceGroup, id.StaticSiteName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Static Site Custom Domain %q (resource group %q) was not found - removing from state", id.CustomDomainName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed making Read request on %s: %+v", id, err)
	}
	d.Set("name", id.CustomDomainName)
	d.Set("static_site_id", parse.NewStaticSiteID(id.SubscriptionId, id.ResourceGroup, id.StaticSiteName).ID())

	if resp.ValidationToken != nil {
		d.Set("validation_token", resp.ValidationToken)
	}

	return nil
}

func resourceStaticSiteCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Static Site Custom Domain %q (resource group %q)", id.CustomDomainName, id.ResourceGroup)

	future, err := client.DeleteStaticSiteCustomDomain(ctx, id.ResourceGroup, id.StaticSiteName, id.CustomDomainName)

	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
