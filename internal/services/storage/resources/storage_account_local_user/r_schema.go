// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_local_user

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalUserName,
		},
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},
		"ssh_key_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			AtLeastOneOf: []string{"ssh_key_enabled", "ssh_password_enabled"},
		},
		"ssh_password_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			AtLeastOneOf: []string{"ssh_key_enabled", "ssh_password_enabled"},
		},
		"home_directory": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"ssh_authorized_key": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			RequiredWith: []string{"ssh_key_enabled"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     computevalidate.SSHKey,
						DiffSuppressFunc: suppress.SSHKey,
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"permission_scope": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"permissions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"read": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"write": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"delete": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"list": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"create": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
					"service": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(
							// Replace the string literal with enum once https://github.com/Azure/azure-rest-api-specs/pull/21845 is merged
							[]string{"blob", "file"},
							false,
						),
					},
					"resource_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sid": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
		"password": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if diff.HasChange("ssh_password_enabled") {
				if err := diff.SetNewComputed("password"); err != nil {
					return err
				}
			}
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}
