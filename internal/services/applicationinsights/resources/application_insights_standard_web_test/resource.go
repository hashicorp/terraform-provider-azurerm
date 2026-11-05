// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_standard_web_test

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
)

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_application_insights_standard_web_test"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webtestsapis.ValidateWebTestID
}
