// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Function App. Limit the function name to 32 characters to avoid naming collisions. For more information about [Function App naming rule](https://learn.microsoft.com/azure/azure-resource-manager/management/resource-name-rules#microsoftweb) and [Host ID Collisions](https://github.com/Azure/azure-functions-host/wiki/Host-IDs#host-id-collisions).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Description:  "The ID of the App Service Plan within which to create this Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"storage_container_type": {
			Description:  "The Storage Container type used for the Function App. The current supported type is `blobContainer`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForFunctionsDeploymentStorageType(), false),
		},

		"storage_container_endpoint": {
			Description: "The backend Storage Container endpoint which will be used by this Function App.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"storage_authentication_type": {
			Description:  "The authentication type which will be used to access the backend Storage Account for the Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAuthenticationType(), false),
		},

		"storage_access_key": {
			Description:  "The access key which will be used to access the backend Storage Account for the Function App.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_user_assigned_identity_id": {
			Description:  "The User Assigned Managed Identity to access the Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"runtime_name": {
			Description:  "The Runtime of the Linux Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForRuntimeName(), false),
		},

		"runtime_version": {
			Description:  "The Runtime version of the Linux Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"instance_memory_in_mb": {
			Description: "The memory size of the instances on which your app runs. Reference the Microsoft Documentation for the [currently supported values](https://learn.microsoft.com/azure/azure-functions/flex-consumption-plan#instance-memory).",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     2048,
		},

		"maximum_instance_count": {
			Description:  "The number of workers this Function App can scale out to.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"http_concurrency": {
			Description:  "The HTTP concurrency of the instances on which your app runs.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		// the name is always being lower-cased by the api: https://github.com/Azure/azure-rest-api-specs/issues/33095
		"always_ready": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:      "The name of the `always_ready` of the Function App.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"instance_count": {
						Description:  "The instance count of the `always_ready` of the Function App. The minimum number is `0`. The total number of `instance_count` should not exceed the `maximum_instance_count`.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
				},
			},
		},

		"site_config": helpers.SiteConfigSchemaFunctionAppFlexConsumption(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"app_settings": {
			Description: "A map of key-value pairs for [App Settings](https://learn.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"client_certificate_enabled": {
			Description: "Whether the Function App should use Client Certificates.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_certificate_mode": {
			Description:  "The mode of the Function App's client certificates requirement for incoming requests.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Description: "Paths to exclude when using client certificates, separated by `;`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Description: "Whether the Function App is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"public_network_access_enabled": {
			Description: "Should public network access be enabled for the Function App?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"https_only": {
			Description: "Whether the Function App can only be accessed via HTTPS.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Description: "Should the default WebDeploy Basic Authentication publishing credentials enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this Function App for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zip_deploy_file": {
			Description:  "The local path and filename of the Zip packaged application to deploy to this Linux Function App.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Description: "The ID of the App Service Environment used by Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Linux Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_address_list": {
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_address_list": {
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}
