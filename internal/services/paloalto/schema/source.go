package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"reflect"
)

type Source struct {
	CIDRS       []string `tfschema:"cidrs"`
	Countries   []string `tfschema:"countries"`
	Feeds       []string `tfschema:"feeds"`
	PrefixLists []string `tfschema:"prefix_lists"`
}

var SourceDefault = localrules.SourceAddr{
	Cidrs: pointer.To([]string{"any"}),
}

func SourceSchema() *pluginsdk.Schema {
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

				"prefix_lists": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: nil, // TODO - Prefix List is another resource type belonging to a RuleStack
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
	if input == nil || reflect.DeepEqual(*input, SourceDefault) {
		return []Source{}
	}

	return []Source{{
		CIDRS:       pointer.From(input.Cidrs),
		Countries:   pointer.From(input.Countries),
		Feeds:       pointer.From(input.Feeds),
		PrefixLists: pointer.From(input.PrefixLists),
	}}
}
