// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_automanage_configuration_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct {
	ArcMachineId    string `tfschema:"arc_machine_id"`
	ConfigurationId string `tfschema:"configuration_id"`
}

func (v Resource) ResourceType() string {
	return "azurerm_arc_machine_automanage_configuration_assignment"
}

func (v Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationprofilehcrpassignments.ValidateProviders2ConfigurationProfileAssignmentID
}

var _ sdk.Resource = &Resource{}
