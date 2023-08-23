// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/packetcorecontrolplane"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PacketCoreControlPlaneModel struct {
	Name                          string                                     `tfschema:"name"`
	ResourceGroupName             string                                     `tfschema:"resource_group_name"`
	ControlPlaneAccessIPv4Address string                                     `tfschema:"control_plane_access_ipv4_address"`
	ControlPlaneAccessIPv4Gateway string                                     `tfschema:"control_plane_access_ipv4_gateway"`
	ControlPlaneAccessIPv4Subnet  string                                     `tfschema:"control_plane_access_ipv4_subnet"`
	ControlPlaneAccessName        string                                     `tfschema:"control_plane_access_name"`
	CoreNetworkTechnology         string                                     `tfschema:"core_network_technology"`
	LocalDiagnosticsAccess        []LocalDiagnosticsAccessConfigurationModel `tfschema:"local_diagnostics_access"`
	Location                      string                                     `tfschema:"location"`
	SiteIds                       []string                                   `tfschema:"site_ids"`
	Platform                      []PlatformConfigurationModel               `tfschema:"platform"`
	Sku                           string                                     `tfschema:"sku"`
	UeMtu                         int64                                      `tfschema:"user_equipment_mtu_in_bytes"`
	InteropSettings               string                                     `tfschema:"interoperability_settings_json"`
	Identity                      []identity.ModelUserAssigned               `tfschema:"identity"`
	Tags                          map[string]string                          `tfschema:"tags"`
	SoftwareVersion               string                                     `tfschema:"software_version"`
}

type LocalDiagnosticsAccessConfigurationModel struct {
	AuthenticationType        string `tfschema:"authentication_type"`
	HttpsServerCertificateUrl string `tfschema:"https_server_certificate_url"`
}

