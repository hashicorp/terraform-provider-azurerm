// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maps_account

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMapsAccountCreate,
		Read:   resourceMapsAccountRead,
		Update: resourceMapsAccountUpdate,
		Delete: resourceMapsAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accounts.ParseAccountID(id)
			return err
		}),

		Schema: mapsAccountSchema(),
	}
}
