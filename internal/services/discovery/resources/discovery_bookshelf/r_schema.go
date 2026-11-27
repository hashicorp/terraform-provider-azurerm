// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`),
				"`name` must be between 3 and 24 characters long and may only contain letters, numbers and hyphens",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"customer_managed_key": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"log_analytics_cluster_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"log_analytics_cluster_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateClusterID,
		},

		"private_endpoint_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"search_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"workload_identity_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_on_behalf_of_broker_resource_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"managed_resource_group_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			if err := validateDiscoveryBookshelfCustomerManagedKeyIdentity(diff); err != nil {
				return err
			}

			if !diff.HasChange("customer_managed_key") {
				return nil
			}

			oldCMK, newCMK := diff.GetChange("customer_managed_key")
			if len(oldCMK.([]any)) != len(newCMK.([]any)) {
				return diff.ForceNew("customer_managed_key")
			}

			if !diff.NewValueKnown("customer_managed_key.0.key_vault_key_id") {
				return nil
			}

			oldKeyId, newKeyId := diff.GetChange("customer_managed_key.0.key_vault_key_id")
			oldKey, err := keyvault.ParseNestedItemID(oldKeyId.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}
			newKey, err := keyvault.ParseNestedItemID(newKeyId.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			if !strings.EqualFold(oldKey.KeyVaultBaseURL, newKey.KeyVaultBaseURL) {
				return diff.ForceNew("customer_managed_key.0.key_vault_key_id")
			}

			return nil
		},
	}
}

func validateDiscoveryBookshelfCustomerManagedKeyIdentity(diff *pluginsdk.ResourceDiff) error {
	identityId, ok := diff.GetOk("customer_managed_key.0.user_assigned_identity_id")
	if !ok || !diff.NewValueKnown("customer_managed_key.0.user_assigned_identity_id") {
		return nil
	}

	rawConfig := diff.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.GetAttr("workload_identity_ids").IsWhollyKnown() {
		return nil
	}

	for _, v := range diff.Get("workload_identity_ids").(*pluginsdk.Set).List() {
		if strings.EqualFold(v.(string), identityId.(string)) {
			return nil
		}
	}

	return fmt.Errorf("`customer_managed_key.0.user_assigned_identity_id` must be one of the identities specified in `workload_identity_ids`")
}
