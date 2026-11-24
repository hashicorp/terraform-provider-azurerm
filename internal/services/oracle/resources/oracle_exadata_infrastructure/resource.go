// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_oracle_exadata_infrastructure"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID
}
