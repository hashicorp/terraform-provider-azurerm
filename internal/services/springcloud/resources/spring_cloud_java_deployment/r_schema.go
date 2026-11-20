// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_java_deployment

import (
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSprintCloudJavaDeploymentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudDeploymentName,
		},

		"spring_cloud_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(appplatform_rm.ParseAppIDInsensitively),
		},

		"environment_variables": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 500),
		},

		"jvm_options": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"quota": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C The value returned in GET will be recalculated by the service if the deprecated "cpu" is honoured, so make this property as Computed.
						Computed: true,
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"memory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C The value returned in GET will be recalculated by the service if the deprecated "memory_in_gb" is honoured, so make this property as Computed.
						Computed: true,
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"runtime_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appplatform.SupportedRuntimeValueJava8),
				string(appplatform.SupportedRuntimeValueJava11),
				string(appplatform.SupportedRuntimeValueJava17),
			}, false),
			Default: appplatform.SupportedRuntimeValueJava8,
		},
	}
}
