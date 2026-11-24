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
			Description:  "The fully qualified name of the Custom Domain.",
		},

		"container_app_id": {
			Description:  "The ID of the Container App to which this Custom Domain should be bound.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerapps.ValidateContainerAppID,
		},

		"container_app_environment_certificate_id": {
			Description:  "The ID of the Container App Environment Certificate to use.",
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
			Description:  "The Certificate Binding type.",
		},
	}
}

func (a Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_managed_certificate_id": {
			Description: "The ID of the Container App Environment Managed Certificate to use.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
