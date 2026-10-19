// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (Resource) Identity() resourceids.ResourceId {
	return &domains.DomainId{}
}

func (Resource) ResourceType() string {
	return "azurerm_email_communication_service_domain"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return domains.ValidateDomainID
}
