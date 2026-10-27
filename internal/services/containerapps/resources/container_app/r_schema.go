// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerAppName,
			Description:  "The name for this Container App.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to host this Container App.",
		},

		"template": helpers.ContainerTemplateSchema(),

		"revision_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForActiveRevisionsMode(), false),
		},

		"ingress": helpers.ContainerAppIngressSchema(),

		"registry": helpers.ContainerAppRegistrySchema(),

		"secret": helpers.SecretsSchema(),

		"dapr": helpers.ContainerDaprSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"workload_profile_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"max_inactive_revisions": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"latest_revision_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the latest Container Revision.",
		},

		"latest_revision_fqdn": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The FQDN of the Latest Revision of the Container App.",
		},

		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The ID of the Custom Domain Verification for this Container App.",
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}
			var app ContainerAppModel
			if err := metadata.DecodeDiff(&app); err != nil {
				return err
			}
			// Ingress traffic weight validations
			if len(app.Ingress) != 0 {
				ingress := app.Ingress[0]

				for i, tw := range ingress.TrafficWeights {
					if !tw.LatestRevision && tw.RevisionSuffix == "" {
						return fmt.Errorf("`either ingress.0.traffic_weight.%[1]d.revision_suffix` or `ingress.0.traffic_weight.%[1]d.latest_revision` should be specified", i)
					}
				}
			}

			for _, s := range app.Secrets {
				if s.KeyVaultSecretId != "" && s.Identity == "" {
					return fmt.Errorf("secret %s must supply identity for key vault secret id", s.Name)
				}
				if s.KeyVaultSecretId == "" && s.Identity != "" {
					return fmt.Errorf("secret %s must supply key vault secret id when specifying identity", s.Name)
				}
			}
			return nil
		},
	}
}
