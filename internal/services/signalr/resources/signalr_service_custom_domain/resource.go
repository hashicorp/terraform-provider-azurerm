// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_domain

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name signalr_service_custom_domain -service-package-name signalr -properties "name" -compare-values "subscription_id:signalr_service_id,resource_group_name:signalr_service_id,signalr_name:signalr_service_id" -test-env-vars "ARM_TEST_DNS_ZONE,ARM_TEST_DATA_RESOURCE_GROUP"

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &signalr.CustomDomainId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_signalr_service_custom_domain"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return signalr.ValidateCustomDomainID
}
