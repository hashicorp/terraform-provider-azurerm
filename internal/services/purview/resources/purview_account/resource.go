// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-12-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePurviewAccountCreate,
		Read:   resourcePurviewAccountRead,
		Update: resourcePurviewAccountUpdate,
		Delete: resourcePurviewAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := account.ParseAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: purviewAccountSchema(),
	}
}
