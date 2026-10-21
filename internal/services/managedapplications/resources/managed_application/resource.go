// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedApplicationCreate,
		Read:   resourceManagedApplicationRead,
		Update: resourceManagedApplicationUpdate,
		Delete: resourceManagedApplicationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applications.ParseApplicationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(resourceManagedApplicationCustomizeDiff),
		),

		Schema: resourceManagedApplicationSchema(),
	}
}
