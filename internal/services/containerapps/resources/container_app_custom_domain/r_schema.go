// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_custom_domain

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (a Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The fully qualified name of the Custom Domain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
			ValidateFunc: managedenvironments.ValidateCertificateID,
			RequiredWith: []string{"certificate_binding_type"},
		},

		"certificate_binding_type": {
			Description:  "The Certificate Binding type.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForBindingType(), false),
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

func (a Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			// Removing an explicit certificate (switching to an Azure-managed certificate) cannot be
			// expressed as an in-place update: the Update path only writes the certificate when one is
			// configured, so clearing it would leave the previous certificate bound at Azure and
			// produce an inconsistent result after apply. Force a new resource for that transition.
			if metadata.ResourceDiff.HasChange("container_app_environment_certificate_id") {
				oldCertID, newCertID := metadata.ResourceDiff.GetChange("container_app_environment_certificate_id")
				if oldCertID.(string) != "" && newCertID.(string) == "" {
					return metadata.ResourceDiff.ForceNew("container_app_environment_certificate_id")
				}
			}

			return nil
		},
	}
}
