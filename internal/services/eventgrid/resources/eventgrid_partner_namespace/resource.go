// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return partnernamespaces.ValidatePartnerNamespaceID
}
