// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_job

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerAppJobName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"container_app_environment_id": commonschema.ResourceIDReferenceRequiredForceNew(&certificates.ManagedEnvironmentId{}),

		"replica_timeout_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"workload_profile_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"template": helpers.JobTemplateSchema(),

		"secret": helpers.SecretsSchema(),

		"replica_retry_limit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"registry": helpers.ContainerAppRegistrySchema(),

		"event_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"scale": helpers.ContainerAppsJobsScaleSchema(),
				},
			},
		},

		"schedule_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cron_expression": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"manual_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"event_stream_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
