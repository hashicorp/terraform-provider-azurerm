// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_automanage_configuration_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct {
	VirtualMachineId string `tfschema:"virtual_machine_id"`
	ConfigurationId  string `tfschema:"configuration_id"`
}

func (v Resource) ResourceType() string {
	return "azurerm_virtual_machine_automanage_configuration_assignment"
}

func (v Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationprofileassignments.ValidateVirtualMachineProviders2ConfigurationProfileAssignmentID
}

var _ sdk.Resource = &Resource{}
