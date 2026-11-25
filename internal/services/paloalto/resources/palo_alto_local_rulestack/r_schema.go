// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalRuleStackName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"vulnerability_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_spyware_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_virus_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"url_filtering_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"file_blocking_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"dns_subscription": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
