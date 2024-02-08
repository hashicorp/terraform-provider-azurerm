// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	FQDNLists   []string `tfschema:"local_rulestack_fqdn_list_ids"`
	PrefixLists []string `tfschema:"local_rulestack_prefix_list_ids"`
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
						"destination.0.local_rulestack_fqdn_list_ids",
						"destination.0.local_rulestack_prefix_list_ids",
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
						"destination.0.local_rulestack_fqdn_list_ids",
						"destination.0.local_rulestack_prefix_list_ids",
					},
				},

				"feeds": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					AtLeastOneOf: []string{
						"destination.0.cidrs",
						"destination.0.countries",
						"destination.0.feeds",
						"destination.0.local_rulestack_fqdn_list_ids",
						"destination.0.local_rulestack_prefix_list_ids",
					},
				},

				"local_rulestack_fqdn_list_ids": {
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
						"destination.0.local_rulestack_fqdn_list_ids",
						"destination.0.local_rulestack_prefix_list_ids",
					},
				},

				"local_rulestack_prefix_list_ids": {
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
						"destination.0.local_rulestack_fqdn_list_ids",
						"destination.0.local_rulestack_prefix_list_ids",
					},
				},
			},
		},
	}
}

func ExpandDestination(input []Destination) (*localrules.DestinationAddr, error) {
	if len(input) == 0 {
		return nil, nil
	}

	d := input[0]
	prefixLists := make([]string, 0)
	if len(d.PrefixLists) > 0 {
		for _, p := range d.PrefixLists {
			id, err := prefixlistlocalrulestack.ParseLocalRulestackPrefixListID(p)
			if err != nil {
				return nil, err
			}
			prefixLists = append(prefixLists, id.PrefixListName)
		}
	}

	fqdnLists := make([]string, 0)
	if len(d.FQDNLists) > 0 {
		for _, p := range d.FQDNLists {
			id, err := fqdnlistlocalrulestack.ParseLocalRulestackFqdnListID(p)
			if err != nil {
				return nil, err
			}
			fqdnLists = append(fqdnLists, id.FqdnListName)
		}
	}

	return &localrules.DestinationAddr{
		Cidrs:       pointer.To(d.CIDRS),
		Countries:   pointer.To(d.Countries),
		Feeds:       pointer.To(d.Feeds),
		FqdnLists:   pointer.To(fqdnLists),
		PrefixLists: pointer.To(prefixLists),
	}, nil
}

func FlattenDestination(input *localrules.DestinationAddr, ruleId localrules.LocalRuleId) []Destination {
	if input == nil {
		return []Destination{}
	}

	prefixLists := make([]string, 0)
	if p := input.PrefixLists; p != nil {
		for _, v := range *p {
			prefixLists = append(prefixLists, prefixlistlocalrulestack.NewLocalRulestackPrefixListID(ruleId.SubscriptionId, ruleId.ResourceGroupName, ruleId.LocalRulestackName, v).ID())
		}
	}

	fqdnLists := make([]string, 0)
	if p := input.FqdnLists; p != nil {
		for _, v := range *p {
			fqdnLists = append(fqdnLists, fqdnlistlocalrulestack.NewLocalRulestackFqdnListID(ruleId.SubscriptionId, ruleId.ResourceGroupName, ruleId.LocalRulestackName, v).ID())
		}
	}

	return []Destination{{
		CIDRS:       pointer.From(input.Cidrs),
		Countries:   pointer.From(input.Countries),
		Feeds:       pointer.From(input.Feeds),
		FQDNLists:   fqdnLists,
		PrefixLists: prefixLists,
	}}
}
