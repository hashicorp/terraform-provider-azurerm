package recoveryservices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceBackupPolicyVm() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBackupPolicyVmRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceBackupPolicyVmRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	log.Printf("[DEBUG] Reading Recovery Service  Policy %q (resource group %q)", name, resourceGroup)

	protectionPolicy, err := client.Get(ctx, vaultName, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(protectionPolicy.Response) {
			return fmt.Errorf("Error: Backup Policy %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Backup Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	id := strings.Replace(*protectionPolicy.ID, "Subscriptions", "subscriptions", 1)
	d.SetId(id)

	return tags.FlattenAndSet(d, protectionPolicy.Tags)
}
