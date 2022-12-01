package mobilenetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/packetcorecontrolplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	edgedevicevalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreControlPlaneModel struct {
	Name                        string                                 `tfschema:"name"`
	ResourceGroupName           string                                 `tfschema:"resource_group_name"`
	ControlPlaneAccessInterface []InterfacePropertiesModel             `tfschema:"control_plane_access_interface"`
	CoreNetworkTechnology       packetcorecontrolplane.CoreNetworkType `tfschema:"core_network_technology"`
	LocalDiagnosticsAccessUrl   string                                 `tfschema:"local_diagnostics_access_certificate_url"`
	Location                    string                                 `tfschema:"location"`
	MobileNetworkId             string                                 `tfschema:"mobile_network_id"`
	Platform                    []PlatformConfigurationModel           `tfschema:"platform"`
	Sku                         packetcorecontrolplane.BillingSku      `tfschema:"sku"`
	InteropSettings             string                                 `tfschema:"interop_settings"`
	Tags                        map[string]string                      `tfschema:"tags"`
	Version                     string                                 `tfschema:"version"`
}

type PlatformConfigurationModel struct {
	AzureStackEdgeDeviceId string                              `tfschema:"edge_device_id"`
	ConnectedClusterId     string                              `tfschema:"connected_cluster_id"`
	CustomLocationId       string                              `tfschema:"custom_location_id"`
	Type                   packetcorecontrolplane.PlatformType `tfschema:"type"`
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

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: mobilenetwork.ValidateMobileNetworkID,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(packetcorecontrolplane.BillingSkuEdgeSiteFourGBPS),
				string(packetcorecontrolplane.BillingSkuMediumPackage),
				string(packetcorecontrolplane.BillingSkuLargePackage),
				string(packetcorecontrolplane.BillingSkuEvaluationPackage),
				string(packetcorecontrolplane.BillingSkuFlagshipStarterPackage),
				string(packetcorecontrolplane.BillingSkuEdgeSiteTwoGBPS),
				string(packetcorecontrolplane.BillingSkuEdgeSiteThreeGBPS),
			}, false),
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: edgedevicevalidate.DeviceID,
					},

					"connected_cluster_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"custom_location_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(packetcorecontrolplane.PlatformTypeAKSNegativeHCI),
							string(packetcorecontrolplane.PlatformTypeBaseVM),
						}, false),
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),
		//it's still in progress, And will only support user assigned identity.

		"interop_settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"local_diagnostics_access_certificate_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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
		Timeout: 30 * time.Minute,
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

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties := &packetcorecontrolplane.PacketCoreControlPlane{
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Properties: packetcorecontrolplane.PacketCoreControlPlanePropertiesFormat{
					CoreNetworkTechnology: &model.CoreNetworkTechnology,
					Sku:                   model.Sku,
					MobileNetwork: packetcorecontrolplane.MobileNetworkResourceId{
						Id: model.MobileNetworkId,
					},
				},
				Tags: &model.Tags,
			}

			controlPlaneAccessInterfaceValue, err := expandPacketCoreControlPlaneInterfacePropertiesModel(model.ControlPlaneAccessInterface)
			if err != nil {
				return err
			}

			if controlPlaneAccessInterfaceValue != nil {
				properties.Properties.ControlPlaneAccessInterface = *controlPlaneAccessInterfaceValue
			}

			if model.InteropSettings != "" {
				var interopSettingsValue interface{}
				err = json.Unmarshal([]byte(model.InteropSettings), &interopSettingsValue)
				if err != nil {
					return err
				}
				properties.Properties.InteropSettings = &interopSettingsValue
			}

			if model.LocalDiagnosticsAccessUrl != "" {
				properties.Properties.LocalDiagnosticsAccess = &packetcorecontrolplane.LocalDiagnosticsAccessConfiguration{
					HTTPSServerCertificate: &packetcorecontrolplane.KeyVaultCertificate{
						CertificateUrl: &model.LocalDiagnosticsAccessUrl,
					},
				}
			}

			platformValue, err := expandPlatformConfigurationModel(model.Platform)
			if err != nil {
				return err
			}

			properties.Properties.Platform = platformValue

			if model.Version != "" {
				properties.Properties.Version = &model.Version
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
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

			var model PacketCoreControlPlaneModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("control_plane_access_interface") {
				controlPlaneAccessInterfaceValue, err := expandPacketCoreControlPlaneInterfacePropertiesModel(model.ControlPlaneAccessInterface)
				if err != nil {
					return err
				}

				if controlPlaneAccessInterfaceValue != nil {
					properties.Properties.ControlPlaneAccessInterface = *controlPlaneAccessInterfaceValue
				}
			}

			if metadata.ResourceData.HasChange("core_network_technology") {
				properties.Properties.CoreNetworkTechnology = &model.CoreNetworkTechnology
			}

			if metadata.ResourceData.HasChange("interop_settings") {
				var interopSettingsValue interface{}
				err := json.Unmarshal([]byte(model.InteropSettings), &interopSettingsValue)
				if err != nil {
					return err
				}

				properties.Properties.InteropSettings = &interopSettingsValue
			}

			if metadata.ResourceData.HasChange("local_diagnostics_access") {
				properties.Properties.LocalDiagnosticsAccess = &packetcorecontrolplane.LocalDiagnosticsAccessConfiguration{
					HTTPSServerCertificate: &packetcorecontrolplane.KeyVaultCertificate{
						CertificateUrl: &model.LocalDiagnosticsAccessUrl,
					},
				}
			}

			if metadata.ResourceData.HasChange("mobile_network_id") {
				properties.Properties.MobileNetwork = packetcorecontrolplane.MobileNetworkResourceId{
					Id: model.MobileNetworkId,
				}
			}

			if metadata.ResourceData.HasChange("platform") {
				platformValue, err := expandPlatformConfigurationModel(model.Platform)
				if err != nil {
					return err
				}

				properties.Properties.Platform = platformValue
			}

			if metadata.ResourceData.HasChange("sku") {
				properties.Properties.Sku = model.Sku
			}

			if metadata.ResourceData.HasChange("version") {
				if model.Version != "" {
					properties.Properties.Version = &model.Version
				} else {
					properties.Properties.Version = nil
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PacketCoreControlPlaneModel{
				Name:              id.PacketCoreControlPlaneName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties
			controlPlaneAccessInterfaceValue, err := flattenPacketCoreControlPlaneInterfacePropertiesModel(&properties.ControlPlaneAccessInterface)
			if err != nil {
				return err
			}

			state.ControlPlaneAccessInterface = controlPlaneAccessInterfaceValue

			if properties.CoreNetworkTechnology != nil {
				state.CoreNetworkTechnology = *properties.CoreNetworkTechnology
			}

			if properties.InteropSettings != nil && *properties.InteropSettings != nil {

				interopSettingsValue, err := json.Marshal(*properties.InteropSettings)
				if err != nil {
					return err
				}

				state.InteropSettings = string(interopSettingsValue)
			}

			if properties.LocalDiagnosticsAccess != nil && properties.LocalDiagnosticsAccess.HTTPSServerCertificate != nil && properties.LocalDiagnosticsAccess.HTTPSServerCertificate.CertificateUrl != nil {
				state.LocalDiagnosticsAccessUrl = *properties.LocalDiagnosticsAccess.HTTPSServerCertificate.CertificateUrl
			}

			state.MobileNetworkId = properties.MobileNetwork.Id

			platformValue, err := flattenPlatformConfigurationModel(properties.Platform)
			if err != nil {
				return err
			}

			state.Platform = platformValue

			state.Sku = properties.Sku

			if properties.Version != nil {
				state.Version = *properties.Version
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PacketCoreControlPlaneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.PacketCoreControlPlaneClient

			id, err := packetcorecontrolplane.ParsePacketCoreControlPlaneID(metadata.ResourceData.Id())
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

func expandPacketCoreControlPlaneInterfacePropertiesModel(inputList []InterfacePropertiesModel) (*packetcorecontrolplane.InterfaceProperties, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := packetcorecontrolplane.InterfaceProperties{}

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

	return &output, nil
}

func expandPlatformConfigurationModel(inputList []PlatformConfigurationModel) (*packetcorecontrolplane.PlatformConfiguration, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := packetcorecontrolplane.PlatformConfiguration{
		Type: input.Type,
	}

	if input.Type == packetcorecontrolplane.PlatformTypeAKSNegativeHCI && input.AzureStackEdgeDeviceId == "" {
		return nil, fmt.Errorf("`edge_device_id` must be specified when `type` is `AKS-HCI`")
	}

	if input.Type == packetcorecontrolplane.PlatformTypeBaseVM && input.AzureStackEdgeDeviceId != "" {
		return nil, fmt.Errorf("`edge_device_id` must not be specified when `type` is `BaseVM`")
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

	if input.CustomLocationId != "" {
		output.CustomLocation = &packetcorecontrolplane.CustomLocationResourceId{
			Id: input.CustomLocationId,
		}
	}

	return &output, nil
}

func flattenPacketCoreControlPlaneInterfacePropertiesModel(input *packetcorecontrolplane.InterfaceProperties) ([]InterfacePropertiesModel, error) {
	var outputList []InterfacePropertiesModel
	if input == nil {
		return outputList, nil
	}

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

	return append(outputList, output), nil
}

func flattenPlatformConfigurationModel(input *packetcorecontrolplane.PlatformConfiguration) ([]PlatformConfigurationModel, error) {
	var outputList []PlatformConfigurationModel
	if input == nil {
		return outputList, nil
	}

	output := PlatformConfigurationModel{
		Type: input.Type,
	}

	if input.AzureStackEdgeDevice != nil {
		output.AzureStackEdgeDeviceId = input.AzureStackEdgeDevice.Id
	}

	if input.ConnectedCluster != nil {
		output.ConnectedClusterId = input.ConnectedCluster.Id
	}

	if input.CustomLocation != nil {
		output.CustomLocationId = input.CustomLocation.Id
	}

	return append(outputList, output), nil
}
