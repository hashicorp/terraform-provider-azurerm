// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConfidentialLedgerCreate,
		Read:   resourceConfidentialLedgerRead,
		Update: resourceConfidentialLedgerUpdate,
		Delete: resourceConfidentialLedgerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := confidentialledger.ParseLedgerID(id)
			return err
		}),

		Schema: confidentialLedgerSchema(),
	}
}
