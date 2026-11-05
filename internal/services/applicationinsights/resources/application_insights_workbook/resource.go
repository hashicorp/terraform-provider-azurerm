// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_application_insights_workbook"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workbooksapis.ValidateWorkbookID
}
