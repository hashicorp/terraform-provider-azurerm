// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_api_management_workspace_certificate"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificate.ValidateWorkspaceCertificateID
}
