package paloalto

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestack struct{}

var _ sdk.ResourceWithUpdate = LocalRulestack{}

const (
	RulestackSecurityServicesCustom       string = "Custom"
	RulestackSecurityServicesNone         string = "None"
	RulestackSecurityServicesBestPractice string = "BestPractice"
)

type LocalRulestackModel struct {
	Name                           string `tfschema:"name"`
	ResourceGroupName              string `tfschema:"resource_group_name"`
	Location                       string `tfschema:"location"`
	AntiSpywareProfile             string `tfschema:"anti_spyware_profile"`
	AntiVirusProfile               string `tfschema:"anti_virus_profile"`
	DNSSubscription                string `tfschema:"dns_subscription"`
	FileBlockingProfile            string `tfschema:"file_blocking_profile"`
	OutboundTrustedCertificateID   string `tfschema:"outbound_trusted_certificate_name"`
	OutboundUntrustedCertificateID string `tfschema:"outbound_untrusted_certificate_name"`
	URLFilteringProfile            string `tfschema:"url_filtering_profile"`
	VulnerabilityProfile           string `tfschema:"vulnerability_profile"`
	Description                    string `tfschema:"description"`
}

func (r LocalRulestack) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localrulestacks.ValidateLocalRulestackID
}

func (r LocalRulestack) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack"
}

func (r LocalRulestack) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalRulestackName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"vulnerability_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"anti_spyware_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"anti_virus_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"url_filtering_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"file_blocking_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"dns_subscription": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RulestackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RulestackSecurityServicesCustom,
				RulestackSecurityServicesNone,
				RulestackSecurityServicesBestPractice,
			}, false),
		},

		"outbound_trusted_certificate_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LocalRulestackCertificateName,
		},

		"outbound_untrusted_certificate_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.LocalRulestackCertificateName,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r LocalRulestack) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRulestack) ModelObject() interface{} {
	return &LocalRulestackModel{}
}

func (r LocalRulestack) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRulestackModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestacks.NewLocalRulestackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			outboundTrustedCert := model.OutboundTrustedCertificateID
			if outboundTrustedCert != "" {
				certId, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(model.OutboundTrustedCertificateID)
				if err != nil {
					return fmt.Errorf("parsing `outbound_trusted_certificate_name` for %s: %+v", id, err)
				}
				outboundTrustedCert = certId.CertificateName
			}

			outboundUntrustedCert := model.OutboundUntrustedCertificateID
			if outboundUntrustedCert != "" {
				certId, err := certificateobjectlocalrulestack.ParseLocalRulestackCertificateID(model.OutboundUntrustedCertificateID)
				if err != nil {
					return fmt.Errorf("parsing `outbound_untrusted_certificate_name` for %s: %+v", id, err)
				}
				outboundUntrustedCert = certId.CertificateName
			}

			localRulestack := localrulestacks.LocalRulestackResource{
				Location: model.Location,
				Properties: localrulestacks.RulestackProperties{
					DefaultMode: pointer.To(localrulestacks.DefaultModeNONE),
					Description: pointer.To(model.Description),
					Scope:       pointer.To(localrulestacks.ScopeTypeLOCAL),
					SecurityServices: &localrulestacks.SecurityServices{
						AntiSpywareProfile:         pointer.To(model.AntiSpywareProfile),
						AntiVirusProfile:           pointer.To(model.AntiVirusProfile),
						DnsSubscription:            pointer.To(model.DNSSubscription),
						FileBlockingProfile:        pointer.To(model.FileBlockingProfile),
						OutboundTrustCertificate:   pointer.To(outboundTrustedCert),
						OutboundUnTrustCertificate: pointer.To(outboundUntrustedCert),
						UrlFilteringProfile:        pointer.To(model.URLFilteringProfile),
						VulnerabilityProfile:       pointer.To(model.VulnerabilityProfile),
					},
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, localRulestack); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r LocalRulestack) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRulestackModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			props := existing.Model.Properties

			state.Name = id.LocalRulestackName
			state.ResourceGroupName = id.ResourceGroupName
			state.Description = pointer.From(props.Description)
			state.Location = location.Normalize(existing.Model.Location)

			if secServices := props.SecurityServices; secServices != nil {
				if cert := pointer.From(secServices.OutboundUnTrustCertificate); cert != "" {
					state.OutboundUntrustedCertificateID = certificateobjectlocalrulestack.NewLocalRulestackCertificateID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, cert).ID()
				}
				if cert := pointer.From(secServices.OutboundTrustCertificate); cert != "" {
					state.OutboundTrustedCertificateID = certificateobjectlocalrulestack.NewLocalRulestackCertificateID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, cert).ID()
				}
				state.VulnerabilityProfile = pointer.From(secServices.VulnerabilityProfile)
				state.AntiSpywareProfile = pointer.From(secServices.AntiSpywareProfile)
				state.AntiVirusProfile = pointer.From(secServices.AntiVirusProfile)
				state.FileBlockingProfile = pointer.From(secServices.FileBlockingProfile)
				state.URLFilteringProfile = pointer.From(secServices.UrlFilteringProfile)
				state.DNSSubscription = pointer.From(secServices.DnsSubscription)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRulestack) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
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

func (r LocalRulestack) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := LocalRulestackModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			localRulestack := *existing.Model
			props := localRulestack.Properties
			update := localrulestacks.LocalRulestackResourceUpdateProperties{
				DefaultMode: props.DefaultMode,
				Description: props.Description,
				PanLocation: props.PanLocation,
				Scope:       props.Scope,
			}

			if metadata.ResourceData.HasChange("description") {
				update.Description = pointer.To(model.Description)
			}

			secServices := pointer.From(props.SecurityServices)

			if metadata.ResourceData.HasChange("dns_subscription") {
				secServices.DnsSubscription = pointer.To(model.DNSSubscription)
			}

			if metadata.ResourceData.HasChange("vulnerability_profile") {
				secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
			}

			if metadata.ResourceData.HasChange("anti_spyware_profile") {
				secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
			}

			if metadata.ResourceData.HasChange("anti_virus_profile") {
				secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
			}

			if metadata.ResourceData.HasChange("url_filtering_profile") {
				secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
			}

			if metadata.ResourceData.HasChange("file_blocking_profile") {
				secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
			}

			if metadata.ResourceData.HasChange("outbound_trusted_certificate_name") {
				secServices.OutboundTrustCertificate = pointer.To(model.OutboundTrustedCertificateID)
			}

			if metadata.ResourceData.HasChange("outbound_untrusted_certificate_name") {
				secServices.OutboundUnTrustCertificate = pointer.To(model.OutboundUntrustedCertificateID)
			}

			update.SecurityServices = pointer.To(secServices)

			if _, err = client.Update(ctx, *id, localrulestacks.LocalRulestackResourceUpdate{Properties: &update}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if _, err = client.Commit(ctx, *id); err != nil {
				return fmt.Errorf("committing config for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
