package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackDataSource struct{}

var _ sdk.DataSource = LocalRulestackDataSource{}

type LocalRulestackDataSourceModel struct {
	Name                       string `tfschema:"name"`
	ResourceGroupName          string `tfschema:"resource_group_name"`
	Location                   string `tfschema:"location"`
	AntiSpywareProfile         string `tfschema:"anti_spyware_profile"`
	AntiVirusProfile           string `tfschema:"anti_virus_profile"`
	DNSSubscription            string `tfschema:"dns_subscription"`
	FileBlockingProfile        string `tfschema:"file_blocking_profile"`
	URLFilteringProfile        string `tfschema:"url_filtering_profile"`
	VulnerabilityProfile       string `tfschema:"vulnerability_profile"`
	OutboundTrustCertificate   string `tfschema:"outbound_trust_certificate"`
	OutboundUnTrustCertificate string `tfschema:"outbound_untrust_certificate"`
	Description                string `tfschema:"description"`
}

func (l LocalRulestackDataSource) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack"
}

func (l LocalRulestackDataSource) ModelObject() interface{} {
	return &LocalRulestackDataSourceModel{}
}

func (l LocalRulestackDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.LocalRuleStackName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"location": commonschema.LocationWithoutForceNew(),
	}
}

func (l LocalRulestackDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"vulnerability_profile": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"anti_spyware_profile": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"anti_virus_profile": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"url_filtering_profile": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"file_blocking_profile": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_subscription": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_trust_certificate": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_untrust_certificate": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (l LocalRulestackDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			var model LocalRulestackDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestacks.NewLocalRulestackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			props := existing.Model.Properties

			model.Description = pointer.From(props.Description)
			if secServices := props.SecurityServices; secServices != nil {
				model.FileBlockingProfile = pointer.From(secServices.FileBlockingProfile)
				model.AntiVirusProfile = pointer.From(secServices.AntiVirusProfile)
				model.AntiSpywareProfile = pointer.From(secServices.AntiSpywareProfile)
				model.URLFilteringProfile = pointer.From(secServices.UrlFilteringProfile)
				model.VulnerabilityProfile = pointer.From(secServices.VulnerabilityProfile)
				model.DNSSubscription = pointer.From(secServices.DnsSubscription)
				model.OutboundTrustCertificate = pointer.From(secServices.OutboundTrustCertificate)
				model.OutboundUnTrustCertificate = pointer.From(secServices.OutboundUnTrustCertificate)
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}
