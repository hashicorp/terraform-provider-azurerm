package mobilenetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2020-12-01/devices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreControlPlaneModel struct {
	Name                        string                                     `tfschema:"name"`
	ResourceGroupName           string                                     `tfschema:"resource_group_name"`
	ControlPlaneAccessInterface []InterfacePropertiesModel                 `tfschema:"control_plane_access_interface"`
	CoreNetworkTechnology       string                                     `tfschema:"core_network_technology"`
	LocalDiagnosticsAccess      []LocalDiagnosticsAccessConfigurationModel `tfschema:"local_diagnostics_access"`
	Location                    string                                     `tfschema:"location"`
	SiteIds                     []string                                   `tfschema:"site_ids"`
	Platform                    []PlatformConfigurationModel               `tfschema:"platform"`
	Sku                         string                                     `tfschema:"sku"`
	UeMtu                       int64                                      `tfschema:"user_equipment_mtu_in_bytes"`
	InteropSettings             string                                     `tfschema:"interop_json"`
	Identity                    []identity.ModelUserAssigned               `tfschema:"identity"`
	Tags                        map[string]string                          `tfschema:"tags"`
	Version                     string                                     `tfschema:"version"`
}

type LocalDiagnosticsAccessConfigurationModel struct {
	AuthenticationType        string `tfschema:"authentication_type"`
	HttpsServerCertificateUrl string `tfschema:"https_server_certificate_url"`
}

type PlatformConfigurationModel struct {
	AzureStackEdgeDeviceId string `tfschema:"edge_device_id"`
	AzureStackHciClusterId string `tfschema:"azure_stack_hci_cluster_id"`
	ConnectedClusterId     string `tfschema:"azure_arc_connected_cluster_id"`
	CustomLocationId       string `tfschema:"custom_location_id"`
	Type                   string `tfschema:"type"`
}

type PacketCoreControlPlaneResource struct{}

var _ sdk.ResourceWithUpdate = PacketCoreControlPlaneResource{}

func (r PacketCoreControlPlaneResource) ResourceType() string {
	return "azurerm_mobile_network_packet_core_control_plane"
}

func (r PacketCoreControlPlaneResource) ModelObject() interface{} {
	return &PacketCoreControlPlaneModel{}
}

func (r PacketCoreControlPlaneResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packetcorecontrolplane.ValidatePacketCoreControlPlaneID
}

func (r PacketCoreControlPlaneResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"control_plane_access_interface": interfacePropertiesSchema(),

		"site_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: site.ValidateSiteID,
			},
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(packetcorecontrolplane.BillingSkuGZero),
				string(packetcorecontrolplane.BillingSkuGOne),
				string(packetcorecontrolplane.BillingSkuGTwo),
				string(packetcorecontrolplane.BillingSkuGThree),
				string(packetcorecontrolplane.BillingSkuGFour),
				string(packetcorecontrolplane.BillingSkuGFive),
				string(packetcorecontrolplane.BillingSkuGOneZero),
			}, false),
		},

		"user_equipment_mtu_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1440,
			ValidateFunc: validation.IntBetween(1280, 1930),
		},

		"core_network_technology": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(packetcorecontrolplane.CoreNetworkTypeFiveGC),
				string(packetcorecontrolplane.CoreNetworkTypeEPC),
			}, false),
		},

		"platform": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"edge_device_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						AtLeastOneOf:     []string{"platform.0.edge_device_id", "platform.0.azure_stack_hci_cluster_id", "platform.0.azure_arc_connected_cluster_id", "platform.0.custom_location_id"},

						ValidateFunc: func(i interface{}, s string) ([]string, []error) {
							// a workaround, sometimes the id from user is with `DataBoxEdgeDevices` instead of `dataBoxEdgeDevices`.
							// while the service accepts both format.
							if v, ok := i.(string); ok {
								i = strings.ReplaceAll(v, "DataBoxEdgeDevices", "dataBoxEdgeDevices")
							}
							return devices.ValidateDataBoxEdgeDeviceID(i, s)
						},
					},

					"azure_stack_hci_cluster_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"platform.0.edge_device_id", "platform.0.azure_stack_hci_cluster_id", "platform.0.azure_arc_connected_cluster_id", "platform.0.custom_location_id"},
					},

					"azure_arc_connected_cluster_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"platform.0.edge_device_id", "platform.0.azure_stack_hci_cluster_id", "platform.0.azure_arc_connected_cluster_id", "platform.0.custom_location_id"},
					},

					"custom_location_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"platform.0.edge_device_id", "platform.0.azure_stack_hci_cluster_id", "platform.0.azure_arc_connected_cluster_id", "platform.0.custom_location_id"},
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(packetcorecontrolplane.PlatformTypeAKSNegativeHCI),
							string(packetcorecontrolplane.PlatformTypeThreePNegativeAZURENegativeSTACKNegativeHCI),
							string("BaseVM"), // a workaround for there might be existing one with this value. allow user to import them.
						}, false),
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"interop_json": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"local_diagnostics_access": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(packetcorecontrolplane.AuthenticationTypeAAD),
							string(packetcorecontrolplane.AuthenticationTypePassword),
						}, false),
					},
					"https_server_certificate_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r PacketCoreControlPlaneResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PacketCoreControlPlaneResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PacketCoreControlPlaneModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := packetcorecontrolplane.NewPacketCoreControlPlaneID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := expandMobileNetworkLegacyToUserAssignedIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties := packetcorecontrolplane.PacketCoreControlPlane{
				Name:     &model.Name,
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
			}

			props := packetcorecontrolplane.PacketCoreControlPlanePropertiesFormat{
				Sku:   packetcorecontrolplane.BillingSku(model.Sku),
				Sites: expandPacketCoreControlPlaneSites(model.SiteIds),
				UeMtu: &model.UeMtu,
			}

			if model.CoreNetworkTechnology != "" {
				value := packetcorecontrolplane.CoreNetworkType(model.CoreNetworkTechnology)
				props.CoreNetworkTechnology = &value
			}

			props.ControlPlaneAccessInterface = expandPacketCoreControlPlaneInterfacePropertiesModel(model.ControlPlaneAccessInterface)

			if model.InteropSettings != "" {
				var interopSettingsValue interface{}
				err = json.Unmarshal([]byte(model.InteropSettings), &interopSettingsValue)
				if err != nil {
					return err
				}
				props.InteropSettings = &interopSettingsValue
			}

			props.LocalDiagnosticsAccess = expandPacketCoreControlLocalDiagnosticsAccessConfiguration(model.LocalDiagnosticsAccess)

			props.Platform, err = expandPlatformConfigurationModel(model.Platform)
			if err != nil {
				return err
			}

			if strings.EqualFold("basevm", string(props.Platform.Type)) {
				return fmt.Errorf("it's not allowed to create new PCCP with BaseVm platform type")
			}

			if model.Version != "" {
				props.Version = &model.Version
			}

			properties.Properties = props

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PacketCoreControlPlaneResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient

			id, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan PacketCoreControlPlaneModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: Model was nil", id)
			}
			model := *resp.Model

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := expandMobileNetworkLegacyToUserAssignedIdentity(plan.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("control_plane_access_interface") {
				model.Properties.ControlPlaneAccessInterface = expandPacketCoreControlPlaneInterfacePropertiesModel(plan.ControlPlaneAccessInterface)
			}

			if metadata.ResourceData.HasChange("core_network_technology") {
				value := packetcorecontrolplane.CoreNetworkType(plan.CoreNetworkTechnology)
				model.Properties.CoreNetworkTechnology = &value
			}

			if metadata.ResourceData.HasChange("interop_json") {
				var interopSettingsValue interface{}
				err := json.Unmarshal([]byte(plan.InteropSettings), &interopSettingsValue)
				if err != nil {
					return err
				}

				model.Properties.InteropSettings = &interopSettingsValue
			}

			if metadata.ResourceData.HasChange("local_diagnostics_access") {
				model.Properties.LocalDiagnosticsAccess = expandPacketCoreControlLocalDiagnosticsAccessConfiguration(plan.LocalDiagnosticsAccess)
			}

			if metadata.ResourceData.HasChange("mobile_network_id") {
				model.Properties.Sites = expandPacketCoreControlPlaneSites(plan.SiteIds)
			}

			if metadata.ResourceData.HasChange("platform") {
				model.Properties.Platform, err = expandPlatformConfigurationModel(plan.Platform)
				if err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("sku") {
				model.Properties.Sku = packetcorecontrolplane.BillingSku(plan.Sku)
			}

			if metadata.ResourceData.HasChange("version") && plan.Version != "" {
				model.Properties.Version = &plan.Version
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &plan.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PacketCoreControlPlaneResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient

			id, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PacketCoreControlPlaneModel{
				Name:              id.PacketCoreControlPlaneName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if model.Properties.UeMtu != nil {
					state.UeMtu = *model.Properties.UeMtu
				}

				state.Identity, err = flattenMobileNetworkUserAssignedToNetworkLegacyIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				properties := model.Properties
				state.ControlPlaneAccessInterface = flattenPacketCoreControlPlaneInterfacePropertiesModel(properties.ControlPlaneAccessInterface)

				if properties.CoreNetworkTechnology != nil {
					state.CoreNetworkTechnology = string(*properties.CoreNetworkTechnology)
				}

				if properties.InteropSettings != nil && *properties.InteropSettings != nil {
					interopSettingsValue, err := json.Marshal(*properties.InteropSettings)
					if err != nil {
						return err
					}

					state.InteropSettings = string(interopSettingsValue)
				}

				state.LocalDiagnosticsAccess = flattenLocalPacketCoreControlLocalDiagnosticsAccessConfiguration(properties.LocalDiagnosticsAccess)

				state.SiteIds = flattenPacketCoreControlPlaneSites(properties.Sites)

				state.Platform = flattenPlatformConfigurationModel(properties.Platform)

				state.Sku = string(properties.Sku)

				if properties.Version != nil {
					state.Version = *properties.Version
				}
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PacketCoreControlPlaneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient

			id, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.Get(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandPacketCoreControlPlaneSites(input []string) []packetcorecontrolplane.SiteResourceId {
	outputs := make([]packetcorecontrolplane.SiteResourceId, 0)
	for _, siteId := range input {
		outputs = append(outputs, packetcorecontrolplane.SiteResourceId{
			Id: siteId,
		})
	}
	return outputs
}

func flattenPacketCoreControlPlaneSites(input []packetcorecontrolplane.SiteResourceId) []string {
	outputs := make([]string, 0)
	for _, site := range input {
		outputs = append(outputs, site.Id)
	}
	return outputs
}

func expandPacketCoreControlLocalDiagnosticsAccessConfiguration(input []LocalDiagnosticsAccessConfigurationModel) packetcorecontrolplane.LocalDiagnosticsAccessConfiguration {
	model := input[0]
	output := packetcorecontrolplane.LocalDiagnosticsAccessConfiguration{
		AuthenticationType: packetcorecontrolplane.AuthenticationType(model.AuthenticationType),
	}
	if model.HttpsServerCertificateUrl != "" {
		output.HTTPSServerCertificate = &packetcorecontrolplane.HTTPSServerCertificate{
			CertificateUrl: model.HttpsServerCertificateUrl,
		}
	}
	return output
}

func flattenLocalPacketCoreControlLocalDiagnosticsAccessConfiguration(input packetcorecontrolplane.LocalDiagnosticsAccessConfiguration) []LocalDiagnosticsAccessConfigurationModel {
	outputs := make([]LocalDiagnosticsAccessConfigurationModel, 0)
	output := LocalDiagnosticsAccessConfigurationModel{
		AuthenticationType: string(input.AuthenticationType),
	}
	if input.HTTPSServerCertificate != nil {
		output.HttpsServerCertificateUrl = input.HTTPSServerCertificate.CertificateUrl
	}
	outputs = append(outputs, output)
	return outputs
}

func expandPacketCoreControlPlaneInterfacePropertiesModel(inputList []InterfacePropertiesModel) packetcorecontrolplane.InterfaceProperties {
	output := packetcorecontrolplane.InterfaceProperties{}
	if len(inputList) == 0 {
		return output
	}

	input := inputList[0]

	if input.IPv4Address != "" {
		output.IPv4Address = &input.IPv4Address
	}

	if input.IPv4Gateway != "" {
		output.IPv4Gateway = &input.IPv4Gateway
	}

	if input.IPv4Subnet != "" {
		output.IPv4Subnet = &input.IPv4Subnet
	}

	if input.Name != "" {
		output.Name = &input.Name
	}

	return output
}

func expandPlatformConfigurationModel(inputList []PlatformConfigurationModel) (packetcorecontrolplane.PlatformConfiguration, error) {
	output := packetcorecontrolplane.PlatformConfiguration{}
	if len(inputList) == 0 {
		return output, nil
	}

	input := inputList[0]

	output.Type = packetcorecontrolplane.PlatformType(input.Type)

	if pass, err := vertifyPlatformConfigurationModel(input); !pass {
		return output, err
	}

	if input.AzureStackEdgeDeviceId != "" {
		output.AzureStackEdgeDevice = &packetcorecontrolplane.AzureStackEdgeDeviceResourceId{
			Id: input.AzureStackEdgeDeviceId,
		}
	}

	if input.ConnectedClusterId != "" {
		output.ConnectedCluster = &packetcorecontrolplane.ConnectedClusterResourceId{
			Id: input.ConnectedClusterId,
		}
	}

	if input.AzureStackHciClusterId != "" {
		output.AzureStackHciCluster = &packetcorecontrolplane.AzureStackHCIClusterResourceId{
			Id: input.AzureStackHciClusterId,
		}
	}

	if input.CustomLocationId != "" {
		output.CustomLocation = &packetcorecontrolplane.CustomLocationResourceId{
			Id: input.CustomLocationId,
		}
	}

	return output, nil
}

func vertifyPlatformConfigurationModel(input PlatformConfigurationModel) (bool, error) {
	idList := make([]string, 0)
	if input.AzureStackEdgeDeviceId != "" {
		idList = append(idList, input.AzureStackEdgeDeviceId)
	}
	if input.AzureStackHciClusterId != "" {
		idList = append(idList, input.AzureStackHciClusterId)
	}
	if input.ConnectedClusterId != "" {
		idList = append(idList, input.ConnectedClusterId)
	}

	if len(idList) == 0 {
		return false, fmt.Errorf("at least one of `azure_arc_connected_cluster_id`, `azure_stack_hci_cluster_id` and `custom_location_id` should be specified")
	}

	firstId := idList[0]
	for _, id := range idList {
		if !strings.EqualFold(firstId, id) {
			return false, fmt.Errorf("if multiple of `azure_arc_connected_cluster_id`, `azure_stack_hci_cluster_id` and `custom_location_id` are specified, they should be consistent with each other")
		}
	}

	return true, nil
}

func flattenPacketCoreControlPlaneInterfacePropertiesModel(input packetcorecontrolplane.InterfaceProperties) []InterfacePropertiesModel {
	var outputList []InterfacePropertiesModel

	output := InterfacePropertiesModel{}

	if input.IPv4Address != nil {
		output.IPv4Address = *input.IPv4Address
	}

	if input.IPv4Gateway != nil {
		output.IPv4Gateway = *input.IPv4Gateway
	}

	if input.IPv4Subnet != nil {
		output.IPv4Subnet = *input.IPv4Subnet
	}

	if input.Name != nil {
		output.Name = *input.Name
	}

	outputList = append(outputList, output)
	return outputList
}

func flattenPlatformConfigurationModel(input packetcorecontrolplane.PlatformConfiguration) []PlatformConfigurationModel {
	var outputList []PlatformConfigurationModel

	output := PlatformConfigurationModel{
		Type: string(input.Type),
	}

	if input.AzureStackEdgeDevice != nil {
		output.AzureStackEdgeDeviceId = input.AzureStackEdgeDevice.Id
	}

	if input.AzureStackHciCluster != nil {
		output.AzureStackHciClusterId = input.AzureStackHciCluster.Id
	}

	if input.ConnectedCluster != nil {
		output.ConnectedClusterId = input.ConnectedCluster.Id
	}

	if input.CustomLocation != nil {
		output.CustomLocationId = input.CustomLocation.Id
	}

	return append(outputList, output)
}
