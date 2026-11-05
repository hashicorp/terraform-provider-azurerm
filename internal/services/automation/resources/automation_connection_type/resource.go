// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = (*Resource)(nil)

func (m Resource) ResourceType() string {
	return "azurerm_automation_connection_type"
}

func (m Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectiontype.ValidateConnectionTypeID
}
