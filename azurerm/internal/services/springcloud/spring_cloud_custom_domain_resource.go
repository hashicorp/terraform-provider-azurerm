package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSpringCloudCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudCustomDomainCreateUpdate,
		Read:   resourceSpringCloudCustomDomainRead,
		Update: resourceSpringCloudCustomDomainCreateUpdate,
		Delete: resourceSpringCloudCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudCustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudCustomDomainName,
			},

			"spring_cloud_app_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"certificate_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"thumbprint"},
			},

			"thumbprint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"certificate_name"},
			},
		},
	}
}

func resourceSpringCloudCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CustomDomainsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	appId, err := parse.SpringCloudAppID(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	resourceId := parse.NewSpringCloudCustomDomainID(appId.SubscriptionId, appId.ResourceGroup, appId.SpringName, appId.AppName, name).ID()
	if d.IsNewResource() {
		existing, err := client.Get(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("making Read request on AzureRM Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_custom_domain", resourceId)
		}
	}

	domain := appplatform.CustomDomainResource{
		Properties: &appplatform.CustomDomainProperties{
			Thumbprint: utils.String(d.Get("thumbprint").(string)),
			CertName:   utils.String(d.Get("certificate_name").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, appId.ResourceGroup, appId.SpringName, appId.AppName, name, domain); err != nil {
		return fmt.Errorf("creating/update Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", name, appId.SpringName, appId.AppName, appId.ResourceGroup, err)
	}
	d.SetId(resourceId)
	return resourceSpringCloudCustomDomainRead(d, meta)
}

func resourceSpringCloudCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Custom Domain %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", id.DomainName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	d.Set("name", id.DomainName)
	d.Set("spring_cloud_app_id", parse.NewSpringCloudAppID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("certificate_name", props.CertName)
		d.Set("thumbprint", props.Thumbprint)
	}

	return nil
}

func resourceSpringCloudCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.CustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DomainName); err != nil {
		return fmt.Errorf("deleting Spring Cloud Custom Domain %q (Spring Cloud service %q / App %q / rcsource group %q): %+v", id.DomainName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}
	return nil
}
