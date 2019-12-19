package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAppServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppServiceCertificateRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

func dataSourceAppServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("App Service Certificate %q (Resource Group %q) does not exist!", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving App Service Certificate %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.CertificateProperties; props != nil {
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
