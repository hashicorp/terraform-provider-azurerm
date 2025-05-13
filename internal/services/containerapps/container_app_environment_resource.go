// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	LogsDestinationLogAnalytics string = "log-analytics"
	LogsDestinationAzureMonitor string = "azure-monitor"
	LogsDestinationNone         string = ""
)

type ContainerAppEnvironmentResource struct{}

type ContainerAppEnvironmentModel struct {
	Name                                    string                         `tfschema:"name"`
	ResourceGroup                           string                         `tfschema:"resource_group_name"`
	Location                                string                         `tfschema:"location"`
	DaprApplicationInsightsConnectionString string                         `tfschema:"dapr_application_insights_connection_string"`
	LogAnalyticsWorkspaceId                 string                         `tfschema:"log_analytics_workspace_id"`
	LogsDestination                         string                         `tfschema:"logs_destination"`
	InfrastructureSubnetId                  string                         `tfschema:"infrastructure_subnet_id"`
	InternalLoadBalancerEnabled             bool                           `tfschema:"internal_load_balancer_enabled"`
	ZoneRedundant                           bool                           `tfschema:"zone_redundancy_enabled"`
	Tags                                    map[string]interface{}         `tfschema:"tags"`
	WorkloadProfiles                        []helpers.WorkloadProfileModel `tfschema:"workload_profile"`
	InfrastructureResourceGroup             string                         `tfschema:"infrastructure_resource_group_name"`
	Mtls                                    bool                           `tfschema:"mutual_tls_enabled"`

	CustomDomainVerificationId string `tfschema:"custom_domain_verification_id"`

	DefaultDomain         string `tfschema:"default_domain"`
	DockerBridgeCidr      string `tfschema:"docker_bridge_cidr"`
	PlatformReservedCidr  string `tfschema:"platform_reserved_cidr"`
	PlatformReservedDnsIP string `tfschema:"platform_reserved_dns_ip_address"`
	StaticIP              string `tfschema:"static_ip_address"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentResource{}

var _ sdk.ResourceWithCustomizeDiff = ContainerAppEnvironmentResource{}

func (r ContainerAppEnvironmentResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentModel{}
}

func (r ContainerAppEnvironmentResource) ResourceType() string {
	return "azurerm_container_app_environment"
}

func (r ContainerAppEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedenvironments.ValidateManagedEnvironmentID
}

func (r ContainerAppEnvironmentResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
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
			DiffSuppressFunc: func(k, oldValue, newValue string, d *pluginsdk.ResourceData) bool { // If this is omitted, and there is a non-consumption profile, then the service generates a value for the required manage resource group.
				if profiles := d.Get("workload_profile").(*pluginsdk.Set).List(); len(profiles) > 0 && newValue == "" {
					for _, profile := range profiles {
						if profile.(map[string]interface{})["workload_profile_type"].(string) != string(helpers.WorkloadProfileSkuConsumption) {
							return true
						}
					}
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

		"workload_profile": helpers.WorkloadProfileSchema(),

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

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		schema["logs_destination"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // NOTE: O+C as the introduction of this property is a behavioural change where we previously set it behind the scenes if `log_analytics_workspace_id` was set.
			ValidateFunc: validation.StringInSlice([]string{
				LogsDestinationAzureMonitor,
				LogsDestinationNone,
				LogsDestinationLogAnalytics,
			}, false),
		}
	}

	return schema
}

func (r ContainerAppEnvironmentResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r ContainerAppEnvironmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var containerAppEnvironment ContainerAppEnvironmentModel

			if err := metadata.Decode(&containerAppEnvironment); err != nil {
				return err
			}

			id := managedenvironments.NewManagedEnvironmentID(subscriptionId, containerAppEnvironment.ResourceGroup, containerAppEnvironment.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			managedEnvironment := managedenvironments.ManagedEnvironment{
				Location: containerAppEnvironment.Location,
				Name:     pointer.To(containerAppEnvironment.Name),
				Properties: &managedenvironments.ManagedEnvironmentProperties{
					AppLogsConfiguration: &managedenvironments.AppLogsConfiguration{
						Destination: pointer.To(containerAppEnvironment.LogsDestination),
					},
					VnetConfiguration: &managedenvironments.VnetConfiguration{},
					ZoneRedundant:     pointer.To(containerAppEnvironment.ZoneRedundant),
					PeerAuthentication: &managedenvironments.ManagedEnvironmentPropertiesPeerAuthentication{
						Mtls: &managedenvironments.Mtls{
							Enabled: pointer.To(containerAppEnvironment.Mtls),
						},
					},
					PeerTrafficConfiguration: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfiguration{
						Encryption: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfigurationEncryption{
							Enabled: pointer.To(containerAppEnvironment.Mtls),
						},
					},
				},
				Tags: tags.Expand(containerAppEnvironment.Tags),
			}

			if containerAppEnvironment.DaprApplicationInsightsConnectionString != "" {
				managedEnvironment.Properties.DaprAIConnectionString = pointer.To(containerAppEnvironment.DaprApplicationInsightsConnectionString)
			}

			if containerAppEnvironment.InfrastructureResourceGroup != "" {
				managedEnvironment.Properties.InfrastructureResourceGroup = pointer.To(containerAppEnvironment.InfrastructureResourceGroup)
			}

			if containerAppEnvironment.LogAnalyticsWorkspaceId != "" {
				customerId, sharedKey, err := getSharedKeyForWorkspace(ctx, metadata, containerAppEnvironment.LogAnalyticsWorkspaceId)
				if err != nil {
					return fmt.Errorf("retrieving access keys to Log Analytics Workspace for %s: %+v", id, err)
				}

				managedEnvironment.Properties.AppLogsConfiguration.Destination = pointer.To(LogsDestinationLogAnalytics)

				managedEnvironment.Properties.AppLogsConfiguration.LogAnalyticsConfiguration = &managedenvironments.LogAnalyticsConfiguration{
					CustomerId: customerId,
					SharedKey:  sharedKey,
				}
			}

			if containerAppEnvironment.InfrastructureSubnetId != "" {
				managedEnvironment.Properties.VnetConfiguration.InfrastructureSubnetId = pointer.To(containerAppEnvironment.InfrastructureSubnetId)
				managedEnvironment.Properties.VnetConfiguration.Internal = pointer.To(containerAppEnvironment.InternalLoadBalancerEnabled)
			}

			managedEnvironment.Properties.WorkloadProfiles = helpers.ExpandWorkloadProfiles(containerAppEnvironment.WorkloadProfiles)

			if err := client.CreateOrUpdateThenPoll(ctx, id, managedEnvironment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerAppEnvironmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentModel

			consumptionDefined := consumptionIsExplicitlyDefined(metadata)

			if model := existing.Model; model != nil {
				state.Name = id.ManagedEnvironmentName
				state.ResourceGroup = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if vnet := props.VnetConfiguration; vnet != nil {
						state.InfrastructureSubnetId = pointer.From(vnet.InfrastructureSubnetId)
						state.InternalLoadBalancerEnabled = pointer.From(vnet.Internal)
						state.DockerBridgeCidr = pointer.From(vnet.DockerBridgeCidr)
						state.PlatformReservedCidr = pointer.From(vnet.PlatformReservedCidr)
						state.PlatformReservedDnsIP = pointer.From(vnet.PlatformReservedDnsIP)
					}

					if appLogsConfig := props.AppLogsConfiguration; appLogsConfig != nil {
						state.LogsDestination = pointer.From(appLogsConfig.Destination)
						if appLogsConfig.LogAnalyticsConfiguration != nil && appLogsConfig.LogAnalyticsConfiguration.CustomerId != nil {
							workspaceId, err := findWorkspaceResourceIDFromCustomerID(ctx, metadata, *appLogsConfig.LogAnalyticsConfiguration.CustomerId)
							if err != nil {
								if v := metadata.ResourceData.GetRawConfig().AsValueMap()["log_analytics_workspace_id"]; !v.IsNull() && v.AsString() != "" {
									state.LogAnalyticsWorkspaceId = v.AsString()
								} else {
									return fmt.Errorf("retrieving Log Analytics Workspace ID for %s: %+v", *appLogsConfig.LogAnalyticsConfiguration.CustomerId, err)
								}
							}

							state.LogAnalyticsWorkspaceId = workspaceId.ID()
						}
					}

					state.CustomDomainVerificationId = pointer.From(props.CustomDomainConfiguration.CustomDomainVerificationId)
					state.ZoneRedundant = pointer.From(props.ZoneRedundant)
					state.StaticIP = pointer.From(props.StaticIP)
					state.DefaultDomain = pointer.From(props.DefaultDomain)
					state.WorkloadProfiles = helpers.FlattenWorkloadProfiles(props.WorkloadProfiles, consumptionDefined)
					state.InfrastructureResourceGroup = pointer.From(props.InfrastructureResourceGroup)
					state.Mtls = pointer.From(props.PeerAuthentication.Mtls.Enabled)
				}
			}

			// `dapr_application_insights_connection_string` is sensitive and not returned by API
			if v := metadata.ResourceData.Get("dapr_application_insights_connection_string").(string); v != "" {
				state.DaprApplicationInsightsConnectionString = v
			}

			if err = metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			workaroundClient := azuresdkhacks.NewManagedEnvironmentWorkaroundClient(client)
			id, err := managedenvironments.ParseManagedEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppEnvironmentModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := azuresdkhacks.ManagedEnvironment{
				Name:       pointer.To(state.Name),
				Location:   state.Location,
				Properties: &azuresdkhacks.ManagedEnvironmentProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("workload_profile") {
				payload.Properties.WorkloadProfiles = helpers.ExpandWorkloadProfiles(state.WorkloadProfiles)
			}

			if metadata.ResourceData.HasChange("mutual_tls_enabled") {
				payload.Properties.PeerAuthentication = &managedenvironments.ManagedEnvironmentPropertiesPeerAuthentication{
					Mtls: &managedenvironments.Mtls{
						Enabled: pointer.To(state.Mtls),
					},
				}
				payload.Properties.PeerTrafficConfiguration = &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfiguration{
					Encryption: &managedenvironments.ManagedEnvironmentPropertiesPeerTrafficConfigurationEncryption{
						Enabled: pointer.To(state.Mtls),
					},
				}
			}

			if metadata.ResourceData.HasChanges("logs_destination", "log_analytics_workspace_id") {
				// For 4.x we need to compensate for the legacy behaviour of setting log destination based on the presence of log_analytics_workspace_id
				if !features.FivePointOh() && metadata.ResourceData.GetRawConfig().AsValueMap()["logs_destination"].IsNull() && state.LogAnalyticsWorkspaceId == "" {
					state.LogsDestination = LogsDestinationNone
				}

				switch state.LogsDestination {
				case LogsDestinationAzureMonitor:
					payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
						Destination:               pointer.To(LogsDestinationAzureMonitor),
						LogAnalyticsConfiguration: nil,
					}
				case LogsDestinationLogAnalytics:
					if state.LogAnalyticsWorkspaceId != "" {
						customerId, sharedKey, err := getSharedKeyForWorkspace(ctx, metadata, state.LogAnalyticsWorkspaceId)
						if err != nil {
							return fmt.Errorf("retrieving access keys to Log Analytics Workspace for %s: %+v", id, err)
						}

						payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
							Destination: pointer.To(LogsDestinationLogAnalytics),
							LogAnalyticsConfiguration: &managedenvironments.LogAnalyticsConfiguration{
								CustomerId: customerId,
								SharedKey:  sharedKey,
							},
						}
					}
				default:
					payload.Properties.AppLogsConfiguration = &azuresdkhacks.AppLogsConfiguration{
						Destination:               pointer.To(LogsDestinationNone),
						LogAnalyticsConfiguration: nil,
					}
				}
			}

			if err := workaroundClient.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func consumptionIsExplicitlyDefined(metadata sdk.ResourceMetaData) bool {
	config := ContainerAppEnvironmentModel{}
	if err := metadata.Decode(&config); err != nil {
		return false
	}
	for _, v := range config.WorkloadProfiles {
		if v.Name == string(helpers.WorkloadProfileSkuConsumption) {
			return true
		}
	}

	return false
}

func (r ContainerAppEnvironmentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			if metadata.ResourceDiff.HasChange("workload_profile") {
				oldProfiles, newProfiles := metadata.ResourceDiff.GetChange("workload_profile")

				oldProfileCount := oldProfiles.(*pluginsdk.Set).Len()
				newProfileCount := newProfiles.(*pluginsdk.Set).Len()
				if oldProfileCount > 0 && newProfileCount == 0 {
					if err := metadata.ResourceDiff.ForceNew("workload_profile"); err != nil {
						return err
					}
				}

				if newProfileCount > 0 && oldProfileCount == 0 {
					if err := metadata.ResourceDiff.ForceNew("workload_profile"); err != nil {
						return err
					}
				}
			}

			if !features.FivePointOh() { // in 4.x `logs_destination` is Computed due to legacy code implying destination from presence of a valid id in `log_analytics_workspace_id` so we need to check explicit config values here
				if metadata.ResourceDiff.HasChanges("logs_destination", "log_analytics_workspace_id") {
					logsDestination := metadata.ResourceDiff.Get("logs_destination").(string)
					logDestinationIsNull := metadata.ResourceDiff.GetRawConfig().AsValueMap()["logs_destination"].IsNull()
					logAnalyticsWorkspaceID := metadata.ResourceDiff.Get("log_analytics_workspace_id").(string)
					logAnalyticsWorkspaceIDIsNull := metadata.ResourceDiff.GetRawConfig().AsValueMap()["log_analytics_workspace_id"].IsNull()

					if !logDestinationIsNull || !logAnalyticsWorkspaceIDIsNull {
						switch logsDestination {
						case LogsDestinationLogAnalytics:
							if logAnalyticsWorkspaceIDIsNull {
								return fmt.Errorf("`log_analytics_workspace_id` must be set when `logs_destination` is set to `log-analytics`")
							}

						case LogsDestinationAzureMonitor, LogsDestinationNone:
							if (logAnalyticsWorkspaceID != "" || !logAnalyticsWorkspaceIDIsNull) && !logDestinationIsNull {
								return fmt.Errorf("`log_analytics_workspace_id` can only be set when `logs_destination` is set to `log-analytics` or omitted")
							}
						}
					}
				}
			} else {
				if metadata.ResourceDiff.HasChanges("logs_destination", "log_analytics_workspace_id") {
					logsDestination := metadata.ResourceDiff.Get("logs_destination").(string)
					logAnalyticsWorkspaceID := metadata.ResourceDiff.Get("log_analytics_workspace_id").(string)

					switch logsDestination {
					case LogsDestinationLogAnalytics:
						if logAnalyticsWorkspaceID == "" {
							return fmt.Errorf("`log_analytics_workspace_id` must be set when `logs_destination` is set to `log-analytics`")
						}

					case LogsDestinationAzureMonitor, LogsDestinationNone:
						if logAnalyticsWorkspaceID != "" {
							return fmt.Errorf("`log_analytics_workspace_id` can only be set when `logs_destination` is set to `log-analytics` or omitted")
						}
					}
				}
			}

			return nil
		},
	}
}

func findWorkspaceResourceIDFromCustomerID(ctx context.Context, meta sdk.ResourceMetaData, customerID string) (*workspaces.WorkspaceId, error) {
	client := meta.Client.LogAnalytics.WorkspaceClient

	subscriptionId := commonids.NewSubscriptionID(meta.Client.Account.SubscriptionId)

	result := &workspaces.WorkspaceId{}

	list, err := client.List(ctx, subscriptionId)
	if err != nil {
		return nil, err
	}

	model := list.Model
	if model == nil {
		return nil, fmt.Errorf("could not resolve Log Analytics Workspace ID for %s, list model was nil", customerID)
	}

	if model.Value == nil || len(*model.Value) == 0 {
		return nil, fmt.Errorf("could not resolve Log Analytics Workspace ID for %s, no Log Analytics Workspaces found in %s", customerID, subscriptionId)
	}

	for _, v := range *list.Model.Value {
		if v.Properties != nil && v.Properties.CustomerId != nil && strings.EqualFold(*v.Properties.CustomerId, customerID) {
			result, err = workspaces.ParseWorkspaceIDInsensitively(pointer.From(v.Id))
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func getSharedKeyForWorkspace(ctx context.Context, meta sdk.ResourceMetaData, workspaceID string) (*string, *string, error) {
	logAnalyticsClient := meta.Client.LogAnalytics.SharedKeyWorkspacesClient

	logAnalyticsId, err := workspaces.ParseWorkspaceID(workspaceID)
	if err != nil {
		return nil, nil, err
	}

	workspace, err := logAnalyticsClient.Get(ctx, *logAnalyticsId)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving %s: %+v", logAnalyticsId, err)
	}

	if workspace.Model == nil || workspace.Model.Properties == nil {
		return nil, nil, fmt.Errorf("reading customer ID from %s", logAnalyticsId)
	}

	if workspace.Model.Properties.CustomerId == nil {
		return nil, nil, fmt.Errorf("reading customer ID from %s, `customer_id` is nil", logAnalyticsId)
	}

	keys, err := logAnalyticsClient.SharedKeysGetSharedKeys(ctx, *logAnalyticsId)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving access keys to %s: %+v", logAnalyticsId, err)
	}
	if keys.Model.PrimarySharedKey == nil {
		return nil, nil, fmt.Errorf("reading shared key for %s", logAnalyticsId)
	}

	return workspace.Model.Properties.CustomerId, keys.Model.PrimarySharedKey, nil
}
