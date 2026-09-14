// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package suppression

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_advisor_suppression"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return suppressions.ValidateScopedSuppressionID
}
