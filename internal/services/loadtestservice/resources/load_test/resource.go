// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package load_test

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return loadtests.ValidateLoadTestID
}

func (r Resource) ResourceType() string {
	return "azurerm_load_test"
}
