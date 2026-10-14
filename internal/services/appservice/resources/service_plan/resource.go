// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_service_plan"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServicePlanID
}

func (r Resource) Identity() resourceids.ResourceId {
	return &commonids.AppServicePlanId{}
}
