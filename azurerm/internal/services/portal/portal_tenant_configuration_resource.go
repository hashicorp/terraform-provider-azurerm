package portal

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePortalTenantConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortalTenantConfigurationCreateUpdate,
		Read:   resourcePortalTenantConfigurationRead,
		Update: resourcePortalTenantConfigurationCreateUpdate,
		Delete: resourcePortalTenantConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PortalTenantConfigurationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"private_markdown_storage_enforced": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourcePortalTenantConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPortalTenantConfigurationID("default")

	if d.IsNewResource() {
		existing, err := client.Get(ctx)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_portal_tenant_configuration", id.ID())
		}
	}

	parameters := portal.Configuration{
		ConfigurationProperties: &portal.ConfigurationProperties{
			EnforcePrivateMarkdownStorage: utils.Bool(d.Get("private_markdown_storage_enforced").(bool)),
		},
	}

	if _, err := client.Create(ctx, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePortalTenantConfigurationRead(d, meta)
}

func resourcePortalTenantConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PortalTenantConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if props := resp.ConfigurationProperties; props != nil {
		d.Set("private_markdown_storage_enforced", props.EnforcePrivateMarkdownStorage)
	}

	return nil
}

func resourcePortalTenantConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PortalTenantConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
