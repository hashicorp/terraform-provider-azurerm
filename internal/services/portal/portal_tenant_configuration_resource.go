package portal

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePortalTenantConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePortalTenantConfigurationCreateUpdate,
		Read:   resourcePortalTenantConfigurationRead,
		Update: resourcePortalTenantConfigurationCreateUpdate,
		Delete: resourcePortalTenantConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := tenantconfiguration.ParseConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"private_markdown_storage_enforced": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourcePortalTenantConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := tenantconfiguration.NewConfigurationID("default")

	if d.IsNewResource() {
		existing, err := client.TenantConfigurationsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_portal_tenant_configuration", id.ID())
		}
	}

	parameters := tenantconfiguration.Configuration{
		Properties: &tenantconfiguration.ConfigurationProperties{
			EnforcePrivateMarkdownStorage: utils.Bool(d.Get("private_markdown_storage_enforced").(bool)),
		},
	}

	if _, err := client.TenantConfigurationsCreate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePortalTenantConfigurationRead(d, meta)
}

func resourcePortalTenantConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tenantconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TenantConfigurationsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("private_markdown_storage_enforced", props.EnforcePrivateMarkdownStorage)
		}
	}

	return nil
}

func resourcePortalTenantConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.TenantConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tenantconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.TenantConfigurationsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
