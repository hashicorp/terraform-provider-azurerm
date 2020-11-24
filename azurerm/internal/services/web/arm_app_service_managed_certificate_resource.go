package web

import (
	"fmt"
	"log"
	"time"

	azvalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func resourceArmAppServiceManagedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceManagedCertificateCreateUpdate,
		Read:   resourceArmAppServiceManagedCertificateRead,
		Update: resourceArmAppServiceManagedCertificateCreateUpdate,
		Delete: resourceArmAppServiceManagedCertificateDelete,
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
				ValidateFunc: validate.AppServiceName,
			},

			"canonical_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azvalidate.DomainName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"app_service_plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Managed Certificate creation.")

	name := d.Get("name").(string)
	appServicePlanID := d.Get("app_service_plan_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	id := parse.NewAppServiceManagedCertificateId(name, resourceGroup)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_managed_certificate", *existing.ID)
		}
	}

	certificate := web.Certificate{
		CertificateProperties: &web.CertificateProperties{
			CanonicalName: utils.String(d.Get("canonical_name").(string)),
			ServerFarmID:  utils.String(appServicePlanID),
			Password:      new(string),
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if resp, err := client.CreateOrUpdate(ctx, resourceGroup, name, certificate); err != nil {
		// API returns 202 where 200 is expected - https://github.com/Azure/azure-sdk-for-go/issues/13665
		if !utils.ResponseWasStatusCode(resp.Response, 202) {
			return fmt.Errorf("Error creating/updating App Service Managed Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	certificateWait := &resource.StateChangeConf{
		Pending:    []string{"NotFound", "Unknown"},
		Target:     []string{"Success"},
		MinTimeout: 10 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return "NotFound", "NotFound", err
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
		return fmt.Errorf("waiting for App Service Managed Certificate %q: %+v", id.Name, err)
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmAppServiceManagedCertificateRead(d, meta)
}

func resourceArmAppServiceManagedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceManagedCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Managed Certificate %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Managed Certificate %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.CertificateProperties; props != nil {
		d.Set("canonical_name", props.CanonicalName)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("subject_name", props.SubjectName)
		d.Set("host_names", props.HostNames)
		d.Set("issuer", props.Issuer)
		d.Set("issue_date", props.IssueDate.Format(time.RFC3339))
		d.Set("expiration_date", props.ExpirationDate.Format(time.RFC3339))
		d.Set("thumbprint", props.Thumbprint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceManagedCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceManagedCertificateID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Certificate %q (Resource Group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Certificate %q (Resource Group %q): %s)", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
