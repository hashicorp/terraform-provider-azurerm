// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dynamic_scope

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (Resource) ResourceType() string {
	return "azurerm_maintenance_assignment_dynamic_scope"
}

func (Resource) IDValidationFunc() func(any, string) ([]string, []error) {
	return configurationassignments.ValidateConfigurationAssignmentID
}
