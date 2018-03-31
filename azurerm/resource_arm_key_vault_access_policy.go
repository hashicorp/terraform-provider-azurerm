package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func resourceArmKeyVaultAccessPolicy() *schema.Resource {
	// The Azure SDK doesn't provide direct access to KV access policies,
	// so changing access policies requires a full Key Vault resource
	// be provided. This is done here by reusing the resourceArmKeyVault
	// schema with some tweaks so that unneed properties need not be
	// specified in the template.

	kvapSchema := resourceArmKeyVault()
	// Reuse
	kvapSchema.Create = resourceArmKeyVaultAccessPolicyCreateOrUpdate
	kvapSchema.Read = resourceArmKeyVaultRead // Redundant, but included for clarity
	kvapSchema.Update = resourceArmKeyVaultAccessPolicyCreateOrUpdate
	kvapSchema.Delete = resourceArmKeyVaultAccessPolicyDelete

	kvapSchema.Schema["sku"].Required = false
	kvapSchema.Schema["sku"].Computed = true
	kvapSchema.Schema["enabled_for_deployment"].Computed = true
	kvapSchema.Schema["enabled_for_disk_encryption"].Computed = true
	kvapSchema.Schema["enabled_for_template_deployment"].Computed = true

	return kvapSchema
}

func resourceArmKeyVaultAccessPolicyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM KeyVault creation.")

	newAccessPolicies := expandKeyVaultAccessPolicies(d)

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	// Reuse all existing Key Vault properties
	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Cannot find Key Vault '%s' in resource group '%s'", name, resGroup)
	}

	d.SetId(*read.ID)
	resourceArmKeyVaultRead(d, meta)

	location := d.Get("location").(string)
	tenantUUID := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	enabledForDeployment := d.Get("enabled_for_deployment").(bool)
	enabledForDiskEncryption := d.Get("enabled_for_disk_encryption").(bool)
	enabledForTemplateDeployment := d.Get("enabled_for_template_deployment").(bool)
	tags := d.Get("tags").(map[string]interface{})

	parameters := keyvault.VaultCreateOrUpdateParameters{
		Location: &location,
		Properties: &keyvault.VaultProperties{
			TenantID:                     &tenantUUID,
			Sku:                          expandKeyVaultSku(d),
			AccessPolicies:               newAccessPolicies,
			EnabledForDeployment:         &enabledForDeployment,
			EnabledForDiskEncryption:     &enabledForDiskEncryption,
			EnabledForTemplateDeployment: &enabledForTemplateDeployment,
		},
		Tags: expandTags(tags),
	}

	_, err = client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	return resourceArmKeyVaultRead(d, meta)
}

func resourceArmKeyVaultAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	d.Set("access_policy", nil)

	return resourceArmKeyVaultAccessPolicyCreateOrUpdate(d, meta)
}
