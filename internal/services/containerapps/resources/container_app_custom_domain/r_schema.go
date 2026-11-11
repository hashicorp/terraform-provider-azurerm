// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (a Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The hostname of the Certificate. Must be the CN or a named SAN in the certificate.",
		},

		"container_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerapps.ValidateContainerAppID,
		},

		"container_app_environment_certificate_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"certificate_binding_type"},
			ValidateFunc: managedenvironments.ValidateCertificateID,
		},

		"certificate_binding_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForBindingType(), false),
			Description:  "The Binding type. Possible values include `Disabled` and `SniEnabled`.",
		},
	}
}

func (a Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_managed_certificate_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
