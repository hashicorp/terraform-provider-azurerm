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
			Description:  "The name of the Logic App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogicAppStandardName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"app_service_plan_id": {
			Description:  "The ID of the App Service Plan within which to create this Logic App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"app_settings": {
			Description: "A map of key-value pairs for [App Settings](https://learn.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"use_extension_bundle": {
			Description: "Whether the Logic App should use the bundled extension package. If true, then application settings for `AzureFunctionsJobHost__extensionBundle__id` and `AzureFunctionsJobHost__extensionBundle__version` will be created.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"bundle_version": {
			Description: "If `use_extension_bundle` is set to `true` this controls the allowed range for bundle versions.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "[1.*, 2.0.0)",
		},

		"client_affinity_enabled": {
			Description: "Whether the Logic App should send session affinity cookies, which route client requests in the same session to the same instance.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"client_certificate_mode": {
			Description:  "The mode of the Logic App's client certificates requirement for incoming requests.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeRequired,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"enabled": {
			Description: "Whether the Logic App is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Description: "Whether the FTP basic authentication publishing profile is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"https_only": {
			Description: "Whether the Logic App can only be accessed via HTTPS.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"scm_publish_basic_authentication_enabled": {
			Description: "Whether the default SCM basic authentication publishing profile is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"site_config": helpers.SchemaLogicAppStandardSiteConfig(),

		"connection_string": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the Connection String.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"type": {
						Description:  "The type of the Connection String.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForConnectionStringType(), false),
					},

					"value": {
						Description: "The value for the Connection String.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"storage_account_name": {
			Description:  "The backend Storage Account name which will be used by this Logic App (e.g. for Stateful workflows data).",
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
			Description:  "The access key which will be used to access the backend Storage Account for the Logic App.",
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
			Description:  "The Key Vault Secret ID, optionally including version, that contains the connection string to the backend Storage Account for the Logic App.",
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
			Description:  "The User Assigned Identity ID used for accessing KeyVault secrets.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access": {
			Description: "Whether Public Network Access should be enabled.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     helpers.PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice([]string{
				helpers.PublicNetworkAccessEnabled,
				helpers.PublicNetworkAccessDisabled,
			}, false),
		},

		"storage_account_share_name": {
			Description: "The name of the share used by the Logic App, if you want to use a custom name. This corresponds to the WEBSITE_CONTENTSHARE appsetting, which this resource will create for you. If you don't specify a name, then this resource will generate a dynamic name. This setting is useful if you want to provision a Storage Account and create a share using `azurerm_storage_share`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"version": {
			Description: "The runtime version associated with the Logic App.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
		},

		"vnet_content_share_enabled": {
			Description: "Whether allow routing traffic between the Logic App and Storage Account content share through a Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"content_share_force_enabled": {
			Description: "Should the content share settings (`WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` and `WEBSITE_CONTENTSHARE`) be enabled? Defaults to `true`. Set to `false` when deploying the Logic App to an App Service Environment (ASE) that uses the internal storage account. Note: when set to `false`, any value supplied for `storage_account_share_name` is ignored.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this resource for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
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
			Description: "An identifier used by App Service to perform domain ownership verification via DNS TXT record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Description: "The Logic App kind.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"site_credential": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Description: "The username which can be used to publish to this App Service.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"password": {
						Description: "The password associated with the username, which can be used to publish to this App Service.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
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
