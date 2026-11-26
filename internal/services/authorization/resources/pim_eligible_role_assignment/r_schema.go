// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package pim_eligible_role_assignment

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Description: "The scope for this eligible Role Assignment, should be a valid resource ID.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				// Elevated access for a global admin is needed to assign roles in this scope:
				// https://learn.microsoft.com/azure/role-based-access-control/elevate-access-global-admin#azure-cli
				// It seems only user account is allowed to be elevated access.
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),

				validate.EnrollmentID,
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},

		"role_definition_id": {
			Description: "The Role Definition ID for this eligible Role Assignment.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"principal_id": {
			Description:  "Object ID of the principal for this eligible Role Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"justification": {
			Description: "The justification of the Role Assignment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"condition": {
			Description:  "The condition that limits the resources that the role can be assigned to. See the [official conditions documentation](https://learn.microsoft.com/azure/role-based-access-control/conditions-overview#what-are-role-assignment-conditions) for details.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"condition_version"},
		},

		"condition_version": {
			Description: "The version of the condition.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
			}, false),
			RequiredWith: []string{"condition"},
		},

		"schedule": {
			Description: "The schedule details for this eligible Role Assignment.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date_time": {
						Description: "The start date/time of the Role Assignment.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Note: O+C if not set, defaults to current time
						Computed: true,
						ForceNew: true,
					},

					"expiration": { // if none specified, it's a permanent assignment
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"duration_days": {
									Description: "The duration of the Role Assignment in days.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true, // azignore:AZS007 - pre-existing violation
									ForceNew:    true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_hours",
										"schedule.0.expiration.0.end_date_time",
									},
								},

								"duration_hours": {
									Description: "The duration of the Role Assignment in hours.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true, // azignore:AZS007 - pre-existing violation
									ForceNew:    true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_days",
										"schedule.0.expiration.0.end_date_time",
									},
								},

								"end_date_time": {
									Description: "The end date/time of the Role Assignment.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									Computed:    true, // azignore:AZS007 - pre-existing violation
									ForceNew:    true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_days",
										"schedule.0.expiration.0.duration_hours",
									},
								},
							},
						},
					},
				},
			},
		},

		"ticket": {
			Description: "Ticket details relating to the eligible assignment.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"number": {
						Description: "User-supplied ticket number to be included with the request.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},

					"system": {
						Description: "User-supplied ticket system name to be included with the request.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Description: "Type of principal to which the role will be assigned.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
