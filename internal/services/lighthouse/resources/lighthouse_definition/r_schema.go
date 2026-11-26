// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lighthouseDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Lighthouse Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"managing_tenant_id": {
			Description:  "The ID of the managing tenant.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"scope": {
			Description:  "The ID of the managed subscription.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"authorization": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Description:  "Principal ID of the security group/Service Principal/user that would be assigned permissions to the projected subscription.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"role_definition_id": {
						Description:  "The Role Definition identifier. This role will define the permissions that are granted to the principal. This cannot be an `Owner` role.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"principal_display_name": {
						Description:  "The display name of the security group/Service Principal/user that would be assigned permissions to the projected subscription.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"delegated_role_definition_ids": {
						Description: "The set of Role Definition IDs which define all the permissions that the principal ID can assign.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},

		"description": {
			Description: "A description of the Lighthouse Definition.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"eligible_authorization": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"principal_id": {
						Description:  "Principal ID of the security group/Service Principal/user that would be assigned permissions to the projected subscription.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"role_definition_id": {
						Description:  "The Principal ID of the Azure built-in role that defines the permissions that the Azure Active Directory will have on the projected scope.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"just_in_time_access_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"multi_factor_auth_provider": {
									Description: "The multi-factor authorization provider to be used for just-in-time access requests.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(registrationdefinitions.MultiFactorAuthProviderAzure),
									}, false),
								},

								"maximum_activation_duration": {
									Description:  "The maximum access duration in ISO 8601 format for just-in-time access requests.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "PT8H",
									ValidateFunc: validation.ISO8601Duration,
								},

								"approver": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"principal_id": {
												Description:  "The Principal ID of the Azure Active Directory principal for the approver.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.IsUUID,
											},

											"principal_display_name": {
												Description:  "The display name of the Azure Active Directory Principal for the approver.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"principal_display_name": {
						Description:  "The display name of the Azure Active Directory Principal.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"lighthouse_definition_id": {
			Description:  "A unique UUID/GUID which identifies this lighthouse definition - one will be generated if not specified.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The plan name of the marketplace offer.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Description:  "The publisher ID of the plan.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Description:  "The product code of the plan.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Description:  "The version of the plan.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
