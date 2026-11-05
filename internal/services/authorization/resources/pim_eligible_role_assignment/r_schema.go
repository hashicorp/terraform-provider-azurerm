// Copyright IBM Corp. 2014, 2025
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
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Scope for this eligible role assignment, should be a valid resource ID",
			ValidateFunc: validation.Any(
				// Elevated access for a global admin is needed to assign roles in this scope:
				// https://docs.microsoft.com/azure/role-based-access-control/elevate-access-global-admin#azure-cli
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
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Role definition ID for this eligible role assignment",
		},

		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Object ID of the principal for this eligible role assignment",
			ValidateFunc: validation.IsUUID,
		},

		"justification": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			Description: "The justification for this eligible role assignment",
		},

		"condition": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition_version"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition"},
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
			}, false),
		},

		"schedule": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			Description: "The schedule details for this eligible role assignment",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date_time": {
						Optional: true,
						// Note: O+C if not set, defaults to current time
						Computed:    true,
						ForceNew:    true,
						Type:        pluginsdk.TypeString,
						Description: "The start date/time",
					},

					"expiration": { // if none specified, it's a permanent assignment
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"duration_days": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									ForceNew: true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_hours",
										"schedule.0.expiration.0.end_date_time",
									},
									Description: "The duration of the eligible role assignment in days",
								},

								"duration_hours": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									ForceNew: true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_days",
										"schedule.0.expiration.0.end_date_time",
									},
									Description: "The duration of the eligible role assignment in hours",
								},

								"end_date_time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true, // azignore:AZS007 - pre-existing violation
									ForceNew: true,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_days",
										"schedule.0.expiration.0.duration_hours",
									},
									Description: "The end date/time of the eligible role assignment",
								},
							},
						},
					},
				},
			},
		},

		"ticket": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
			Description: "Ticket details relating to the eligible assignment",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"number": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "User-supplied ticket number to be included with the request",
					},

					"system": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "User-supplied ticket system name to be included with the request",
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "Type of principal to which the role will be assigned",
		},
	}
}
