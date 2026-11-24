// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_trust

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const DomainServiceResourceName = "azurerm_active_directory_domain_service"

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_active_directory_domain_service_trust"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DomainServiceTrustID
}
