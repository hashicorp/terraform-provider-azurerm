// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func springCloudAppSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceName,
		},

		"addon_json": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"custom_persistent_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"storage_name": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"mount_path": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"share_name": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"mount_options": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"read_only_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"is_public": { // azignore:AZS023 - Azure Spring Apps is being retired, so it is left as it is
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"ingress_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend_protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(appplatform.BackendProtocolDefault),
						ValidateFunc: validation.StringInEnumSlice(appplatform.PossibleBackendProtocolValues(), false),
					},

					"read_timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"send_timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      60,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"session_affinity": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(appplatform.SessionAffinityNone),
						ValidateFunc: validation.StringInEnumSlice(appplatform.PossibleSessionAffinityValues(), false),
					},

					"session_cookie_max_age": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"persistent_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size_in_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 50),
					},

					"mount_path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "/persistent",
						ValidateFunc: validate.MountPath,
					},
				},
			},
		},

		"public_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tls_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
