// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint_application_security_group_association

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (p Resource) ResourceType() string {
	return "azurerm_private_endpoint_application_security_group_association"
}

func (p Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.PrivateEndpointApplicationSecurityGroupAssociationIDValidation
}
