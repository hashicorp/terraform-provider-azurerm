package paloalto

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	helpersValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVHubPanoramaResource struct{}

type NextGenerationFirewallVHubPanoramaModel struct {
	Name                 string                      `tfschema:"name"`
	ResourceGroupName    string                      `tfschema:"resource_group_name"`
	PanoramaBase64Config string                      `tfschema:"panorama_base64_config"`
	Location             string                      `tfschema:"location"`
	NetworkProfile       []schema.NetworkProfileVHub `tfschema:"network_profile"`
	DNSSettings          []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd             []schema.DestinationNAT     `tfschema:"destination_nat"`
	PanoramaConfig       []schema.Panorama           `tfschema:"panorama_config"`

	// Computed
	PlanData []schema.Plan `tfschema:"plan"`
	PanEtag  string        `tfschema:"pan_etag"`

	Tags map[string]interface{} `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVHubPanoramaResource{}

func (r NextGenerationFirewallVHubPanoramaResource) ModelObject() interface{} {
	return &NextGenerationFirewallVHubPanoramaModel{}
}

func (r NextGenerationFirewallVHubPanoramaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVHubPanoramaResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_vhub_panorama"
}

func (r NextGenerationFirewallVHubPanoramaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"panorama_base64_config": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: helpersValidate.Base64EncodedString,
		},

		"network_profile": schema.VHubNetworkProfileSchema(),

		// Optional
		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"panorama_config": schema.PanoramaSchema(),

		"plan": schema.PlanSchema(),

		"pan_etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

			var model NextGenerationFirewallVHubPanoramaModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := firewalls.NewFirewallID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			firewall := firewalls.FirewallResource{
				Location: location.Normalize(model.Location),
				Properties: firewalls.FirewallDeploymentProperties{
					PanoramaConfig: &firewalls.PanoramaConfig{
						ConfigString: model.PanoramaBase64Config,
					},
					IsPanoramaManaged: pointer.To(firewalls.BooleanEnumTRUE),
					DnsSettings:       schema.ExpandDNSSettings(model.DNSSettings),
					MarketplaceDetails: firewalls.MarketplaceDetails{
						OfferId:     "pan_swfw_cloud_ngfw", // TODO - Will just supplying the offer ID `panw-cloud-ngfw-payg` work?
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVHub(model.NetworkProfile),
					PlanData: firewalls.PlanData{
						BillingCycle: firewalls.BillingCycleMONTHLY,
						PlanId:       "panw-cloud-ngfw-payg",
					},
					FrontEndSettings: schema.ExpandDestinationNAT(model.FrontEnd),
				},
				Tags: tags.Expand(model.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVHubPanoramaModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.FirewallName
			state.ResourceGroupName = id.ResourceGroupName

			props := existing.Model.Properties
			state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

			netProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
			if err != nil {
				return fmt.Errorf("parsing Network Profile for %s: %+v", *id, err)
			}

			state.NetworkProfile = []schema.NetworkProfileVHub{*netProfile}

			if feSettings := pointer.From(props.FrontEndSettings); len(feSettings) != 0 {
				fes := make([]schema.DestinationNAT, 0)
				for _, v := range feSettings {
					bePort, _ := strconv.Atoi(v.BackendConfiguration.Port)
					fePort, _ := strconv.Atoi(v.FrontendConfiguration.Port)
					fe := schema.DestinationNAT{
						Name:     v.Name,
						Protocol: string(v.Protocol),
						BackendConfiguration: []schema.BackendEndpointConfiguration{{
							PublicIP: pointer.From(v.BackendConfiguration.Address.Address),
							Port:     bePort,
						}},
						FrontendConfiguration: []schema.FrontendEndpointConfiguration{{
							PublicIPID: pointer.From(v.FrontendConfiguration.Address.ResourceId),
							Port:       fePort,
						}},
					}

					fes = append(fes, fe)
				}
				state.FrontEnd = fes
			}

			if panoramaConfig := props.PanoramaConfig; panoramaConfig != nil {
				state.PanoramaBase64Config = panoramaConfig.ConfigString
				state.PanoramaConfig = []schema.Panorama{{
					Name:            pointer.From(panoramaConfig.CgName),
					DeviceGroupName: pointer.From(panoramaConfig.DgName),
					HostName:        pointer.From(panoramaConfig.HostName),
					PanoramaServer:  pointer.From(panoramaConfig.PanoramaServer),
					PanoramaServer2: pointer.From(panoramaConfig.PanoramaServer2),
					TplName:         pointer.From(panoramaConfig.TplName),
					VMAuthKey:       pointer.From(panoramaConfig.VMAuthKey),
				}}
			}

			state.PanEtag = pointer.From(props.PanEtag)

			state.PlanData = schema.FlattenPlanData(props.PlanData)

			state.Tags = tags.Flatten(existing.Model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
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

func (r NextGenerationFirewallVHubPanoramaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			return nil
		},
	}
}
