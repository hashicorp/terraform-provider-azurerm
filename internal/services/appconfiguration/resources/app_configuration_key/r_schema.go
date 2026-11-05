// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (k Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// User-specified segments are lowercased in the API response
			// tracked in https://github.com/Azure/azure-rest-api-specs/issues/24337
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     configurationstores.ValidateConfigurationStoreID,
		},
		"key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C We set some values in this field depending on the `type` so this needs to remain Computed
			Computed: true,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			// NOTE: O+C The value of this is updated anytime the resource changes
			Optional: true,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"value": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"vault_key_reference",
			},
			// if `type` is set to `vault`, then `value` will be set by `vault_key_reference`
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return d.Get("type").(string) == KeyTypeVault && d.Get("vault_key_reference").(string) != "" && old != "" && new == ""
			},
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      KeyTypeKV,
			ValidateFunc: validation.StringInSlice([]string{KeyTypeVault, KeyTypeKV}, false),
		},
		"vault_key_reference": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			ConflictsWith: []string{
				"value",
			},
		},
		"tags": commonschema.Tags(),
	}
}

func (k Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			keyType := rd.Get("type").(string)

			if keyType == KeyTypeVault {
				contentType := rd.Get("content_type").(string)
				if rd.HasChange("content_type") && contentType != VaultKeyContentType {
					return fmt.Errorf("key type %q cannot have content type other than %q (found %q)", KeyTypeVault, VaultKeyContentType, contentType)
				}

				if rd.HasChange("value") && rd.Get("value").(string) != "" {
					return fmt.Errorf("'value' should only be set when key type is set to %q", KeyTypeKV)
				}
			}

			if keyType == KeyTypeKV && rd.Get("vault_key_reference").(string) != "" {
				return fmt.Errorf("'vault_key_reference' should only be set when key type is set to %q", KeyTypeVault)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
