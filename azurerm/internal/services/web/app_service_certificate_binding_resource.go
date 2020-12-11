package web

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var appServiceHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceArmAppServiceCertificateBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCertificateBindingCreate,
		Read:   resourceArmAppServiceCertificateBindingRead,
		Delete: resourceArmAppServiceCertificateBindingDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CertificateBindingID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"hostname_binding_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.HostnameBindingID,
			},

			"certificate_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CertificateID,
			},

			"ssl_state": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.SslStateIPBasedEnabled),
					string(web.SslStateSniEnabled),
				}, false),
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmAppServiceCertificateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	certClient := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	hostnameBindingID, err := parse.HostnameBindingID(d.Get("hostname_binding_id").(string))
	if err != nil {
		return err
	}

	certificateID, err := parse.CertificateID(d.Get("certificate_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewCertificateBindingId(*hostnameBindingID, *certificateID)
	if err != nil {
		return fmt.Errorf("could not parse ID: %+v", err)
	}

	certDetails, err := certClient.Get(ctx, id.CertificateId.ResourceGroup, id.CertificateId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(certDetails.Response) {
			return fmt.Errorf("retrieving App Service Certificate %q (Resource Group %q), not found", id.CertificateId.Name, id.CertificateId.ResourceGroup)
		}
		return fmt.Errorf("failed reading App Service Certificate %q (Resource Group %q): %+v", id.CertificateId.Name, id.CertificateId.ResourceGroup, err)
	}

	if certDetails.Thumbprint == nil {
		return fmt.Errorf("could not read thumbprint from certificate %q (resource group %q): %+v", id.CertificateId.Name, id.CertificateId.ResourceGroup, err)
	}
	thumbprint := certDetails.Thumbprint

	binding, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(binding.Response) {
			return fmt.Errorf("retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup, err)
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.CertificateId.Name, id.HostnameBindingId.ResourceGroup, err)
	}

	props := binding.HostNameBindingProperties
	if props != nil {
		if props.Thumbprint != nil && *props.Thumbprint == *thumbprint {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_binding", id.ID())
		}
	}

	locks.ByName(id.HostnameBindingId.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostnameBindingId.SiteName, appServiceHostnameBindingResourceName)

	binding.HostNameBindingProperties.SslState = web.SslState(d.Get("ssl_state").(string))
	binding.HostNameBindingProperties.Thumbprint = thumbprint

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name, binding); err != nil {
		return fmt.Errorf("creating/updating Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.CertificateId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceArmAppServiceCertificateBindingRead(d, meta)
}

func resourceArmAppServiceCertificateBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	props := resp.HostNameBindingProperties
	if props == nil || props.Thumbprint == nil {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("hostname_binding_id", id.HostnameBindingId.ID())
	d.Set("certificate_id", id.CertificateId.ID())
	d.Set("ssl_state", string(props.SslState))
	d.Set("thumbprint", props.Thumbprint)
	d.Set("hostname", id.HostnameBindingId.Name)
	d.Set("app_service_name", id.HostnameBindingId.SiteName)

	return nil
}

func resourceArmAppServiceCertificateBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateBindingID(d.Id())
	if err != nil {
		return err
	}

	binding, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(binding.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	props := binding.HostNameBindingProperties
	if props == nil || props.Thumbprint == nil {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup)
		d.SetId("")
		return nil
	}

	locks.ByName(id.HostnameBindingId.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostnameBindingId.SiteName, appServiceHostnameBindingResourceName)

	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q / Resource Group %q)", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup)

	binding.HostNameBindingProperties.SslState = web.SslStateDisabled
	binding.HostNameBindingProperties.Thumbprint = nil

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.HostnameBindingId.SiteName, id.HostnameBindingId.Name, binding); err != nil {
		return fmt.Errorf("deleting Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.HostnameBindingId.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	return nil
}
