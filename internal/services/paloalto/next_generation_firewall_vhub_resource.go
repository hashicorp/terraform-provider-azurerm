package paloalto

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVHubResource struct{}

type NextGenerationFirewallVHubModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	// TODO - Bring VHub ID up to top level?
	Location       string                      `tfschema:"location"` // TODO RG Location only, or other OK?
	NetworkProfile []schema.NetworkProfileVHub `tfschema:"network_profile"`
	RuleStackId    string                      `tfschema:"rule_stack_id"`
	DNSSettings    []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd       []schema.FrontEnd           `tfschema:"front_end"`
	PanoramaConfig []schema.Panorama           `tfschema:"panorama"`

	// Computed
	PlanData []schema.Plan `tfschema:"plan"`
	PanEtag  string        `tfschema:"pan_etag"`

	Tags map[string]interface{} `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVHubResource{}

func (r NextGenerationFirewallVHubResource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetModel{}
}

func (r NextGenerationFirewallVHubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVHubResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateFunc:     location.EnhancedValidate,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			ConflictsWith:    []string{"rule_stack_id"},
		},

		"network_profile": schema.VHubNetworkProfileSchema(),

		// Optional
		"dns_settings": schema.DNSSettingsSchema(),

		"rule_stack_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: localrulestacks.ValidateLocalRulestackID,
			ExactlyOneOf: []string{
				"panorama",
				"rule_stack_id",
			},
		},

		"panorama": schema.PanoramaSchema(),

		"front_end": schema.FrontEndSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r NextGenerationFirewallVHubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"plan": schema.PlanSchema(),

		"pan_etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NextGenerationFirewallVHubResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_vhub"
}

func (r NextGenerationFirewallVHubResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient
			localRulestackClient := metadata.Client.PaloAlto.LocalRulestacksClient
			loc := ""
			var model NextGenerationFirewallVHubModel

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
				Properties: firewalls.FirewallDeploymentProperties{
					DnsSettings: firewalls.DNSSettings{
						//EnableDnsProxy: pointer.To(firewalls.DNSProxyDISABLED),
						//EnabledDnsType: pointer.To(firewalls.EnabledDNSTypeCUSTOM),
					},
					MarketplaceDetails: firewalls.MarketplaceDetails{
						OfferId:     "pan_swfw_cloud_ngfw", // TODO - Will just supplying the offer ID `panw-cloud-ngfw-payg` work?
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVHub(model.NetworkProfile),
					PlanData: firewalls.PlanData{
						BillingCycle:  firewalls.BillingCycleMONTHLY,
						PlanId:        "panw-cloud-ngfw-payg",
						EffectiveDate: pointer.To("0001-01-01T00:00:00Z"),
					},
				},
				Tags: tags.Expand(model.Tags),
			}

			if len(model.DNSSettings) > 0 {
				dns := model.DNSSettings[0]
				dnsSettings := firewalls.DNSSettings{
					EnableDnsProxy: pointer.To(firewalls.DNSProxyENABLED),
				}

				if len(dns.DnsServers) > 0 {
					dnsSettings.EnabledDnsType = pointer.To(firewalls.EnabledDNSTypeCUSTOM)
					dnsServers := make([]firewalls.IPAddress, 0)
					for _, v := range dns.DnsServers {
						dnsServers = append(dnsServers, firewalls.IPAddress{
							Address: pointer.To(v),
						})
					}
					dnsSettings.DnsServers = pointer.To(dnsServers)
				}

				if dns.AzureDNS {
					dnsSettings.EnabledDnsType = pointer.To(firewalls.EnabledDNSTypeAZURE)
				}

				firewall.Properties.DnsSettings = dnsSettings
			}

			if model.RuleStackId != "" {
				ruleStackID, err := localrulestacks.ParseLocalRulestackID(model.RuleStackId)
				if err != nil {
					return err
				}

				ruleStack, err := localRulestackClient.Get(ctx, *ruleStackID)
				if err != nil {
					return fmt.Errorf("reading %s for %s: %+v", ruleStackID, id, err)
				}

				firewall.Location = location.Normalize(ruleStack.Model.Location)
				loc = location.Normalize(ruleStack.Model.Location)
				firewall.Properties.AssociatedRulestack = &firewalls.RulestackDetails{
					ResourceId: pointer.To(ruleStackID.ID()),
					Location:   pointer.To(loc),
				}
			}

			if len(model.PanoramaConfig) > 0 {
				firewall.Location = location.Normalize(model.Location)
				firewall.Properties.IsPanoramaManaged = pointer.To(firewalls.BooleanEnumTRUE)
				firewall.Properties.PanoramaConfig = &firewalls.PanoramaConfig{ConfigString: model.PanoramaConfig[0].B64Config}
			}

			if len(model.FrontEnd) > 0 {
				fes := make([]firewalls.FrontendSetting, 0)
				for _, v := range model.FrontEnd {
					fe := firewalls.FrontendSetting{
						Name:                  v.Name,
						Protocol:              firewalls.ProtocolType(v.Protocol),
						BackendConfiguration:  firewalls.EndpointConfiguration{},
						FrontendConfiguration: firewalls.EndpointConfiguration{},
					}

					if len(v.FrontendConfiguration) > 0 {
						fec := v.FrontendConfiguration[0]
						fe.BackendConfiguration = firewalls.EndpointConfiguration{
							Address: firewalls.IPAddress{
								ResourceId: pointer.To(fec.PublicIPID),
							},
							Port: strconv.Itoa(fec.Port),
						}
					}

					if len(v.BackendConfiguration) > 0 {
						bec := v.BackendConfiguration[0]
						fe.BackendConfiguration = firewalls.EndpointConfiguration{
							Address: firewalls.IPAddress{
								ResourceId: pointer.To(bec.PublicIPID),
							},
							Port: strconv.Itoa(bec.Port),
						}
					}

					fes = append(fes, fe)
				}

				firewall.Properties.FrontEndSettings = pointer.To(fes)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVHubResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVHubModel

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
			dns := props.DnsSettings

			if dnsServers := pointer.From(dns.DnsServers); len(dnsServers) > 0 {
				dnsSettings := make([]string, 0)
				for _, v := range dnsServers {
					dnsSettings = append(dnsSettings, pointer.From(v.Address))
				}
			}

			netProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
			if err != nil {
				return fmt.Errorf("parsing Network Profile for %s: %+v", *id, err)
			}

			state.NetworkProfile = []schema.NetworkProfileVHub{*netProfile}

			if feSettings := pointer.From(props.FrontEndSettings); len(feSettings) != 0 {
				fes := make([]schema.FrontEnd, 0)
				for _, v := range feSettings {
					bePort, _ := strconv.Atoi(v.BackendConfiguration.Port)
					fePort, _ := strconv.Atoi(v.FrontendConfiguration.Port)
					fe := schema.FrontEnd{
						Name:     v.Name,
						Protocol: string(v.Protocol),
						BackendConfiguration: []schema.EndpointConfiguration{{
							PublicIPID: pointer.From(v.BackendConfiguration.Address.ResourceId),
							Port:       bePort,
						}},
						FrontendConfiguration: []schema.EndpointConfiguration{{
							PublicIPID: pointer.From(v.FrontendConfiguration.Address.ResourceId),
							Port:       fePort,
						}},
					}

					fes = append(fes, fe)
				}
				state.FrontEnd = fes
			}

			if panoramaConfig := props.PanoramaConfig; panoramaConfig != nil {
				state.PanoramaConfig = []schema.Panorama{{
					B64Config:       panoramaConfig.ConfigString,
					Name:            pointer.From(panoramaConfig.CgName),
					DeviceGroupName: pointer.From(panoramaConfig.DgName),
					HostName:        pointer.From(panoramaConfig.HostName),
					PanoramaServer:  pointer.From(panoramaConfig.PanoramaServer),
					PanoramaServer2: pointer.From(panoramaConfig.PanoramaServer2),
					TplName:         pointer.From(panoramaConfig.TplName),
					VMAuthKey:       pointer.From(panoramaConfig.VMAuthKey),
				}}
			}

			state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)

			state.PanEtag = pointer.From(props.PanEtag)

			state.PlanData = schema.FlattenPlanData(props.PlanData)

			state.Tags = tags.Flatten(existing.Model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVHubResource) Delete() sdk.ResourceFunc {
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

func (r NextGenerationFirewallVHubResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			return nil
		},
	}
}