type PlatformConfigurationModel struct {
	AzureStackEdgeDeviceId string `tfschema:"edge_device_id"`
	AzureStackHciClusterId string `tfschema:"stack_hci_cluster_id"`
	ConnectedClusterId     string `tfschema:"arc_kubernetes_cluster_id"`
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

		"control_plane_access_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"control_plane_access_ipv4_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"control_plane_access_ipv4_subnet": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.CIDR,
		},

		"control_plane_access_ipv4_gateway": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

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
						AtLeastOneOf: []string{
							"platform.0.edge_device_id",
							"platform.0.stack_hci_cluster_id",
							"platform.0.arc_kubernetes_cluster_id",
							"platform.0.custom_location_id",
						},
						ValidateFunc: devices.ValidateDataBoxEdgeDeviceID,
					},

					"stack_hci_cluster_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: clusters.ValidateClusterID,
						AtLeastOneOf: []string{
							"platform.0.edge_device_id",
							"platform.0.stack_hci_cluster_id",
							"platform.0.arc_kubernetes_cluster_id",
							"platform.0.custom_location_id",
						},
					},

					"arc_kubernetes_cluster_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateKubernetesClusterID,
						AtLeastOneOf: []string{
							"platform.0.edge_device_id",
							"platform.0.stack_hci_cluster_id",
							"platform.0.arc_kubernetes_cluster_id",
							"platform.0.custom_location_id",
						},
					},

					"custom_location_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID, // TODO: use the resource validate function after custom location onboarded.
						AtLeastOneOf: []string{
							"platform.0.edge_device_id",
							"platform.0.stack_hci_cluster_id",
							"platform.0.arc_kubernetes_cluster_id",
							"platform.0.custom_location_id",
						},
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(packetcorecontrolplane.PlatformTypeAKSNegativeHCI),
							string(packetcorecontrolplane.PlatformTypeThreePNegativeAZURENegativeSTACKNegativeHCI),
							"BaseVM", // tracked on https://github.com/Azure/azure-rest-api-specs/issues/23243
						}, false),
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"interoperability_settings_json": {
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

		"software_version": {
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

			controlPlane := packetcorecontrolplane.PacketCoreControlPlane{
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
				props.CoreNetworkTechnology = pointer.To(packetcorecontrolplane.CoreNetworkType(model.CoreNetworkTechnology))
			}

			props.ControlPlaneAccessInterface = packetcorecontrolplane.InterfaceProperties{}

			if model.ControlPlaneAccessName != "" {
				props.ControlPlaneAccessInterface.Name = &model.ControlPlaneAccessName
			}

			if model.ControlPlaneAccessIPv4Address != "" {
				props.ControlPlaneAccessInterface.IPv4Address = &model.ControlPlaneAccessIPv4Address
			}

			if model.ControlPlaneAccessIPv4Subnet != "" {
				props.ControlPlaneAccessInterface.IPv4Subnet = &model.ControlPlaneAccessIPv4Subnet
			}

			if model.ControlPlaneAccessIPv4Gateway != "" {
				props.ControlPlaneAccessInterface.IPv4Gateway = &model.ControlPlaneAccessIPv4Gateway
			}

			if model.InteropSettings != "" {
				var interopSettingsValue interface{}
				err = json.Unmarshal([]byte(model.InteropSettings), &interopSettingsValue)
				if err != nil {
					return err
				}
				props.InteropSettings = &interopSettingsValue
			}

			props.LocalDiagnosticsAccess = expandPacketCoreControlLocalDiagnosticsAccessConfiguration(model.LocalDiagnosticsAccess)

			props.Platform = expandPlatformConfigurationModel(model.Platform)

			if model.SoftwareVersion != "" {
				props.Version = &model.SoftwareVersion
			}

			controlPlane.Properties = props

			if err := client.CreateOrUpdateThenPoll(ctx, id, controlPlane); err != nil {
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

			if metadata.ResourceData.HasChange("control_plane_access_name") {
				model.Properties.ControlPlaneAccessInterface.Name = &plan.ControlPlaneAccessName
			}

			if metadata.ResourceData.HasChange("control_plane_access_ipv4_address") {
				model.Properties.ControlPlaneAccessInterface.IPv4Address = &plan.ControlPlaneAccessIPv4Address
			}

			if metadata.ResourceData.HasChange("control_plane_access_ipv4_subnet") {
				model.Properties.ControlPlaneAccessInterface.IPv4Subnet = &plan.ControlPlaneAccessIPv4Subnet
			}

			if metadata.ResourceData.HasChange("control_plane_access_ipv4_gateway") {
				model.Properties.ControlPlaneAccessInterface.IPv4Gateway = &plan.ControlPlaneAccessIPv4Gateway
			}

			if metadata.ResourceData.HasChange("core_network_technology") {
				model.Properties.CoreNetworkTechnology = pointer.To(packetcorecontrolplane.CoreNetworkType(plan.CoreNetworkTechnology))
			}

			if metadata.ResourceData.HasChange("interoperability_settings_json") {
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
				model.Properties.Platform = expandPlatformConfigurationModel(plan.Platform)
			}

			if metadata.ResourceData.HasChange("sku") {
				model.Properties.Sku = packetcorecontrolplane.BillingSku(plan.Sku)
			}

			if metadata.ResourceData.HasChange("version") && plan.SoftwareVersion != "" {
				model.Properties.Version = &plan.SoftwareVersion
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

				state.Identity, err = flattenMobileNetworkUserAssignedToNetworkLegacyIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				properties := model.Properties

				state.UeMtu = pointer.From(properties.UeMtu)

				state.ControlPlaneAccessIPv4Address = pointer.From(properties.ControlPlaneAccessInterface.IPv4Address)

				state.ControlPlaneAccessIPv4Gateway = pointer.From(properties.ControlPlaneAccessInterface.IPv4Gateway)

				state.ControlPlaneAccessIPv4Subnet = pointer.From(properties.ControlPlaneAccessInterface.IPv4Subnet)

				state.ControlPlaneAccessName = pointer.From(properties.ControlPlaneAccessInterface.Name)

				// it still needs a nil check because it needs to do type conversion
				if properties.CoreNetworkTechnology != nil {
					state.CoreNetworkTechnology = string(pointer.From(properties.CoreNetworkTechnology))
				}

				// Marshal on a nil interface{} may get random result.
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

				state.SoftwareVersion = pointer.From(properties.Version)
				state.Tags = pointer.From(model.Tags)
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

			// a workaround for that some child resources may still exist for seconds before it fully deleted.
			// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22691
			// it will cause the error "Can not delete resource before nested resources are deleted."
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			stateConf := &pluginsdk.StateChangeConf{
				Delay:   5 * time.Minute,
				Pending: []string{"409"},
				Target:  []string{"200", "202"},
				Refresh: func() (result interface{}, state string, err error) {
					resp, err := client.Delete(ctx, *id)
					if err != nil {
						if resp.HttpResponse.StatusCode == http.StatusConflict {
							return nil, "409", nil
						}
						return nil, "", err
					}
					return resp, "200", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if future, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deleting of %s: %+v", id, err)
			} else {
				poller := future.(packetcorecontrolplane.DeleteOperationResponse).Poller
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
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
	for _, s := range input {
		outputs = append(outputs, s.Id)
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

func expandPlatformConfigurationModel(inputList []PlatformConfigurationModel) packetcorecontrolplane.PlatformConfiguration {
	output := packetcorecontrolplane.PlatformConfiguration{}
	if len(inputList) == 0 {
		return output
	}

	input := inputList[0]

	output.Type = packetcorecontrolplane.PlatformType(input.Type)

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

	return output
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
