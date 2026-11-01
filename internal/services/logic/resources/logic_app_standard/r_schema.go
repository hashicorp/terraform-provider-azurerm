// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogicAppStandardName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"app_service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"use_extension_bundle": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"bundle_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "[1.*, 2.0.0)",
		},

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeRequired,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"scm_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"site_config": helpers.SchemaLogicAppStandardSiteConfig(),

		"connection_string": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForConnectionStringType(), false),
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			RequiredWith: []string{
				"storage_account_access_key",
			},
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			RequiredWith: []string{
				"storage_account_name",
			},
			ConflictsWith: []string{
				"storage_key_vault_secret_id",
			},
		},

		"storage_key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		// Once this property is set, it can not be removed while identity is UserAssigned.
		// tracked on https://github.com/Azure/azure-rest-api-specs/issues/37553
		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  helpers.PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice([]string{
				helpers.PublicNetworkAccessEnabled,
				helpers.PublicNetworkAccessDisabled,
			}, false),
		},

		"storage_account_share_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "~4",
		},

		"vnet_content_share_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"possible_outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"site_credential": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// If `key_vault_reference_identity_id` is not set in config, check if ForceNew is needed
			if metadata.ResourceDiff.GetRawConfig().GetAttr("key_vault_reference_identity_id").IsNull() {
				// If it was previously set and identity is purely UserAssigned, ForceNew is required
				// because the API does not support unsetting this property
				old, _ := metadata.ResourceDiff.GetChange("key_vault_reference_identity_id")
				if old.(string) != "" {
					identityType := metadata.ResourceDiff.Get("identity.0.type").(string)
					if strings.EqualFold(identityType, "UserAssigned") {
						if err := metadata.ResourceDiff.ForceNew("key_vault_reference_identity_id"); err != nil {
							return err
						}
					}
				}
				return nil
			}

			keyVaultReferenceIdentityId := metadata.ResourceDiff.Get("key_vault_reference_identity_id").(string)
			if keyVaultReferenceIdentityId == "" {
				return nil
			}

			identityIdsRaw := metadata.ResourceDiff.Get("identity.0.identity_ids").(*pluginsdk.Set)
			identityIds := identityIdsRaw.List()

			for _, id := range identityIds {
				if strings.EqualFold(id.(string), keyVaultReferenceIdentityId) {
					return nil
				}
			}

			return fmt.Errorf("`key_vault_reference_identity_id` must be an identity assigned to this resource in the `identity` block, got `%s`", keyVaultReferenceIdentityId)
		},
	}
}
