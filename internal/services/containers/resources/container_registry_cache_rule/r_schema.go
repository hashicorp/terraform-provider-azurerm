// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_cache_rule

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "The name of the Container Registry Cache Rule. Only Alphanumeric characters allowed.",
			ValidateFunc: validate.ContainerRegistryCacheRuleName,
		},
		"container_registry_id": {
			Description:  "The ID of the Container Registry where the Cache Rule should apply.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},
		"credential_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Description:  "The ARM resource ID of the Credential Store which is associated with the Cache Rule.",
			ValidateFunc: credentialsets.ValidateCredentialSetID,
		},
		"source_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the source repository path.",
		},

		"target_repo": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the new repository path to store artifacts.",
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("credential_set_id"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("credential_set_id"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
