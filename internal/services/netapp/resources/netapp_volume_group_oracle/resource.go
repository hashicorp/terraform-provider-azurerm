// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_group_oracle

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_netapp_volume_group_oracle"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}
