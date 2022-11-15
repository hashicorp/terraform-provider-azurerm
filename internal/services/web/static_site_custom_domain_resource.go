package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"domain_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"validation_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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

	id := parse.NewStaticSiteCustomDomainID(staticSiteId.SubscriptionId, staticSiteId.ResourceGroup, staticSiteId.Name, d.Get("domain_name").(string))
	_, err = client.GetStaticSite(ctx, id.ResourceGroup, id.StaticSiteName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *staticSiteId, err)
	}

	existing, err := client.GetStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, id.StaticSiteName, id.CustomDomainName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_static_site_custom_domain", staticSiteId.ID())
	}

	validationMethod := d.Get("validation_type").(string)
	if validationMethod == "" {
		return fmt.Errorf("`validation_type` can't be empty string")
	}

	siteEnvelope := web.StaticSiteCustomDomainRequestPropertiesARMResource{
		StaticSiteCustomDomainRequestPropertiesARMResourceProperties: &web.StaticSiteCustomDomainRequestPropertiesARMResourceProperties{
			ValidationMethod: &validationMethod,
		},
	}

	future, err := client.CreateOrUpdateStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, id.StaticSiteName, id.CustomDomainName, siteEnvelope)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// we can't wait for the future to be complete for txt validation as we need to give the user the validation token
	if validationMethod == cnameValidationType {
		err := future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("waiting for creation of %s: %+v", id, err)
		}
	}

	// we need to poll till the validation token is ready
	if validationMethod == txtValidationType {
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("context was missing a deadline")
		}
		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				string(web.CustomDomainStatusRetrievingValidationToken),
			},
			Target: []string{
				string(web.CustomDomainStatusValidating),
			},
			MinTimeout: 20 * time.Second,
			Timeout:    time.Until(deadline),
			Refresh: func() (interface{}, string, error) {
				domain, err := client.GetStaticSiteCustomDomain(ctx, staticSiteId.ResourceGroup, id.StaticSiteName, id.CustomDomainName)
				if err != nil {
					return domain, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
				}

				if domain.StaticSiteCustomDomainOverviewARMResourceProperties == nil {
					return nil, "Failed", fmt.Errorf("`properties` was missing from the response")
				}
				return domain, string(domain.StaticSiteCustomDomainOverviewARMResourceProperties.Status), nil
			},
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for DNS Validation after Creation of %s %+v", id, err)
		}
	}

	d.SetId(id.ID())

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
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("domain_name", id.CustomDomainName)
	d.Set("static_site_id", parse.NewStaticSiteID(id.SubscriptionId, id.ResourceGroup, id.StaticSiteName).ID())

	if props := resp.StaticSiteCustomDomainOverviewARMResourceProperties; props != nil {
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
