// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_network_mapping

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (s Resource) ResourceType() string {
	return "azurerm_site_recovery_hyperv_network_mapping"
}

func (s Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationnetworkmappings.ValidateReplicationNetworkMappingID
}
