// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_automanage_configuration_assignment

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (v Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"arc_machine_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationprofileassignments.ValidateMachineID,
		},
		"configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
		},
	}
}

func (v Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
