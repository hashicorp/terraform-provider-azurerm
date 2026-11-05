// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk_sas_token

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedDiskSasTokenCreate,
		Read:   resourceManagedDiskSasTokenRead,
		Delete: resourceManagedDiskSasTokenDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseManagedDiskID(id)
			return err
		}),

		Schema: managedDiskSasTokenSchema(),
	}
}
