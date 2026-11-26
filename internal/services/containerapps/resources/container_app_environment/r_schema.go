// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Container Apps Managed Environment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"dapr_application_insights_connection_string": {
			Description:  "Application Insights connection string used by Dapr to export Service to Service communication telemetry.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Description:  "The ID for the Log Analytics Workspace to link this Container Apps Managed Environment to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"logs_destination": {
			Description: "Where the application logs will be saved for this Container Apps Managed Environment. Omitting this value will result in logs being streamed only.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     LogsDestinationNone,
			ValidateFunc: validation.StringInSlice([]string{
				LogsDestinationAzureMonitor,
				LogsDestinationLogAnalytics,
			}, false),
		},

		"infrastructure_resource_group_name": {
			Description:  "Name of the platform-managed Resource Group created for the Managed Environment to host infrastructure resources.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *pluginsdk.ResourceData) bool { // If this is omitted and workload_profile is set, then the service generates a value for the required manage resource group.
				if profiles := d.Get("workload_profile").(*pluginsdk.Set).List(); len(profiles) > 0 && newValue == "" {
					return true
				}
				return false
			},
			DiffSuppressOnRefresh: true,
			RequiredWith:          []string{"workload_profile"},
		},

		"infrastructure_subnet_id": {
			Description:  "The existing Subnet to use for the Container Apps Control Plane.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"internal_load_balancer_enabled": {
			Description:  "Whether the Container Environment should operate in Internal Load Balancing Mode.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			RequiredWith: []string{"infrastructure_subnet_id"},
		},

		"public_network_access": {
			Description: "The public network access setting for the Container App Environment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because Azure sets a different value depending on whether the CAE is provisioned with its own virtual network.
			Computed:     true,
			ValidateFunc: validation.StringInSlice(managedenvironments.PossibleValuesForPublicNetworkAccess(), false),
		},

		"workload_profile": WorkloadProfileSchema(),

		"ingress_configuration": IngressConfigurationSchema(),

		"zone_redundancy_enabled": {
			Description:  "Whether the Container App Environment should be created with Zone Redundancy enabled.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			RequiredWith: []string{"infrastructure_subnet_id"},
		},

		"mutual_tls_enabled": {
			Description: "Whether mutual transport layer security (mTLS) should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Description: "The ID of the Custom Domain Verification for this Container App Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"default_domain": {
			Description: "The default publicly resolvable name of this Container App Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"docker_bridge_cidr": {
			Description: "The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"platform_reserved_cidr": {
			Description: "The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"platform_reserved_dns_ip_address": {
			Description: "The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"static_ip_address": {
			Description: "The Static IP address of the Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ContainerAppEnvironmentModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceDiff.HasChange("workload_profile") {
				o, n := metadata.ResourceDiff.GetChange("workload_profile")

				oldProfiles := o.(*pluginsdk.Set)
				newProfiles := n.(*pluginsdk.Set)

				if oldProfiles.Len() > 0 && newProfiles.Len() == 0 && !OneAdditionalConsumptionProfileReturnedByAPI(oldProfiles, newProfiles) {
					if err := metadata.ResourceDiff.ForceNew("workload_profile"); err != nil {
						return err
					}
				}

				if newProfiles.Len() > 0 && oldProfiles.Len() == 0 {
					if err := metadata.ResourceDiff.ForceNew("workload_profile"); err != nil {
						return err
					}
				}
			}

			for _, v := range model.WorkloadProfiles {
				if IsConsumptionProfileType(v.WorkloadProfileType) && (v.MinimumCount != 0 || v.MaximumCount != 0) {
					return fmt.Errorf("`minimum_count` and `maximum_count` cannot be set for `workload_profile` %q as they are not supported for the `workload_profile_type` %q", v.Name, v.WorkloadProfileType)
				}
			}

			if metadata.ResourceDiff.HasChanges("logs_destination", "log_analytics_workspace_id") {
				logsDestination := metadata.ResourceDiff.Get("logs_destination").(string)
				logAnalyticsWorkspaceID := metadata.ResourceDiff.Get("log_analytics_workspace_id").(string)
				logAnalyticsWorkspaceIDUnknown := !metadata.ResourceDiff.GetRawConfig().AsValueMap()["log_analytics_workspace_id"].IsKnown()

				switch logsDestination {
				case LogsDestinationLogAnalytics:
					if logAnalyticsWorkspaceID == "" && !logAnalyticsWorkspaceIDUnknown {
						return errors.New("`log_analytics_workspace_id` must be set when `logs_destination` is set to `log-analytics`")
					}

				case LogsDestinationAzureMonitor, LogsDestinationNone:
					if logAnalyticsWorkspaceID != "" || logAnalyticsWorkspaceIDUnknown {
						return errors.New("`log_analytics_workspace_id` can only be set when `logs_destination` is set to `log-analytics` or `\"\"`")
					}
				}
			}

			if model.InternalLoadBalancerEnabled && model.PublicNetworkAccess == string(managedenvironments.PublicNetworkAccessEnabled) {
				return errors.New("`public_network_access` cannot be `Enabled` when `internal_load_balancer_enabled` is set to `true`")
			}

			if len(model.IngressConfiguration) > 0 {
				ic := model.IngressConfiguration[0]
				if ic.MaximumNodeCount < ic.MinimumNodeCount {
					return errors.New("`maximum_node_count` in `ingress_configuration` must be greater than or equal to `minimum_node_count`")
				}
				for _, wp := range model.WorkloadProfiles {
					if wp.Name == ic.WorkloadProfileName {
						return fmt.Errorf("`ingress_configuration.workload_profile_name` %q conflicts with a user-defined `workload_profile` block with the same name. The dedicated ingress workload profile is auto-managed by `ingress_configuration` and must not be duplicated in `workload_profile`", ic.WorkloadProfileName)
					}
				}
			}

			return nil
		},
	}
}
