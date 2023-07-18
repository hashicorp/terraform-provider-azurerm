package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type Source struct {
	CIDRS       []string `tfschema:"cidrs"`
	Countries   []string `tfschema:"countries"`
	Feeds       []string `tfschema:"feeds"`
	PrefixLists []string `tfschema:"prefix_list_ids"`
}

func SourceSchema() *pluginsdk.Schema {
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
						"source.0.cidrs",
						"source.0.countries",
						"source.0.feeds",
						"source.0.prefix_list_ids",
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
						"source.0.cidrs",
						"source.0.countries",
						"source.0.feeds",
						"source.0.prefix_list_ids",
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
						"source.0.cidrs",
						"source.0.countries",
						"source.0.feeds",
						"source.0.prefix_list_ids",
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
						"source.0.cidrs",
						"source.0.countries",
						"source.0.feeds",
						"source.0.prefix_list_ids",
					},
				},
			},
		},
	}
}

func ExpandSource(input []Source) *localrules.SourceAddr {
	if len(input) == 0 {
		return nil
	}

	d := input[0]

	return &localrules.SourceAddr{
		Cidrs:       pointer.To(d.CIDRS),
		Countries:   pointer.To(d.Countries),
		Feeds:       pointer.To(d.Feeds),
		PrefixLists: pointer.To(d.PrefixLists),
	}
}

func FlattenSource(input *localrules.SourceAddr) []Source {
	if input == nil {
		return []Source{}
	}

	return []Source{{
		CIDRS:       pointer.From(input.Cidrs),
		Countries:   pointer.From(input.Countries),
		Feeds:       pointer.From(input.Feeds),
		PrefixLists: pointer.From(input.PrefixLists),
	}}
}
