package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceManagedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceManagedCertificateCreateUpdate,
		Read:   resourceArmAppServiceManagedCertificateRead,
		Update: resourceArmAppServiceManagedCertificateCreateUpdate,
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
				ValidateFunc: validate.CustomHostnameBindingID,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAppServiceManagedCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	appServiceClient := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	customHostnameBindingId, err := parse.AppServiceCustomHostnameBindingID(d.Get("custom_hostname_binding_id").(string))
	if err != nil {
		return err
	}

	appService, err := appServiceClient.Get(ctx, customHostnameBindingId.ResourceGroup, customHostnameBindingId.AppServiceName)
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

	t := d.Get("tags").(map[string]interface{})

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
		Tags:     tags.Expand(t),
	}

	if resp, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.CertificateName, certificate); err != nil {
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

	return resourceArmAppServiceManagedCertificateRead(d, meta)
}

func resourceArmAppServiceManagedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceManagedCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedCertificateID(d.Id())
	if err != nil {
		return err
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
