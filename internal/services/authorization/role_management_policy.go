package authorization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = RoleManagementPolicyResource{}

var _ sdk.ResourceWithUpdate = RoleManagementPolicyResource{}

var _ sdk.ResourceWithCustomizeDiff = RoleManagementPolicyResource{}

type RoleManagementPolicyResource struct{}

func (r RoleManagementPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The scope.",
		},
		"role_definition_id": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "The role definition id.",
		},
		"activation": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The activation settings for PIM on the scope and role definition.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"maximum_duration_hours": {
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: "The maximum duration in hours for an activation.",
					},
					"require_multi_factor_authentication": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Is Multi Factor Authentication required for an activation?",
					},
					"require_justification": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Is Justification required for an activation?",
					},
					"require_ticket_information": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true,
						Description: "Is Ticket Information required for an activation?",
					},
					"approvers": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Description: "A list of approvers for activation.",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"group": {
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Description: "An approval group",
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"id": {
												Type:        pluginsdk.TypeString,
												Required:    true,
												Description: "The object id of the AAD group.",
											},
											"name": {
												Type:        pluginsdk.TypeString,
												Required:    true,
												Description: "The name of the AAD group.",
											},
										},
									},
								},
								"user": {
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Description: "An approval user",
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"id": {
												Type:        pluginsdk.TypeString,
												Required:    true,
												Description: "The object id of a user.",
											},
											"name": {
												Type:        pluginsdk.TypeString,
												Required:    true,
												Description: "The name of a user.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"assignment": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The assignment settings for PIM on the scope and role definition.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"eligible": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Description: "The eligible settings for an assignment.",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allow_permanent": {
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Computed:    true,
									Description: "Allow permanent eligible assignment.",
									ConflictsWith: []string{
										"assignment.0.eligible.0.expire_after_hours",
										"assignment.0.eligible.0.expire_after_days",
									},
								},
								"expire_after_hours": {
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true,
									Description: "The number of hours after an eligible assignments is expired.",
									ConflictsWith: []string{
										"assignment.0.eligible.0.allow_permanent",
										"assignment.0.eligible.0.expire_after_days",
									},
								},
								"expire_after_days": {
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true,
									Description: "The number of days after an eligible assignments is expired.",
									ConflictsWith: []string{
										"assignment.0.eligible.0.allow_permanent",
										"assignment.0.eligible.0.expire_after_hours",
									},
								},
							},
						},
					},
					"active": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Description: "The active settings for an assignment.",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allow_permanent": {
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Computed:    true,
									Description: "Allow permanent active assignment.",
									ConflictsWith: []string{
										"assignment.0.active.0.expire_after_days",
										"assignment.0.active.0.expire_after_hours",
									},
								},
								"expire_after_hours": {
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true,
									Description: "The number of hours after an active assignments is expired.",
									ConflictsWith: []string{
										"assignment.0.active.0.allow_permanent",
										"assignment.0.active.0.expire_after_days",
									},
								},
								"expire_after_days": {
									Type:        pluginsdk.TypeInt,
									Optional:    true,
									Computed:    true,
									Description: "The number of days after an active assignments is expired.",
									ConflictsWith: []string{
										"assignment.0.active.0.allow_permanent",
										"assignment.0.active.0.expire_after_hours",
									},
								},
								"require_multi_factor_authentication": {
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Computed:    true,
									Description: "Is Multi Factor Authentication required for an active assignment?",
								},
								"require_justification": {
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Computed:    true,
									Description: "Is Justification required for an active assignment?",
								},
							},
						},
					},
				},
			},
		},
		"notifications": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "The notification settings for PIM on the scope and role definition.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"member_assigned_eligible": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Computed:    true,
						Description: "Notifications settings when members are assigned as eligible to this role.",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"role_assignment_alert": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Role assignment alerts.",
									Elem:        notificationConfiguration(),
								},
								"assigned_user": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the assigned user (assignee).",
									Elem:        notificationConfiguration(),
								},
								"request_for_extension_or_approval": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Request to approve a role assignment renewal/extension.",
									Elem:        notificationConfiguration(),
								},
							},
						},
					},
					"member_assigned_active": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Computed:    true,
						Description: "Notifications settings when members are assigned as active to this role.						",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"role_assignment_alert": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Role assignment alerts.",
									Elem:        notificationConfiguration(),
								},
								"assigned_user": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the assigned user (assignee).",
									Elem:        notificationConfiguration(),
								},
								"request_for_extension_or_approval": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Request to approve a role assignment renewal/extension.",
									Elem:        notificationConfiguration(),
								},
							},
						},
					},
					"eligible_member_activate": {
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						Optional:    true,
						Computed:    true,
						Description: "Notifications settings when eligible members activate this role.",
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"role_assignment_alert": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Role assignment alerts.",
									Elem:        notificationConfiguration(),
								},
								"assigned_user": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the assigned user (assignee).",
									Elem:        notificationConfiguration(),
								},
								"request_for_extension_or_approval": {
									Type:        pluginsdk.TypeList,
									MaxItems:    1,
									Optional:    true,
									Computed:    true,
									Description: "Notification settings for the Request to approve a role assignment renewal/extension.",
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"default_recipients": {
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												Computed:    true,
												Description: "Will notifications be sent to the default recipients?",
											},
											"critical_emails_only": {
												Type:        pluginsdk.TypeBool,
												Optional:    true,
												Computed:    true,
												Description: "Will critical emails only be sent?",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r RoleManagementPolicyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			// custom validation for multiple properties

			if rd.HasChange("notifications.0.eligible_member_activate.0.request_for_extension_or_approval") {
				_, criticalEmailsAfter := rd.GetChange("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.critical_emails_only")
				_, defaultRecipientsAfter := rd.GetChange("notifications.0.eligible_member_activate.0.request_for_extension_or_approval.0.default_recipients")

				if criticalEmailsAfter.(bool) && !defaultRecipientsAfter.(bool) {
					return fmt.Errorf("cannot enable critical emails and disable default recipients for `request_for_extension_or_approval` on `eligible_member_activate`")
				}
			}

			if rd.HasChange("assignment.0.eligible") {
				_, assignmentEligibleAllowPermanent := rd.GetChange("assignment.0.eligible.0.allow_permanent")
				_, assignmentEligibleExpireAfterHours := rd.GetChange("assignment.0.eligible.0.expire_after_hours")
				_, assignmentEligibleExpireAfterDays := rd.GetChange("assignment.0.eligible.0.expire_after_days")

				if !assignmentEligibleAllowPermanent.(bool) && assignmentEligibleExpireAfterHours.(int) == 0 && assignmentEligibleExpireAfterDays.(int) == 0 {
					return fmt.Errorf("cannot set allow permanent to false when expire after days and hours is 0 on `assignment` for `eligible`")
				}
			}

			if rd.HasChange("assignment.0.active") {
				_, assignmentActiveAllowPermanent := rd.GetChange("assignment.0.active.0.allow_permanent")
				_, assignmentActiveExpireAfterHours := rd.GetChange("assignment.0.active.0.expire_after_hours")
				_, assignmentActiveExpireAfterDays := rd.GetChange("assignment.0.active.0.expire_after_days")

				if !assignmentActiveAllowPermanent.(bool) && assignmentActiveExpireAfterHours.(int) == 0 && assignmentActiveExpireAfterDays.(int) == 0 {
					return fmt.Errorf("cannot set allow permanent to false when expire after days and hours is 0 on `assignment` for `active`")
				}
			}

			return nil
		},
	}
}

func notificationConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"default_recipients": {
				Type:        pluginsdk.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Will notifications be sent to the default recipients?",
			},
			"additional_recipients": {
				Type:        pluginsdk.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "List of additional recipients to email notifications",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"critical_emails_only": {
				Type:        pluginsdk.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Will critical emails only be sent?",
			},
		},
	}
}

func (r RoleManagementPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RoleManagementPolicyResource) ModelObject() interface{} {
	return &RoleManagementPolicyResourceSchema{}
}

func (r RoleManagementPolicyResource) ResourceType() string {
	return "azurerm_role_management_policy"
}

func (r RoleManagementPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: r.createUpdate,
	}
}

func (r RoleManagementPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: r.createUpdate,
	}
}

func (r RoleManagementPolicyResource) createUpdate(ctx context.Context, metadata sdk.ResourceMetaData) error {
	client := metadata.Client.Authorization.RoleManagementPoliciesClient

	var id rolemanagementpolicies.ScopedRoleManagementPolicyId

	roleDefinitionId := metadata.ResourceData.Get("role_definition_id").(string)
	scope := metadata.ResourceData.Get("scope").(string)

	lockId := fmt.Sprintf("%s|%s", scope, roleDefinitionId)

	locks.ByName(lockId, r.ResourceType())
	defer locks.UnlockByName(lockId, r.ResourceType())

	scopeId := &commonids.ScopeId{
		Scope: scope,
	}
	filter := "roleDefinitionId%20eq%20'" + roleDefinitionId + "'"

	// filter by role definition
	result, err := client.ListForScopeComplete(ctx, *scopeId, filter)
	if err != nil {
		return fmt.Errorf("loading finding role management policy %q: %+v", scopeId, err)
	}

	if len(result.Items) != 1 {
		return fmt.Errorf("loading finding role management policy %q: %+v", scopeId, err)
	}

	roleManagementPolicyId := *result.Items[0].Name
	id = rolemanagementpolicies.NewScopedRoleManagementPolicyID(scope, roleManagementPolicyId)

	var config RoleManagementPolicyResourceSchema
	if err := metadata.Decode(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	var payload rolemanagementpolicies.RoleManagementPolicy

	if err := r.mapRoleManagementPolicyResourceSchemaToRoleManagementPolicy(config, &payload); err != nil {
		return fmt.Errorf("mapping schema model to sdk model: %+v", err)
	}

	res, err := client.Update(ctx, id, payload)
	if err != nil {
		return fmt.Errorf("updating %q: %+v", id.ID(), err)
	}

	if res.Model == nil {
		return fmt.Errorf("could not get role management policy")
	}
	policyID := res.Model.Name
	stateId := parse.NewRoleManagementPolicyID(scope, *policyID, roleDefinitionId)

	metadata.ResourceData.SetId(stateId.ID())

	return nil
}

func (r RoleManagementPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleManagementPoliciesClient
			schema := RoleManagementPolicyResourceSchema{}

			id, err := parse.RoleManagementPolicyId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ScopedRoleManagementPolicyId())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] Role Management Policy %q was not found - removing from state", id)
					err = metadata.MarkAsGone(id)
					if err != nil {
						return err
					}

					return nil
				}

				return fmt.Errorf("loading Role Management Policy %q: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				schema.Scope = id.Scope
				schema.RoleDefinitionID = id.RoleDefinitionId

				var config RoleManagementPolicyResourceSchema
				if err := metadata.Decode(&config); err != nil {
					return fmt.Errorf("decoding: %+v", err)
				}

				if err := r.mapRoleManagementPolicyToRoleManagementPolicyResourceSchema(*model, &schema, &config); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r RoleManagementPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// nothing to delete
			id, err := parse.RoleManagementPolicyId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("error parsing role management id: %v", err)
			}

			err = metadata.MarkAsGone(id)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func (r RoleManagementPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidateRoleManagementPolicyId
}

