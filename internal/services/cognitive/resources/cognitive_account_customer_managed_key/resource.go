// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_customer_managed_key

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceCognitiveAccountCustomerManagedKeyRead,
		Update: resourceCognitiveAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceCognitiveAccountCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cognitiveservicesaccounts.ParseAccountID(id)
			return err
		}),

		Schema: cognitiveAccountCustomerManagedKeySchema(),
	}
}
