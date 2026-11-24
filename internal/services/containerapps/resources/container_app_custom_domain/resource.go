// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (a Resource) ResourceType() string {
	return "azurerm_container_app_custom_domain"
}

func (a Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ContainerAppCustomDomainId
}
