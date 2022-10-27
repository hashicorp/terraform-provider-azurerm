package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentDataSource struct{}

type ContainerAppEnvironmentDataSourceModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`

	Location                    string                 `tfschema:"location"`
	LogAnalyticsWorkspaceName   string                 `tfschema:"log_analytics_workspace_name"`
	ControlPlaneSubnetId        string                 `tfschema:"control_plane_subnet_id"`
	AppsSubnetId                string                 `tfschema:"apps_subnet_id"`
	InternalLoadBalancerEnabled bool                   `tfschema:"internal_load_balancer_enabled"`
	Tags                        map[string]interface{} `tfschema:"tags"`

	DefaultDomain         string `tfschema:"default_domain"`
	DockerBridgeCidr      string `tfschema:"docker_bridge_cidr"`
	PlatformReservedCidr  string `tfschema:"platform_reserved_cidr"`
	PlatformReservedDnsIP string `tfschema:"platform_reserved_dns_ip"`
	StaticIP              string `tfschema:"static_ip"`

	// System Data - R/O
	CreatedAt          string `tfschema:"created_at"`
	CreatedBy          string `tfschema:"created_by"`
	CreatedByType      string `tfschema:"created_by_type"`
	LastModifiedAt     string `tfschema:"last_modified_at"`
	LastModifiedBy     string `tfschema:"last_modified_by"`
	LastModifiedByType string `tfschema:"last_modified_by_type"`
}

var _ sdk.DataSource = ContainerAppEnvironmentDataSource{}

func (r ContainerAppEnvironmentDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentDataSourceModel{}
}

func (r ContainerAppEnvironmentDataSource) ResourceType() string {
	return "azurerm_container_app_environment"
}

func (r ContainerAppEnvironmentDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty, // There are no meaningful indicators for validation here, seems any non-empty string is valid at the Portal?
			Description:  "The name of the Container Apps Managed Environment.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r ContainerAppEnvironmentDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"log_analytics_workspace_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Log Analytics Workspace this Container Apps Managed Environment is linked to.",
		},

		"control_plane_subnet_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The existing Subnet in use by the Container Apps Control Plane.",
		},

		"apps_subnet_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The existing Subnet in use by the Container Apps runtime.",
		},

		"internal_load_balancer_enabled": {
			Type:        pluginsdk.TypeBool,
			Computed:    true,
			Description: "Does the Container Environment operate in Internal Load Balancing Mode?",
		},

		"tags": commonschema.TagsDataSource(),

		"created_at": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The time and date at which this Container App Environment was created.",
		},

		"created_by": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The user or principal which created this Container App Environment.",
		},

		"created_by_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The type of account which created this Container App Environment.",
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

		"last_modified_at": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The time and date at which this Container App Environment was last modified.",
		},

		"last_modified_by": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The user or principal which last modified this Container App Environment.",
		},

		"last_modified_by_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The type of account which last modified this Container App Environment.",
		},

		"platform_reserved_cidr": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The IP range, in CIDR notation, that is reserved for environment infrastructure IP addresses.",
		},

		"platform_reserved_dns_ip": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The IP address from the IP range defined by `platform_reserved_cidr` that is reserved for the internal DNS server.",
		},

		"static_ip": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Static IP of the Environment.",
		},
	}
}

func (r ContainerAppEnvironmentDataSource) Read() sdk.ResourceFunc {
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
				environment.Name = id.EnvironmentName
				environment.ResourceGroup = id.ResourceGroupName
				environment.Location = location.Normalize(model.Location)
				environment.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if vnet := props.VnetConfiguration; vnet != nil {
						environment.ControlPlaneSubnetId = pointer.From(vnet.InfrastructureSubnetId)
						environment.AppsSubnetId = pointer.From(vnet.RuntimeSubnetId)
						environment.InternalLoadBalancerEnabled = pointer.From(vnet.Internal)
						environment.DockerBridgeCidr = pointer.From(vnet.DockerBridgeCidr)
						environment.PlatformReservedCidr = pointer.From(vnet.PlatformReservedCidr)
						environment.PlatformReservedDnsIP = pointer.From(vnet.PlatformReservedDnsIP)
					}

					if appsLogs := props.AppLogsConfiguration; appsLogs != nil && appsLogs.LogAnalyticsConfiguration != nil {
						environment.LogAnalyticsWorkspaceName = pointer.From(appsLogs.LogAnalyticsConfiguration.CustomerId)
					}

					environment.StaticIP = pointer.From(props.StaticIP)
					environment.DefaultDomain = pointer.From(props.DefaultDomain)
				}

				if sysData := model.SystemData; sysData != nil {
					environment.CreatedAt = sysData.CreatedAt
					environment.CreatedBy = sysData.CreatedBy
					environment.CreatedByType = sysData.CreatedByType
					environment.LastModifiedAt = sysData.LastModifiedAt
					environment.LastModifiedBy = sysData.LastModifiedBy
					environment.LastModifiedByType = sysData.LastModifiedbyType
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&environment)
		},
	}
}
