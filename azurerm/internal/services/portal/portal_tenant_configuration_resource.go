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

func resourceTenantConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTenantConfigurationCreateUpdate,
		Read:   resourceTenantConfigurationRead,
		Update: resourceTenantConfigurationCreateUpdate,
		Delete: resourceTenantConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TenantConfigurationID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"enforce_private_markdown_storage": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceTenantConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTenantConfigurationID("default")

	if d.IsNewResource() {
		existing, err := client.Get(ctx)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_tenant_configuration", id.ID())
		}
	}

	tenantConfiguration := portal.Configuration{
		ConfigurationProperties: &portal.ConfigurationProperties{
			EnforcePrivateMarkdownStorage: utils.Bool(d.Get("enforce_private_markdown_storage").(bool)),
		},
	}

	if _, err := client.Create(ctx, tenantConfiguration); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceTenantConfigurationRead(d, meta)
}

func resourceTenantConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantConfigurationID(d.Id())
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
		d.Set("enforce_private_markdown_storage", props.EnforcePrivateMarkdownStorage)
	}

	return nil
}

func resourceTenantConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TenantConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
