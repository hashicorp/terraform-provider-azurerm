// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_resource_guard

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionResourceGuardCreateUpdate,
		Read:   resourceDataProtectionResourceGuardRead,
		Update: resourceDataProtectionResourceGuardCreateUpdate,
		Delete: resourceDataProtectionResourceGuardDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&resourceguardresources.ResourceGuardId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&resourceguardresources.ResourceGuardId{}),
		},

		Schema: dataProtectionResourceGuardSchema(),
	}
}
