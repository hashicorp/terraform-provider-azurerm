// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSearchServiceCreate,
		Read:   resourceSearchServiceRead,
		Update: resourceSearchServiceUpdate,
		Delete: resourceSearchServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := services.ParseSearchServiceID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(validateSearchServiceSKUUpdate),
			pluginsdk.CustomizeDiffShim(validateSearchServiceApiAccessControlRbac),
		),

		Schema: searchServiceSchema(),
	}
}
