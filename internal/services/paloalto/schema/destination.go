package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Destination struct {
	CIDRS       []string `tfschema:"cidrs"`
	Countries   []string `tfschema:"countries"`
	Feeds       []string `tfschema:"feeds"`
	FQDNLists   []string `tfschema:"fqdn_list_ids"`
	PrefixLists []string `tfschema:"prefix_list_ids"`
}

func DestinationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cidrs": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.Any(
							validation.IsCIDR,
							validation.StringInSlice([]string{"any"}, false),
						),
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.fqdn_list_ids",
						"destination.0.prefix_list_ids",
					},
				},

				"countries": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.ISO3361CountryCode,
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.fqdn_list_ids",
						"destination.0.prefix_list_ids",
					},
				},

				"feeds": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: nil, // TODO - This is another resource type?
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.fqdn_list_ids",
						"destination.0.prefix_list_ids",
					},
				},

				"fqdn_list_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: fqdnlistlocalrulestack.ValidateLocalRulestackFqdnListID,
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.fqdn_list_ids",
						"destination.0.prefix_list_ids",
					},
				},

				"prefix_list_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: prefixlistlocalrulestack.ValidateLocalRulestackPrefixListID,
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.fqdn_list_ids",
						"destination.0.prefix_list_ids",
					},
				},
			},
		},
	}
}

func ExpandDestination(input []Destination) *localrules.DestinationAddr {
	if len(input) == 0 {
		return nil
	}

	d := input[0]

	return &localrules.DestinationAddr{
		Cidrs:       pointer.To(d.CIDRS),
		Countries:   pointer.To(d.Countries),
		Feeds:       pointer.To(d.Feeds),
		FqdnLists:   pointer.To(d.FQDNLists),
		PrefixLists: pointer.To(d.PrefixLists),
	}
}

func FlattenDestination(input *localrules.DestinationAddr) []Destination {
	if input == nil {
		return []Destination{}
	}

	return []Destination{{
		CIDRS:       pointer.From(input.Cidrs),
		Countries:   pointer.From(input.Countries),
		Feeds:       pointer.From(input.Feeds),
		FQDNLists:   pointer.From(input.FqdnLists),
		PrefixLists: pointer.From(input.PrefixLists),
	}}
}
