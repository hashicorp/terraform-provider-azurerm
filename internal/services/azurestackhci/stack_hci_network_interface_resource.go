package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = StackHCINetworkInterfaceResource{}
	_ sdk.ResourceWithUpdate = StackHCINetworkInterfaceResource{}
)

type StackHCINetworkInterfaceResource struct{}

func (StackHCINetworkInterfaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkinterfaces.ValidateNetworkInterfaceID
}

func (StackHCINetworkInterfaceResource) ResourceType() string {
	return "azurerm_stack_hci_network_interface"
}

func (StackHCINetworkInterfaceResource) ModelObject() interface{} {
	return &StackHCINetworkInterfaceResourceModel{}
}

type StackHCINetworkInterfaceResourceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	CustomLocationId  string                         `tfschema:"custom_location_id"`
	DNSServers        []string                       `tfschema:"dns_servers"`
	IPConfiguration   []StackHCIIPConfigurationModel `tfschema:"ip_configuration"`
	MACAddress        string                         `tfschema:"mac_address"`
	Tags              map[string]interface{}         `tfschema:"tags"`
}

type StackHCIIPConfigurationModel struct {
	Gateway          string `tfschema:"gateway"`
	PrefixLength     string `tfschema:"prefix_length"`
	PrivateIPAddress string `tfschema:"private_ip_address"`
	SubnetID         string `tfschema:"subnet_id"`
}

func (StackHCINetworkInterfaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
				"name must be between 2 and 80 characters and can only contain alphanumberic characters, hyphen, dot and underline",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: logicalnetworks.ValidateLogicalNetworkID,
					},

					"private_ip_address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPv4Address,
					},

					"gateway": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"prefix_length": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"mac_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (StackHCINetworkInterfaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCINetworkInterfaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			var config StackHCINetworkInterfaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := networkinterfaces.NewNetworkInterfaceID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := networkinterfaces.NetworkInterfaces{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &networkinterfaces.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(networkinterfaces.ExtendedLocationTypesCustomLocation),
				},
				Properties: &networkinterfaces.NetworkInterfaceProperties{
					IPConfigurations: expandStackHCINetworkInterfaceIPConfiguration(config.IPConfiguration),
				},
			}

			if config.MACAddress != "" {
				payload.Properties.MacAddress = pointer.To(config.MACAddress)
			}

			if len(config.DNSServers) != 0 {
				payload.Properties.DnsSettings = &networkinterfaces.InterfaceDNSSettings{
					DnsServers: pointer.To(config.DNSServers),
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			time.Sleep(2 * time.Minute)
			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCINetworkInterfaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			id, err := networkinterfaces.ParseNetworkInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCINetworkInterfaceResourceModel{
				Name:              id.NetworkInterfaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.MACAddress = pointer.From(props.MacAddress)

					ipConfiguration, err := flattenStackHCINetworkInterfaceIPConfiguration(props.IPConfigurations)
					if err != nil {
						return err
					}
					schema.IPConfiguration = ipConfiguration

					if props.DnsSettings != nil {
						schema.DNSServers = pointer.From(props.DnsSettings.DnsServers)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCINetworkInterfaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			id, err := networkinterfaces.ParseNetworkInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StackHCINetworkInterfaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r StackHCINetworkInterfaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			id, err := networkinterfaces.ParseNetworkInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandStackHCINetworkInterfaceIPConfiguration(input []StackHCIIPConfigurationModel) *[]networkinterfaces.IPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]networkinterfaces.IPConfiguration, 0)
	for _, v := range input {
		result := networkinterfaces.IPConfiguration{
			Properties: &networkinterfaces.IPConfigurationProperties{
				Subnet: &networkinterfaces.IPConfigurationPropertiesSubnet{
					Id: pointer.To(v.SubnetID),
				},
			},
		}

		if v.PrivateIPAddress != "" {
			result.Properties.PrivateIPAddress = pointer.To(v.PrivateIPAddress)
		}

		results = append(results, result)
	}

	return &results
}

func flattenStackHCINetworkInterfaceIPConfiguration(input *[]networkinterfaces.IPConfiguration) ([]StackHCIIPConfigurationModel, error) {
	if input == nil {
		return make([]StackHCIIPConfigurationModel, 0), nil
	}

	results := make([]StackHCIIPConfigurationModel, 0)
	for _, v := range *input {
		result := StackHCIIPConfigurationModel{}

		if v.Properties != nil {
			var subnetId string
			if v.Properties.Subnet != nil && v.Properties.Subnet.Id != nil {
				parsedSubnetId, err := logicalnetworks.ParseLogicalNetworkIDInsensitively(*v.Properties.Subnet.Id)
				if err != nil {
					return results, err
				}

				subnetId = parsedSubnetId.ID()
			}

			result.Gateway = pointer.From(v.Properties.Gateway)
			result.PrefixLength = pointer.From(v.Properties.PrefixLength)
			result.PrivateIPAddress = pointer.From(v.Properties.PrivateIPAddress)
			result.SubnetID = subnetId

			results = append(results, result)
		}
	}

	return results, nil
}
