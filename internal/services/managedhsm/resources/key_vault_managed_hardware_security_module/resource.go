// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmKeyVaultManagedHardwareSecurityModuleCreate,
		Read:   resourceArmKeyVaultManagedHardwareSecurityModuleRead,
		Delete: resourceArmKeyVaultManagedHardwareSecurityModuleDelete,
		Update: resourceArmKeyVaultManagedHardwareSecurityModuleUpdate,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedhsms.ParseManagedHSMID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(keyVaultHSMCustomizeDiff),

		Schema: keyVaultManagedHardwareSecurityModuleSchema(),
	}
}
