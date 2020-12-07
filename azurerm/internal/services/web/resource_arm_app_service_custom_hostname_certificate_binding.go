package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var appServiceCustomHostnameCertificateBindingResourceName = "azurerm_app_service_custom_hostname_certificate_binding"

func resourceArmAppServiceCustomHostnameCertificateBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCustomHostnameCertificateBindingCreate,
		Read:   resourceArmAppServiceCustomHostnameCertificateBindingRead,
		Delete: resourceArmAppServiceCustomHostnameCertificateBindingDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceCustomHostnameBindingID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"hostname_binding_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"resource_group_name": azure.SchemaResourceGroupNameOptionalComputed(),

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

func resourceArmAppServiceCustomHostnameCertificateBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	certClient := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	hostnameBindingID := d.Get("hostname_binding_id").(string)
	certificateID := d.Get("certificate_id").(string)

	certID, err := parse.CertificateID(certificateID)
	if err != nil {
		return err
	}
	certDetails, err := certClient.Get(ctx, certID.ResourceGroup, certID.Name)
	if err != nil {
		return fmt.Errorf("App Service Certificate %q (Resource Group %q) does not exist", certID.Name, certID.ResourceGroup)
	}

	hostname := certID.Name

	id, err := parse.AppServiceCustomHostnameBindingID(hostnameBindingID)
	if err != nil {
		return err
	}
	resp, err := client.GetHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, hostname)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cannot bind certificate as hostname binding doesn't exist. Error retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", hostname, id.AppServiceName, id.ResourceGroup, err)
		}
		return fmt.Errorf("Error retrieving Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", id.Name, hostname, id.AppServiceName, id.ResourceGroup, err)
	}

	thumbprint := certDetails.Thumbprint

	appServiceName := id.AppServiceName
	resourceGroup := id.ResourceGroup

	locks.ByName(appServiceName, appServiceCustomHostnameCertificateBindingResourceName)
	defer locks.UnlockByName(appServiceName, appServiceCustomHostnameCertificateBindingResourceName)

	if d.IsNewResource() {
		props := resp.HostNameBindingProperties
		existing := false
		if (props != nil) && (props.Thumbprint != nil) {
			existing = true
		}

		if existing {
			return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_certificate_binding", hostnameBindingID)
		}
	}

	sslState := d.Get("ssl_state").(string)

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SslState:   web.SslState(sslState),
			Thumbprint: thumbprint,
		},
	}

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, resourceGroup, appServiceName, hostname, properties); err != nil {
		return fmt.Errorf("Error creating/updating Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", hostname, hostname, appServiceName, resourceGroup, err)
	}

	read, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
	if err != nil {
		return fmt.Errorf("Error retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", hostname, appServiceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Hostname Binding %q (App Service %q / Resource Group %q) ID", hostname, appServiceName, resourceGroup)
	}

	d.Set("thumbprint", thumbprint)
	d.Set("hostname", hostname)
	d.Set("app_service_name", appServiceName)
	d.Set("resource_group_name", resourceGroup)
	d.SetId(*read.ID)

	return resourceArmAppServiceCustomHostnameCertificateBindingRead(d, meta)
}

func resourceArmAppServiceCustomHostnameCertificateBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
	}

	props := resp.HostNameBindingProperties
	if (props == nil) || (props.Thumbprint == nil) {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("ssl_state", props.SslState)
	d.Set("thumbprint", props.Thumbprint)
	d.Set("hostname", id.Name)
	d.Set("app_service_name", id.AppServiceName)
	d.Set("resource_group_name", id.ResourceGroup)

	return nil
}

func resourceArmAppServiceCustomHostnameCertificateBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
	}

	props := resp.HostNameBindingProperties
	if (props == nil) || (props.Thumbprint == nil) {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	locks.ByName(id.AppServiceName, appServiceCustomHostnameCertificateBindingResourceName)
	defer locks.UnlockByName(id.AppServiceName, appServiceCustomHostnameCertificateBindingResourceName)

	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q / Resource Group %q)", id.Name, id.AppServiceName, id.ResourceGroup)

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SslState:   web.SslStateDisabled,
			Thumbprint: nil,
		},
	}

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name, properties); err != nil {
		return fmt.Errorf("Error deleting Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
	}

	d.SetId("")
	return nil
}