// functions to convert data from terraform structs to azure structs

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaToRoleManagementPolicy(input RoleManagementPolicyResourceSchema, output *rolemanagementpolicies.RoleManagementPolicy) error {

	if output.Properties == nil {
		output.Properties = &rolemanagementpolicies.RoleManagementPolicyProperties{}
	}
	if err := r.mapRoleManagementPolicyResourceSchemaToProperties(input, output.Properties); err != nil {
		return err
	}
	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaToProperties(input RoleManagementPolicyResourceSchema, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	output.Scope = &input.Scope

	if output.Rules == nil {
		rules := make([]rolemanagementpolicies.RoleManagementPolicyRule, 0)
		output.Rules = &rules
	}

	if len(input.Activation) == 1 {
		if err := r.mapRoleManagementPolicyResourceActivationSchemaToRoleManagementPolicyProperties(input.Activation[0], output); err != nil {
			return err
		}
	}
	if len(input.Assignment) == 1 {
		if err := r.mapRoleManagementPolicyResourceAssignmentSchemaToRoleManagementPolicyProperties(input.Assignment[0], output); err != nil {
			return err
		}
	}
	if len(input.Notifications) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationsToRoleManagementPolicyProperties(input.Notifications[0], output); err != nil {
			return err
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceActivationSchemaToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaActivation, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	if input.MaximumDurationHours > 0 {
		expirationRule := ExpirationRule{
			Id:                   "Expiration_EndUser_Assignment",
			RuleType:             "RoleManagementPolicyExpirationRule",
			MaximumDuration:      fmt.Sprintf("PT%dH", input.MaximumDurationHours),
			IsExpirationRequired: true,
			Target: &Target{
				Caller:              "EndUser",
				Operations:          []string{"All"},
				Level:               "Eligibility",
				TargetObjects:       nil,
				InheritableSettings: nil,
				EnforceSettings:     nil,
			},
		}
		*output.Rules = append(*output.Rules, expirationRule)
	}

	enabledRules := make([]string, 0)

	if input.RequireJustification {
		enabledRules = append(enabledRules, "Justification")
	}
	if input.RequireMultiFactorAuthentication {
		enabledRules = append(enabledRules, "MultiFactorAuthentication")
	}
	if input.RequireTicketInformation {
		enabledRules = append(enabledRules, "Ticketing")
	}

	enablementRule := EnablementRule{
		Id:       "Enablement_EndUser_Assignment",
		RuleType: "RoleManagementPolicyEnablementRule",

		EnabledRules: enabledRules,

		Target: &Target{
			Caller:              "EndUser",
			Operations:          []string{"All"},
			Level:               "Assignment",
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, enablementRule)

	primaryApprovers := make([]PrimaryApprovers, 0)

	if input.Approvers != nil && len(input.Approvers) > 0 {
		approvers := input.Approvers[0]

		if approvers.Groups != nil {
			for _, approver := range approvers.Groups {
				primaryApprovers = append(primaryApprovers, PrimaryApprovers{
					Id:          approver.Id,
					Description: approver.Name,
					IsBackup:    false,
					UserType:    "Group",
				})
			}
		}

		if approvers.Users != nil {
			for _, approver := range approvers.Users {
				primaryApprovers = append(primaryApprovers, PrimaryApprovers{
					Id:          approver.Id,
					Description: approver.Name,
					IsBackup:    false,
					UserType:    "User",
				})
			}
		}
	}

	isApprovalRequired := len(primaryApprovers) > 0
	approvalStages := make([]ApprovalStages, 1)
	approvalStages[0] = ApprovalStages{
		ApprovalStageTimeOutInDays:      1,
		EscalationTimeInMinutes:         0,
		IsEscalationEnabled:             false,
		IsApproverJustificationRequired: true,
		PrimaryApprovers:                &primaryApprovers,
	}
	approvalRule := ApprovalRule{
		Id:       "Approval_EndUser_Assignment",
		RuleType: "RoleManagementPolicyApprovalRule",
		Setting: &Settings{
			ApprovalMode:                     "SingleStage",
			IsApprovalRequired:               isApprovalRequired,
			IsApprovalRequiredForExtension:   false,
			IsRequestorJustificationRequired: true,
			ApprovalStages:                   &approvalStages,
		},
		Target: &Target{
			Caller:              "EndUser",
			Operations:          []string{"All"},
			Level:               "Assignment",
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, approvalRule)

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceAssignmentSchemaToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaAssignment, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	if len(input.Active) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaAssignmentActiveToRoleManagementPolicyProperties(input.Active[0], output); err != nil {
			return err
		}
	}

	if len(input.Eligible) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaAssignmentEligibleToRoleManagementPolicyProperties(input.Eligible[0], output); err != nil {
			return err
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaAssignmentActiveToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaAssignmentActive, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	maximumDuration := ""
	if input.ExpireAfterDays > 0 {
		maximumDuration = fmt.Sprintf("P%dD", input.ExpireAfterDays)
	} else if input.ExpireAfterHours > 0 {
		maximumDuration = fmt.Sprintf("PT%dH", input.ExpireAfterHours)
	}
	expirationRule := ExpirationRule{
		Id:                   "Expiration_Admin_Assignment",
		RuleType:             "RoleManagementPolicyExpirationRule",
		MaximumDuration:      maximumDuration,
		IsExpirationRequired: !input.AllowPermanent,
		Target: &Target{
			Caller:              "Admin",
			Operations:          []string{"All"},
			Level:               "Assignment",
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, expirationRule)

	enabledRules := make([]string, 0)

	if input.RequireJustification {
		enabledRules = append(enabledRules, "Justification")
	}
	if input.RequireMultiFactorAuthentication {
		enabledRules = append(enabledRules, "MultiFactorAuthentication")
	}

	enablementRule := EnablementRule{
		Id:       "Enablement_Admin_Assignment",
		RuleType: "RoleManagementPolicyEnablementRule",

		EnabledRules: enabledRules,

		Target: &Target{
			Caller:              "Admin",
			Operations:          []string{"All"},
			Level:               "Assignment",
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, enablementRule)

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaAssignmentEligibleToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaAssignmentEligible, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	maximumDuration := ""
	if input.ExpireAfterDays > 0 {
		maximumDuration = fmt.Sprintf("P%dD", input.ExpireAfterDays)
	} else if input.ExpireAfterHours > 0 {
		maximumDuration = fmt.Sprintf("PT%dH", input.ExpireAfterHours)
	}
	expirationRule := ExpirationRule{
		Id:                   "Expiration_Admin_Eligibility",
		RuleType:             "RoleManagementPolicyExpirationRule",
		MaximumDuration:      maximumDuration,
		IsExpirationRequired: !input.AllowPermanent,
		Target: &Target{
			Caller:              "Admin",
			Operations:          []string{"All"},
			Level:               "Eligibility",
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, expirationRule)

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaNotificationsToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaNotifications, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	if len(input.MembersAssignedEligible) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationsEligibleToRoleManagementPolicyProperties(input.MembersAssignedEligible[0], "Eligibility", output); err != nil {
			return err
		}
	}
	if len(input.MembersAssignedActive) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationsEligibleToRoleManagementPolicyProperties(input.MembersAssignedActive[0], "Assignment", output); err != nil {
			return err
		}
	}
	if len(input.EligibleMemberActivate) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationsEligibleMemberActivateToRoleManagementPolicyProperties(input.EligibleMemberActivate[0], output); err != nil {
			return err
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaNotificationsEligibleToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaNotificationsEligible, level string, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	if len(input.RoleAssignmentAlert) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input.RoleAssignmentAlert[0], "Admin", "Admin", level, output); err != nil {
			return err
		}
	}
	if len(input.AssignedUser) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input.AssignedUser[0], "Requestor", "Admin", level, output); err != nil {
			return err
		}
	}
	if len(input.RequestForExtensionOrApproval) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input.RequestForExtensionOrApproval[0], "Approver", "Admin", level, output); err != nil {
			return err
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaNotificationsEligibleMemberActivateToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaNotificationsActivate, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	if len(input.RoleAssignmentAlert) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input.RoleAssignmentAlert[0], "Admin", "EndUser", "Assignment", output); err != nil {
			return err
		}
	}
	if len(input.AssignedUser) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input.AssignedUser[0], "Requestor", "EndUser", "Assignment", output); err != nil {
			return err
		}
	}
	if len(input.RequestForExtensionOrApproval) == 1 {
		if err := r.mapRoleManagementPolicyResourceSchemaNotificationWithoutRecipientsToRoleManagementPolicyProperties(input.RequestForExtensionOrApproval[0], "Approver", "EndUser", "Assignment", output); err != nil {
			return err
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaNotificationToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaNotification, recipientType string, caller string, level string, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	id := fmt.Sprintf("Notification_%s_%s_%s", recipientType, caller, level)
	notificationLevel := "All"
	if input.CriticalEmailsOnly {
		notificationLevel = "Critical"
	}

	notificationRule := NotificationRule{
		Id:                         id,
		RuleType:                   "RoleManagementPolicyNotificationRule",
		NotificationRecipients:     input.AdditionalRecipients,
		RecipientType:              recipientType,
		NotificationLevel:          notificationLevel,
		NotificationType:           "Email",
		IsDefaultRecipientsEnabled: input.DefaultRecipients,
		Target: &Target{
			Caller:              caller,
			Operations:          []string{"All"},
			Level:               level,
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, notificationRule)

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyResourceSchemaNotificationWithoutRecipientsToRoleManagementPolicyProperties(input RoleManagementPolicyResourceSchemaNotificationWithoutRecipients, recipientType string, caller string, level string, output *rolemanagementpolicies.RoleManagementPolicyProperties) error {
	id := fmt.Sprintf("Notification_%s_%s_%s", recipientType, caller, level)
	notificationLevel := "All"
	if input.CriticalEmailsOnly {
		notificationLevel = "Critical"
	}

	notificationRule := NotificationRule{
		Id:                         id,
		RuleType:                   "RoleManagementPolicyNotificationRule",
		RecipientType:              recipientType,
		NotificationLevel:          notificationLevel,
		IsDefaultRecipientsEnabled: input.DefaultRecipients,
		NotificationType:           "Email",
		Target: &Target{
			Caller:              caller,
			Operations:          []string{"All"},
			Level:               level,
			TargetObjects:       nil,
			InheritableSettings: nil,
			EnforceSettings:     nil,
		},
	}
	*output.Rules = append(*output.Rules, notificationRule)

	return nil
}

// functions to convert data from azure structs to terraform structs

func (r RoleManagementPolicyResource) mapRoleManagementPolicyToRoleManagementPolicyResourceSchema(input rolemanagementpolicies.RoleManagementPolicy, output *RoleManagementPolicyResourceSchema, config *RoleManagementPolicyResourceSchema) error {
	if input.Properties == nil {
		input.Properties = &rolemanagementpolicies.RoleManagementPolicyProperties{}
	}
	if err := r.mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchema(input.Properties, output, config); err != nil {
		return err
	}
	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchema(input *rolemanagementpolicies.RoleManagementPolicyProperties, output *RoleManagementPolicyResourceSchema, config *RoleManagementPolicyResourceSchema) error {
	if config.Activation != nil && len(config.Activation) == 1 {
		tmpActivation := &RoleManagementPolicyResourceSchemaActivation{}
		if err := r.mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaActivation(input, tmpActivation, &config.Activation[0]); err != nil {
			return err
		} else {
			output.Activation = make([]RoleManagementPolicyResourceSchemaActivation, 0)
			output.Activation = append(output.Activation, *tmpActivation)
		}
	}

	if config.Assignment != nil && len(config.Assignment) == 1 {
		tmpAssignment := &RoleManagementPolicyResourceSchemaAssignment{}
		if err := r.mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaAssignment(input, tmpAssignment); err != nil {
			return err
		} else {
			output.Assignment = make([]RoleManagementPolicyResourceSchemaAssignment, 0)
			output.Assignment = append(output.Assignment, *tmpAssignment)
		}
	}

	if config.Notifications != nil && len(config.Notifications) == 1 {
		tmpNotifications := &RoleManagementPolicyResourceSchemaNotifications{}
		if err := r.mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaNotifications(input, tmpNotifications, &config.Notifications[0]); err != nil {
			return err
		} else {
			output.Notifications = make([]RoleManagementPolicyResourceSchemaNotifications, 0)
			output.Notifications = append(output.Notifications, *tmpNotifications)
		}
	}

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaActivation(input *rolemanagementpolicies.RoleManagementPolicyProperties, output *RoleManagementPolicyResourceSchemaActivation, config *RoleManagementPolicyResourceSchemaActivation) error {
	if input.Rules != nil {
		for _, r := range *input.Rules {
			rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
			switch rule.Values["id"].(string) {
			case "Expiration_EndUser_Assignment":
				if config.MaximumDurationHours == 0 {
					continue
				}

				maximumDurationRaw := rule.Values["maximumDuration"].(string)
				re := regexp.MustCompile(`\d+`)
				maxDuration, err := strconv.Atoi(re.FindString(maximumDurationRaw))
				if err != nil {
					return err
				}
				output.MaximumDurationHours = maxDuration

			case "Enablement_EndUser_Assignment":
				enabledRules := rule.Values["enabledRules"].([]interface{})
				for _, r := range enabledRules {
					switch r {
					case "Justification":
						output.RequireJustification = true

					case "MultiFactorAuthentication":
						output.RequireMultiFactorAuthentication = true

					case "Ticketing":
						output.RequireTicketInformation = true
					}
				}
			case "Approval_EndUser_Assignment":
				if config.Approvers == nil || len(config.Approvers) == 0 {
					continue
				}

				output.Approvers = make([]RoleManagementPolicyResourceSchemaApprovers, 1)
				approvers := &RoleManagementPolicyResourceSchemaApprovers{}
				approvers.Groups = make([]RoleManagementPolicyResourceSchemaApproversApprover, 0)
				approvers.Users = make([]RoleManagementPolicyResourceSchemaApproversApprover, 0)

				setting := rule.Values["setting"].(map[string]interface{})
				approvalStages := setting["approvalStages"].([]interface{})
				if len(approvalStages) == 1 {
					approvalStage := approvalStages[0].(map[string]interface{})
					if approvalStage["primaryApprovers"] != nil {
						for _, approver := range approvalStage["primaryApprovers"].([]interface{}) {
							approverMap := approver.(map[string]interface{})
							switch approverMap["userType"] {
							case "User":
								approvers.Users = append(approvers.Users, RoleManagementPolicyResourceSchemaApproversApprover{
									Id:   approverMap["id"].(string),
									Name: approverMap["description"].(string),
								})
							case "Group":
								approvers.Groups = append(approvers.Groups, RoleManagementPolicyResourceSchemaApproversApprover{
									Id:   approverMap["id"].(string),
									Name: approverMap["description"].(string),
								})
							}
						}
					}
				}
				output.Approvers[0] = *approvers
			}
		}
	}
	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaAssignment(input *rolemanagementpolicies.RoleManagementPolicyProperties, output *RoleManagementPolicyResourceSchemaAssignment) error {
	tmpActive := &RoleManagementPolicyResourceSchemaAssignmentActive{}
	tmpEligible := &RoleManagementPolicyResourceSchemaAssignmentEligible{}

	if input.Rules != nil {
		for _, r := range *input.Rules {
			rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)
			switch rule.Values["id"].(string) {
			case "Expiration_Admin_Eligibility":
				maximumDurationRaw := rule.Values["maximumDuration"].(string)
				reHours := regexp.MustCompile(`PT(\d+)H`)
				hours := reHours.FindStringSubmatch(maximumDurationRaw)
				if len(hours) == 2 {
					maxDurationHours, err := strconv.Atoi(hours[1])
					if err != nil {
						return err
					}
					tmpEligible.ExpireAfterHours = maxDurationHours
				}
				reDays := regexp.MustCompile(`P(\d+)D`)
				days := reDays.FindStringSubmatch(maximumDurationRaw)
				if len(days) == 2 {
					maxDurationDays, err := strconv.Atoi(days[1])
					if err != nil {
						return err
					}
					tmpEligible.ExpireAfterDays = maxDurationDays
				}

				tmpEligible.AllowPermanent = !rule.Values["isExpirationRequired"].(bool)

			case "Expiration_Admin_Assignment":
				maximumDurationRaw := rule.Values["maximumDuration"].(string)
				reHours := regexp.MustCompile(`PT(\d+)H`)
				hours := reHours.FindStringSubmatch(maximumDurationRaw)
				if len(hours) == 2 {
					maxDurationHours, err := strconv.Atoi(hours[1])
					if err != nil {
						return err
					}
					tmpActive.ExpireAfterHours = maxDurationHours
				}
				reDays := regexp.MustCompile(`P(\d+)D`)
				days := reDays.FindStringSubmatch(maximumDurationRaw)
				if len(days) == 2 {
					maxDurationDays, err := strconv.Atoi(days[1])
					if err != nil {
						return err
					}
					tmpActive.ExpireAfterDays = maxDurationDays
				}

				tmpActive.AllowPermanent = !rule.Values["isExpirationRequired"].(bool)

			case "Enablement_Admin_Assignment":
				enabledRules := rule.Values["enabledRules"].([]interface{})
				for _, r := range enabledRules {
					switch r {
					case "Justification":
						tmpActive.RequireJustification = true

					case "MultiFactorAuthentication":
						tmpActive.RequireMultiFactorAuthentication = true
					}
				}

			}
		}
	}

	output.Active = make([]RoleManagementPolicyResourceSchemaAssignmentActive, 0)
	output.Active = append(output.Active, *tmpActive)

	output.Eligible = make([]RoleManagementPolicyResourceSchemaAssignmentEligible, 0)
	output.Eligible = append(output.Eligible, *tmpEligible)

	return nil
}

func (r RoleManagementPolicyResource) mapRoleManagementPolicyPropertiesToRoleManagementPolicyResourceSchemaNotifications(input *rolemanagementpolicies.RoleManagementPolicyProperties, output *RoleManagementPolicyResourceSchemaNotifications, config *RoleManagementPolicyResourceSchemaNotifications) error {
	var tmpMembersAssignedEligible *RoleManagementPolicyResourceSchemaNotificationsEligible
	var tmpMembersAssignedActive *RoleManagementPolicyResourceSchemaNotificationsEligible
	var tmpEligibleMemberActivate *RoleManagementPolicyResourceSchemaNotificationsActivate

	if config.MembersAssignedEligible != nil && len(config.MembersAssignedEligible) == 1 {
		tmpMembersAssignedEligible = &RoleManagementPolicyResourceSchemaNotificationsEligible{}
	}
	if config.MembersAssignedActive != nil && len(config.MembersAssignedActive) == 1 {
		tmpMembersAssignedActive = &RoleManagementPolicyResourceSchemaNotificationsEligible{}
	}
	if config.EligibleMemberActivate != nil && len(config.EligibleMemberActivate) == 1 {
		tmpEligibleMemberActivate = &RoleManagementPolicyResourceSchemaNotificationsActivate{}
	}

	if input.Rules != nil {
		for _, r := range *input.Rules {
			rule := r.(rolemanagementpolicies.RawRoleManagementPolicyRuleImpl)

			if rule.Values["ruleType"].(string) != "RoleManagementPolicyNotificationRule" {
				continue
			}

			notificationLevel := rule.Values["notificationLevel"].(string)
			criticalEmailsOnly := false
			if notificationLevel == "Critical" {
				criticalEmailsOnly = true
			}

			notificationRecipients := rule.Values["notificationRecipients"]
			additionalRecipients := make([]string, 0)
			if notificationRecipients != nil {
				for _, rec := range notificationRecipients.([]interface{}) {
					additionalRecipients = append(additionalRecipients, rec.(string))
				}
			}

			switch rule.Values["id"].(string) {

			case "Notification_Admin_Admin_Eligibility":
				if config.MembersAssignedEligible != nil && len(config.MembersAssignedEligible) == 1 && config.MembersAssignedEligible[0].RoleAssignmentAlert != nil && len(config.MembersAssignedEligible[0].RoleAssignmentAlert) == 1 {
					tmpMembersAssignedEligible.RoleAssignmentAlert = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedEligible.RoleAssignmentAlert[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedEligible.RoleAssignmentAlert[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedEligible.RoleAssignmentAlert[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Requestor_Admin_Eligibility":
				if config.MembersAssignedEligible != nil && len(config.MembersAssignedEligible) == 1 && config.MembersAssignedEligible[0].AssignedUser != nil && len(config.MembersAssignedEligible[0].AssignedUser) == 1 {
					tmpMembersAssignedEligible.AssignedUser = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedEligible.AssignedUser[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedEligible.AssignedUser[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedEligible.AssignedUser[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Approver_Admin_Eligibility":
				if config.MembersAssignedEligible != nil && len(config.MembersAssignedEligible) == 1 && config.MembersAssignedEligible[0].RequestForExtensionOrApproval != nil && len(config.MembersAssignedEligible[0].RequestForExtensionOrApproval) == 1 {
					tmpMembersAssignedEligible.RequestForExtensionOrApproval = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedEligible.RequestForExtensionOrApproval[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedEligible.RequestForExtensionOrApproval[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedEligible.RequestForExtensionOrApproval[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Admin_Admin_Assignment":
				if config.MembersAssignedActive != nil && len(config.MembersAssignedActive) == 1 && config.MembersAssignedActive[0].RoleAssignmentAlert != nil && len(config.MembersAssignedActive[0].RoleAssignmentAlert) == 1 {
					tmpMembersAssignedActive.RoleAssignmentAlert = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedActive.RoleAssignmentAlert[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedActive.RoleAssignmentAlert[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedActive.RoleAssignmentAlert[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Requestor_Admin_Assignment":
				if config.MembersAssignedActive != nil && len(config.MembersAssignedActive) == 1 && config.MembersAssignedActive[0].AssignedUser != nil && len(config.MembersAssignedActive[0].AssignedUser) == 1 {
					tmpMembersAssignedActive.AssignedUser = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedActive.AssignedUser[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedActive.AssignedUser[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedActive.AssignedUser[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Approver_Admin_Assignment":
				if config.MembersAssignedActive != nil && len(config.MembersAssignedActive) == 1 && config.MembersAssignedActive[0].RequestForExtensionOrApproval != nil && len(config.MembersAssignedActive[0].RequestForExtensionOrApproval) == 1 {
					tmpMembersAssignedActive.RequestForExtensionOrApproval = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpMembersAssignedActive.RequestForExtensionOrApproval[0].AdditionalRecipients = additionalRecipients
					tmpMembersAssignedActive.RequestForExtensionOrApproval[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpMembersAssignedActive.RequestForExtensionOrApproval[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Admin_EndUser_Assignment":
				if config.EligibleMemberActivate != nil && len(config.EligibleMemberActivate) == 1 && config.EligibleMemberActivate[0].RoleAssignmentAlert != nil && len(config.EligibleMemberActivate[0].RoleAssignmentAlert) == 1 {
					tmpEligibleMemberActivate.RoleAssignmentAlert = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpEligibleMemberActivate.RoleAssignmentAlert[0].AdditionalRecipients = additionalRecipients
					tmpEligibleMemberActivate.RoleAssignmentAlert[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpEligibleMemberActivate.RoleAssignmentAlert[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Requestor_EndUser_Assignment":
				if config.EligibleMemberActivate != nil && len(config.EligibleMemberActivate) == 1 && config.EligibleMemberActivate[0].AssignedUser != nil && len(config.EligibleMemberActivate[0].AssignedUser) == 1 {
					tmpEligibleMemberActivate.AssignedUser = make([]RoleManagementPolicyResourceSchemaNotification, 1)
					tmpEligibleMemberActivate.AssignedUser[0].AdditionalRecipients = additionalRecipients
					tmpEligibleMemberActivate.AssignedUser[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpEligibleMemberActivate.AssignedUser[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			case "Notification_Approver_EndUser_Assignment":
				if config.EligibleMemberActivate != nil && len(config.EligibleMemberActivate) == 1 && config.EligibleMemberActivate[0].RequestForExtensionOrApproval != nil && len(config.EligibleMemberActivate[0].RequestForExtensionOrApproval) == 1 {
					tmpEligibleMemberActivate.RequestForExtensionOrApproval = make([]RoleManagementPolicyResourceSchemaNotificationWithoutRecipients, 1)
					tmpEligibleMemberActivate.RequestForExtensionOrApproval[0].CriticalEmailsOnly = criticalEmailsOnly
					tmpEligibleMemberActivate.RequestForExtensionOrApproval[0].DefaultRecipients = rule.Values["isDefaultRecipientsEnabled"].(bool)
				}

			}

		}
	}

	if tmpMembersAssignedEligible != nil {
		output.MembersAssignedEligible = make([]RoleManagementPolicyResourceSchemaNotificationsEligible, 0)
		output.MembersAssignedEligible = append(output.MembersAssignedEligible, *tmpMembersAssignedEligible)
	}
	if tmpMembersAssignedActive != nil {
		output.MembersAssignedActive = make([]RoleManagementPolicyResourceSchemaNotificationsEligible, 0)
		output.MembersAssignedActive = append(output.MembersAssignedActive, *tmpMembersAssignedActive)
	}
	if tmpEligibleMemberActivate != nil {
		output.EligibleMemberActivate = make([]RoleManagementPolicyResourceSchemaNotificationsActivate, 0)
		output.EligibleMemberActivate = append(output.EligibleMemberActivate, *tmpEligibleMemberActivate)
	}
	return nil
}

type RoleManagementPolicyResourceSchema struct {
	RoleDefinitionID string                                            `tfschema:"role_definition_id"`
	Scope            string                                            `tfschema:"scope"`
	Activation       []RoleManagementPolicyResourceSchemaActivation    `tfschema:"activation"`
	Assignment       []RoleManagementPolicyResourceSchemaAssignment    `tfschema:"assignment"`
	Notifications    []RoleManagementPolicyResourceSchemaNotifications `tfschema:"notifications"`
}

type RoleManagementPolicyResourceSchemaActivation struct {
	MaximumDurationHours             int                                           `tfschema:"maximum_duration_hours"`
	RequireMultiFactorAuthentication bool                                          `tfschema:"require_multi_factor_authentication"`
	RequireJustification             bool                                          `tfschema:"require_justification"`
	RequireTicketInformation         bool                                          `tfschema:"require_ticket_information"`
	Approvers                        []RoleManagementPolicyResourceSchemaApprovers `tfschema:"approvers"`
}

type RoleManagementPolicyResourceSchemaApprovers struct {
	Groups []RoleManagementPolicyResourceSchemaApproversApprover `tfschema:"group"`
	Users  []RoleManagementPolicyResourceSchemaApproversApprover `tfschema:"user"`
}

type RoleManagementPolicyResourceSchemaApproversApprover struct {
	Id   string `tfschema:"id"`
	Name string `tfschema:"name"`
}

type RoleManagementPolicyResourceSchemaAssignment struct {
	Active   []RoleManagementPolicyResourceSchemaAssignmentActive   `tfschema:"active"`
	Eligible []RoleManagementPolicyResourceSchemaAssignmentEligible `tfschema:"eligible"`
}

type RoleManagementPolicyResourceSchemaAssignmentActive struct {
	AllowPermanent                   bool `tfschema:"allow_permanent"`
	ExpireAfterDays                  int  `tfschema:"expire_after_days"`
	ExpireAfterHours                 int  `tfschema:"expire_after_hours"`
	RequireMultiFactorAuthentication bool `tfschema:"require_multi_factor_authentication"`
	RequireJustification             bool `tfschema:"require_justification"`
}

type RoleManagementPolicyResourceSchemaAssignmentEligible struct {
	AllowPermanent   bool `tfschema:"allow_permanent"`
	ExpireAfterDays  int  `tfschema:"expire_after_days"`
	ExpireAfterHours int  `tfschema:"expire_after_hours"`
}

type RoleManagementPolicyResourceSchemaNotifications struct {
	MembersAssignedEligible []RoleManagementPolicyResourceSchemaNotificationsEligible `tfschema:"member_assigned_eligible"`
	MembersAssignedActive   []RoleManagementPolicyResourceSchemaNotificationsEligible `tfschema:"member_assigned_active"`
	EligibleMemberActivate  []RoleManagementPolicyResourceSchemaNotificationsActivate `tfschema:"eligible_member_activate"`
}

type RoleManagementPolicyResourceSchemaNotificationsEligible struct {
	RoleAssignmentAlert           []RoleManagementPolicyResourceSchemaNotification `tfschema:"role_assignment_alert"`
	AssignedUser                  []RoleManagementPolicyResourceSchemaNotification `tfschema:"assigned_user"`
	RequestForExtensionOrApproval []RoleManagementPolicyResourceSchemaNotification `tfschema:"request_for_extension_or_approval"`
}

type RoleManagementPolicyResourceSchemaNotificationsActivate struct {
	RoleAssignmentAlert           []RoleManagementPolicyResourceSchemaNotification                  `tfschema:"role_assignment_alert"`
	AssignedUser                  []RoleManagementPolicyResourceSchemaNotification                  `tfschema:"assigned_user"`
	RequestForExtensionOrApproval []RoleManagementPolicyResourceSchemaNotificationWithoutRecipients `tfschema:"request_for_extension_or_approval"`
}

type RoleManagementPolicyResourceSchemaNotification struct {
	AdditionalRecipients []string `tfschema:"additional_recipients"`
	CriticalEmailsOnly   bool     `tfschema:"critical_emails_only"`
	DefaultRecipients    bool     `tfschema:"default_recipients"`
}

type RoleManagementPolicyResourceSchemaNotificationWithoutRecipients struct {
	CriticalEmailsOnly bool `tfschema:"critical_emails_only"`
	DefaultRecipients  bool `tfschema:"default_recipients"`
}

type ApprovalRule struct {
	Id       string  `json:"id"`
	RuleType string  `json:"ruleType"`
	Target   *Target `json:"target"`

	Setting *Settings `json:"setting"`
}

type ExpirationRule struct {
	Id       string  `json:"id"`
	RuleType string  `json:"ruleType"`
	Target   *Target `json:"target"`

	IsExpirationRequired bool   `json:"isExpirationRequired"`
	MaximumDuration      string `json:"maximumDuration"`
}

type EnablementRule struct {
	Id       string  `json:"id"`
	RuleType string  `json:"ruleType"`
	Target   *Target `json:"target"`

	EnabledRules []string `json:"enabledRules"`
}

type NotificationRule struct {
	Id       string  `json:"id"`
	RuleType string  `json:"ruleType"`
	Target   *Target `json:"target"`

	NotificationType           string   `json:"notificationType"`
	RecipientType              string   `json:"recipientType"`
	IsDefaultRecipientsEnabled bool     `json:"isDefaultRecipientsEnabled"`
	NotificationLevel          string   `json:"notificationLevel"`
	NotificationRecipients     []string `json:"notificationRecipients"`
}

type Target struct {
	Caller              string      `json:"caller"`
	Operations          []string    `json:"operations"`
	Level               string      `json:"level"`
	TargetObjects       interface{} `json:"targetObjects"`
	InheritableSettings interface{} `json:"inheritableSettings"`
	EnforceSettings     interface{} `json:"enforceSettings"`
}

type Settings struct {
	ApprovalMode                     string            `json:"approvalMode"`
	IsApprovalRequired               bool              `json:"isApprovalRequired"`
	IsApprovalRequiredForExtension   bool              `json:"isApprovalRequiredForExtension"`
	IsRequestorJustificationRequired bool              `json:"isRequestorJustificationRequired"`
	ApprovalStages                   *[]ApprovalStages `json:"approvalStages"`
}

type ApprovalStages struct {
	ApprovalStageTimeOutInDays      int                 `json:"approvalStageTimeOutInDays"`
	EscalationTimeInMinutes         int                 `json:"escalationTimeInMinutes"`
	IsApproverJustificationRequired bool                `json:"isApproverJustificationRequired"`
	IsEscalationEnabled             bool                `json:"isEscalationEnabled"`
	PrimaryApprovers                *[]PrimaryApprovers `json:"primaryApprovers"`
}

type PrimaryApprovers struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	IsBackup    bool   `json:"isBackup"`
	UserType    string `json:"userType"`
}
