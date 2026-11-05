// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainAssociationCreate,
		Read:   resourceCdnFrontDoorCustomDomainAssociationRead,
		Update: resourceCdnFrontDoorCustomDomainAssociationUpdate,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.FrontDoorCustomDomainAssociationID(id)
			return err
		}, importCdnFrontDoorCustomDomainAssociation()),

		Schema: cdnFrontdoorCustomDomainAssociationSchema(),
	}
}
