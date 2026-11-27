// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

type ContainerAppEnvironmentDataSourceModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`

	Location                    string                      `tfschema:"location"`
	LogAnalyticsWorkspaceName   string                      `tfschema:"log_analytics_workspace_name"`
	InfrastructureSubnetId      string                      `tfschema:"infrastructure_subnet_id"`
	InternalLoadBalancerEnabled bool                        `tfschema:"internal_load_balancer_enabled"`
	Tags                        map[string]any              `tfschema:"tags"`
	IngressConfiguration        []IngressConfigurationModel `tfschema:"ingress_configuration"`

	CustomDomainVerificationId string `tfschema:"custom_domain_verification_id"`

	DefaultDomain         string `tfschema:"default_domain"`
	DockerBridgeCidr      string `tfschema:"docker_bridge_cidr"`
	PlatformReservedCidr  string `tfschema:"platform_reserved_cidr"`
	PlatformReservedDnsIP string `tfschema:"platform_reserved_dns_ip_address"`
	PublicNetworkAccess   string `tfschema:"public_network_access"`
	StaticIP              string `tfschema:"static_ip_address"`
}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ModelObject() any {
	return &ContainerAppEnvironmentDataSourceModel{}
}

func (r DataSource) ResourceType() string {
	return "azurerm_container_app_environment"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Container App Environment."
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Container Apps Managed Environment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty, // There are no meaningful indicators for validation here, seems any non-empty string is valid at the Portal?
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"log_analytics_workspace_name": {
			Description: "The name of the Log Analytics Workspace this Container Apps Managed Environment is linked to.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"infrastructure_subnet_id": {
			Description: "The ID of the Subnet in use by the Container Apps Control Plane.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"internal_load_balancer_enabled": {
			Description: "Whether the Container App Environment operates in Internal Load Balancing Mode.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),

		"custom_domain_verification_id": {
			Description: "The ID of the Custom Domain Verification for this Container App Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"default_domain": {
			Description: "This is generated at creation time to be globally unique.",
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

		"public_network_access": {
			Description: "The public network access setting for this Container App Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"static_ip_address": {
			Description: "The Static IP address of the Environment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ingress_configuration": IngressConfigurationSchemaComputed(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ManagedEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var environment ContainerAppEnvironmentDataSourceModel
			if err := metadata.Decode(&environment); err != nil {
				return err
			}

			id := managedenvironments.NewManagedEnvironmentID(subscriptionId, environment.ResourceGroup, environment.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				environment.Name = id.ManagedEnvironmentName
				environment.ResourceGroup = id.ResourceGroupName
				environment.Location = location.Normalize(model.Location)
				environment.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if vnet := props.VnetConfiguration; vnet != nil {
						environment.InfrastructureSubnetId = pointer.From(vnet.InfrastructureSubnetId)
						environment.InternalLoadBalancerEnabled = pointer.From(vnet.Internal)
						environment.DockerBridgeCidr = pointer.From(vnet.DockerBridgeCidr)
						environment.PlatformReservedCidr = pointer.From(vnet.PlatformReservedCidr)
						environment.PlatformReservedDnsIP = pointer.From(vnet.PlatformReservedDnsIP)
					}

					if appsLogs := props.AppLogsConfiguration; appsLogs != nil && appsLogs.LogAnalyticsConfiguration != nil {
						lawClient := metadata.Client.LogAnalytics.SharedKeyWorkspacesClient
						lawName, err := findLogAnalyticsWorkspaceName(ctx, lawClient, subscriptionId, pointer.From(appsLogs.LogAnalyticsConfiguration.CustomerId))
						if err != nil {
							return fmt.Errorf("retrieving Log Analytics Workspace: %+v", err)
						}
						environment.LogAnalyticsWorkspaceName = lawName
					}

					environment.StaticIP = pointer.From(props.StaticIP)
					environment.DefaultDomain = pointer.From(props.DefaultDomain)
					environment.CustomDomainVerificationId = pointer.From(props.CustomDomainConfiguration.CustomDomainVerificationId)
					environment.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
					environment.IngressConfiguration = FlattenIngressConfiguration(props.IngressConfiguration, props.WorkloadProfiles)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&environment)
		},
	}
}

func findLogAnalyticsWorkspaceName(ctx context.Context, client *workspaces.WorkspacesClient, subscriptionId, targetCustomerId string) (string, error) {
	parsedSubscriptionId := commonids.NewSubscriptionID(subscriptionId)

	resp, err := client.List(ctx, parsedSubscriptionId)
	if err != nil {
		return "", err
	}

	if resp.Model == nil {
		return "", errors.New("model was nil")
	}

	if resp.Model.Value == nil {
		return "", errors.New("value was nil")
	}

	for _, law := range *resp.Model.Value {
		if law.Properties != nil && law.Properties.CustomerId != nil && *law.Properties.CustomerId == targetCustomerId && law.Name != nil {
			return *law.Name, nil
		}
	}

	return "", nil
}
