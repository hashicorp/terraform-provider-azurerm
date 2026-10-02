// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMarketplaceAgreementCreate,
		Read:   resourceMarketplaceAgreementRead,
		Delete: resourceMarketplaceAgreementDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := agreements.ParsePlanID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: marketplaceAgreementSchema(),
	}
}
