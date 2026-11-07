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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentName,
			Description:  "The name of the Container Apps Managed Environment.",
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"dapr_application_insights_connection_string": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Application Insights connection string used by Dapr to export Service to Service communication telemetry.",
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
			Description:  "The ID for the Log Analytics Workspace to link this Container Apps Managed Environment to.",
		},

		"logs_destination": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  LogsDestinationNone,
			ValidateFunc: validation.StringInSlice([]string{
				LogsDestinationAzureMonitor,
				LogsDestinationLogAnalytics,
			}, false),
			Description: "The destination for the application logs. Possible values include `log-analytics` and `azure-monitor`.  Omitting this value will result in logs being streamed only.",
		},

		"infrastructure_resource_group_name": {
			Type:                  pluginsdk.TypeString,
			Optional:              true,
			ForceNew:              true,
			RequiredWith:          []string{"workload_profile"},
			ValidateFunc:          resourcegroups.ValidateName,
			DiffSuppressOnRefresh: true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *pluginsdk.ResourceData) bool { // If this is omitted and workload_profile is set, then the service generates a value for the required manage resource group.
				if profiles := d.Get("workload_profile").(*pluginsdk.Set).List(); len(profiles) > 0 && newValue == "" {
					return true
				}
				return false
			},
			Description: "Name of the platform-managed resource group created for the Managed Environment to host infrastructure resources. **Note:** Only valid if a `workload_profile` is specified. If `infrastructure_subnet_id` is specified, this resource group will be created in the same subscription as `infrastructure_subnet_id`.",
		},

		"infrastructure_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
			Description:  "The existing Subnet to use for the Container Apps Control Plane. **NOTE:** The Subnet must have a `/21` or larger address space.",
		},

		"internal_load_balancer_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			RequiredWith: []string{"infrastructure_subnet_id"},
			Description:  "Should the Container Environment operate in Internal Load Balancing Mode? Defaults to `false`. **Note:** can only be set to `true` if `infrastructure_subnet_id` is specified.",
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure sets a different value depending on whether the CAE is provisioned with its own virtual network.
			Computed:     true,
			Description:  "The public network access setting for the Container App Environment.",
			ValidateFunc: validation.StringInSlice(managedenvironments.PossibleValuesForPublicNetworkAccess(), false),
		},

		"workload_profile": WorkloadProfileSchema(),

		"zone_redundancy_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			RequiredWith: []string{"infrastructure_subnet_id"},
		},

		"mutual_tls_enabled": {
			Description: "Should mutual transport layer security (mTLS) be enabled? Defaults to `false`. **Note:** This feature is in public preview. Enabling mTLS for your applications may increase response latency and reduce maximum throughput in high-load scenarios.",
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
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The ID of the Custom Domain Verification for this Container App Environment.",
		},

		"default_domain": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The default publicly resolvable name of this Container App Environment",
		},

		"docker_bridge_cidr": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The network addressing in which the Container Apps in this Container App Environment will reside in CIDR notation.",
		},

		"platform_reserved_cidr": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.",
		},

		"platform_reserved_dns_ip_address": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.",
		},

		"static_ip_address": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Static IP Address of the Environment.",
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

			return nil
		},
	}
}
