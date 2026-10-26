// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_registration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_eventgrid_partner_registration"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return partnerregistrations.ValidatePartnerRegistrationID
}
