package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultManagedHardwareSecurityModule() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceKeyVaultManagedHardwareSecurityModuleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ManagedHardwareSecurityModuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_object_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"purge_protection_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"soft_delete_retention_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"hsm_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultManagedHardwareSecurityModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewManagedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku_name", skuName)

	if props := resp.Properties; props != nil {
		tenantId := ""
		if tid := props.TenantID; tid != nil {
			tenantId = tid.String()
		}
		d.Set("tenant_id", tenantId)
		d.Set("admin_object_ids", utils.FlattenStringSlice(props.InitialAdminObjectIds))
		d.Set("hsm_uri", props.HsmURI)
		d.Set("purge_protection_enabled", props.EnablePurgeProtection)
		d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
