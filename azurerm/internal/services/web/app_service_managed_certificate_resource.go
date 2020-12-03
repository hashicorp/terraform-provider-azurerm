package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceManagedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceManagedCertificateCreate,
		Read:   resourceArmAppServiceManagedCertificateRead,
		Update: resourceArmAppServiceManagedCertificateUpdate,
		Delete: resourceArmAppServiceManagedCertificateDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagedCertificateID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"custom_hostname_binding_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceCustomHostnameBindingID,
			},

			"ssl_state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.SslStateIPBasedEnabled),
					string(web.SslStateSniEnabled),
				}, false),
			},

			"canonical_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"host_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"issue_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// TODO Remove in 3.0
			"tags": {
				Type:       schema.TypeMap,
				Optional:   true,
				Deprecated: "Tags are not stored by the service and will be ignored, this property will be removed in a future version of the provider",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmAppServiceManagedCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	appServiceClient := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	customHostnameBindingId, err := parse.HostnameBindingID(d.Get("custom_hostname_binding_id").(string))
	if err != nil {
		return err
	}

	binding, err := appServiceClient.GetHostNameBinding(ctx, customHostnameBindingId.ResourceGroup, customHostnameBindingId.SiteName, customHostnameBindingId.Name)
	if err != nil {
		return fmt.Errorf("failed retrieving Hostname Binding to update for Managed Certificate")
	}

	appService, err := appServiceClient.Get(ctx, customHostnameBindingId.ResourceGroup, customHostnameBindingId.SiteName)
	if err != nil {
		return fmt.Errorf("could not retrieve App Service Custom Hostname details for %q", customHostnameBindingId.Name)
	}

	name := customHostnameBindingId.Name
	appServicePlanID := ""
	if appService.SiteProperties == nil || appService.SiteProperties.ServerFarmID == nil {
		return fmt.Errorf("could not get App Service Plan ID for Custom Hostname Binding %q (resource group %q)", customHostnameBindingId.Name, customHostnameBindingId.ResourceGroup)
	}
	appServicePlanID = *appService.SiteProperties.ServerFarmID

	appServiceLocation := ""
	if appService.Location != nil {
		appServiceLocation = location.Normalize(*appService.Location)
	}

	id := parse.NewManagedCertificateID(subscriptionId, customHostnameBindingId.ResourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.CertificateName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing App Service Certificate %q (Resource Group %q): %s", id.CertificateName, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_managed_certificate", *existing.ID)
		}
	}

	certificate := web.Certificate{
		CertificateProperties: &web.CertificateProperties{
			CanonicalName: utils.String(customHostnameBindingId.Name),
			ServerFarmID:  utils.String(appServicePlanID),
			Password:      new(string),
		},
		Location: utils.String(appServiceLocation),
	}

	resp, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.CertificateName, certificate)
	if err != nil {
		// API returns 202 where 200 is expected - https://github.com/Azure/azure-sdk-for-go/issues/13665
		if !utils.ResponseWasStatusCode(resp.Response, 202) {
			return fmt.Errorf("Error creating/updating App Service Managed Certificate %q (Resource Group %q): %s", id.CertificateName, id.ResourceGroup, err)
		}
	}

	certificateWait := &resource.StateChangeConf{
		Pending:    []string{"NotFound", "Unknown"},
		Target:     []string{"Success"},
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id.ResourceGroup, id.CertificateName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return "NotFound", "NotFound", nil
				}
				return "Unknown", "Unknown", err
			}
			if utils.ResponseWasStatusCode(resp.Response, 200) {
				return "Success", "Success", nil
			}
			return "Unknown", "Unknown", err
		},
	}

	if !d.IsNewResource() {
		certificateWait.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := certificateWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for App Service Managed Certificate %q: %+v", id.CertificateName, err)
	}

	d.SetId(id.ID(""))

	// Get the cert again, create doesn't return the Thumbprint
	resp, err = client.Get(ctx, id.ResourceGroup, id.CertificateName)
	if err != nil {
		return fmt.Errorf("could not read Managed Certificate %q (resource group %q) after Creation: %+v", id.CertificateName, id.ResourceGroup, err)
	}

	// Update the binding with the new Cert
	sslState := d.Get("ssl_state").(string)
	if resp.Thumbprint != nil {
		if binding.HostNameBindingProperties != nil {
			binding.HostNameBindingProperties.SslState = web.SslState(sslState)
			binding.HostNameBindingProperties.Thumbprint = resp.Thumbprint
		} else {
			return fmt.Errorf("failed to read Custom Hostname Binding properties for %q (resource group %q)", customHostnameBindingId.Name, customHostnameBindingId.ResourceGroup)
		}
	} else {
		return fmt.Errorf("could not read Thumbprint for Managed Certificate %q (resource group %q) to apply to Custom Hostname Binsing", id.CertificateName, id.ResourceGroup)
	}

	_, err = appServiceClient.CreateOrUpdateHostNameBinding(ctx, customHostnameBindingId.ResourceGroup, customHostnameBindingId.SiteName, customHostnameBindingId.Name, binding)
	if err != nil {
		return fmt.Errorf("failed to update Hostname Binding for %q (resource group %q): %+v", customHostnameBindingId.Name, customHostnameBindingId.ResourceGroup, err)
	}

	return resourceArmAppServiceManagedCertificateRead(d, meta)
}

func resourceArmAppServiceManagedCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	appServiceClient := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("ssl_state") {
		sslState := web.SslState(d.Get("ssl_state").(string))
		hostnameBindingId, err := parse.HostnameBindingID(d.Get("custom_hostname_binding_id").(string))
		if err != nil {
			return err
		}
		// Get the current binding
		binding, err := appServiceClient.GetHostNameBinding(ctx, hostnameBindingId.ResourceGroup, hostnameBindingId.SiteName, hostnameBindingId.Name)
		if err != nil {
			return fmt.Errorf("could not retrieve Hostname Binding %q (resource group %q) for update: %+v", hostnameBindingId.Name, hostnameBindingId.ResourceGroup, err)
		}
		if binding.HostNameBindingProperties != nil {
			binding.HostNameBindingProperties.SslState = sslState
		}

		// Remove binding
		if _, err := appServiceClient.DeleteHostNameBinding(ctx, hostnameBindingId.ResourceGroup, hostnameBindingId.SiteName, hostnameBindingId.Name); err != nil {
			return fmt.Errorf("could not remove Hostname Binding from %q (resource group %q) to change `ssl_state`: %+v", hostnameBindingId.Name, hostnameBindingId.ResourceGroup, err)
		}

		// re-add binding with new SSL State
		if _, err = appServiceClient.CreateOrUpdateHostNameBinding(ctx, hostnameBindingId.ResourceGroup, hostnameBindingId.SiteName, hostnameBindingId.Name, binding); err != nil {
			return fmt.Errorf("failed to update Hostname Binding %q (resource group %q) for new SSL State: %+v", hostnameBindingId.ResourceGroup, hostnameBindingId.Name, err)
		}
	}

	return resourceArmAppServiceManagedCertificateRead(d, meta)
}

func resourceArmAppServiceManagedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	appServicesClient := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Managed Certificate %q (Resource Group %q) was not found - removing from state", id.CertificateName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Managed Certificate %q (Resource Group %q): %+v", id.CertificateName, id.ResourceGroup, err)
	}

	if props := resp.CertificateProperties; props != nil {
		d.Set("canonical_name", props.CanonicalName)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("subject_name", props.SubjectName)
		d.Set("host_names", props.HostNames)
		d.Set("issuer", props.Issuer)
		d.Set("issue_date", props.IssueDate.Format(time.RFC3339))
		expirationDate := ""
		if props.ExpirationDate != nil {
			expirationDate = props.ExpirationDate.Format(time.RFC3339)
		}
		d.Set("expiration_date", expirationDate)
		d.Set("thumbprint", props.Thumbprint)
	}

	// @jackofallops - Best effort here, won't set a value on import as the setting is actually on the binding itself and the certificate ID doesn't provide enough information to get it
	sslState := ""
	if hostnameBindingIdRaw, ok := d.GetOk("custom_hostname_binding_id"); ok {
		hostnameBindingId, err := parse.HostnameBindingID(hostnameBindingIdRaw.(string))
		if err != nil {
			return fmt.Errorf("could not parse ID for Hostname Binding: %+v", err)
		}
		binding, err := appServicesClient.GetHostNameBinding(ctx, hostnameBindingId.ResourceGroup, hostnameBindingId.SiteName, hostnameBindingId.Name)
		if err != nil {
			return fmt.Errorf("could not get details of Hostname Binding: %+v", err)
		}
		if props := binding.HostNameBindingProperties; props != nil {
			sslState = string(props.SslState)
		}
	}
	d.Set("ssl_state", sslState)

	return nil
}

func resourceArmAppServiceManagedCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	appServiceClient := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedCertificateID(d.Id())
	if err != nil {
		return err
	}

	// We need to detach the cert before we can delete it, since this is a best effort since we can't guarantee resource data here
	hostnameBindingRaw, ok := d.GetOk("custom_hostname_binding_id")
	if !ok {
		return fmt.Errorf("could not remove certificate from Hostname Binding, missing `custom_hostname_binding_id`")
	}

	hostnameBinding, err := parse.HostnameBindingID(hostnameBindingRaw.(string))
	if err != nil {
		return err
	}

	binding, err := appServiceClient.GetHostNameBinding(ctx, hostnameBinding.ResourceGroup, hostnameBinding.SiteName, hostnameBinding.Name)
	if err != nil {
		return err
	}

	if binding.HostNameBindingProperties != nil {
		binding.HostNameBindingProperties.SslState = web.SslStateDisabled
		binding.HostNameBindingProperties.Thumbprint = nil
	}
	_, err = appServiceClient.CreateOrUpdateHostNameBinding(ctx, hostnameBinding.ResourceGroup, hostnameBinding.SiteName, hostnameBinding.Name, binding)
	if err != nil {
		return fmt.Errorf("could not remove Managed Certificate %q from Hostname Binding %q (resource group %q): %+v", id.CertificateName, hostnameBinding.Name, hostnameBinding.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Deleting App Service Certificate %q (Resource Group %q)", id.CertificateName, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.CertificateName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Certificate %q (Resource Group %q): %s)", id.CertificateName, id.ResourceGroup, err)
		}
	}

	return nil
}
