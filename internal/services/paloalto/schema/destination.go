package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"reflect"
)

type Destination struct {
	CIDRS       []string `tfschema:"cidrs"`
	Countries   []string `tfschema:"countries"`
	Feeds       []string `tfschema:"feeds"`
	FQDNLists   []string `tfschema:"fqdn_lists"`
	PrefixLists []string `tfschema:"prefix_lists"`
}

var DestinationDefault = localrules.DestinationAddr{
	Cidrs: pointer.To([]string{"any"}),
}

func DesintationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
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
				},

				"countries": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.ISO3361CountryCode,
					},
				},

				"feeds": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: nil, // TODO - This is another resource type?
					},
				},

				"fqdn_lists": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: nil, // TODO - FQDN List is another resource type belonging to a RuleStack, use the name validation when it exists
					},
				},

				"prefix_lists": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: nil, // TODO - Prefix List is another resource type belonging to a RuleStack, use the name validation when it exists
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
	if input == nil || reflect.DeepEqual(*input, DestinationDefault) {
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
