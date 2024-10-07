// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	// nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedules"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	billingValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = PimEligibleRoleAssignmentResource{}

type PimEligibleRoleAssignmentResource struct{}

type PimEligibleRoleAssignmentModel struct {
	RoleDefinitionId string                                  `tfschema:"role_definition_id"`
	Scope            string                                  `tfschema:"scope"`
	PrincipalId      string                                  `tfschema:"principal_id"`
	PrincipalType    string                                  `tfschema:"principal_type"`
	Justification    string                                  `tfschema:"justification"`
	TicketInfo       []PimEligibleRoleAssignmentTicketInfo   `tfschema:"ticket"`
	ScheduleInfo     []PimEligibleRoleAssignmentScheduleInfo `tfschema:"schedule"`
}

type PimEligibleRoleAssignmentTicketInfo struct {
	TicketNumber string `tfschema:"number"`
	TicketSystem string `tfschema:"system"`
}

type PimEligibleRoleAssignmentScheduleInfo struct {
	StartDateTime string                                            `tfschema:"start_date_time"`
	Expiration    []PimEligibleRoleAssignmentScheduleInfoExpiration `tfschema:"expiration"`
}

type PimEligibleRoleAssignmentScheduleInfoExpiration struct {
	DurationDays  int64  `tfschema:"duration_days"`
	DurationHours int64  `tfschema:"duration_hours"`
	EndDateTime   string `tfschema:"end_date_time"`
}

func (PimEligibleRoleAssignmentResource) ModelObject() interface{} {
	return &PimEligibleRoleAssignmentModel{}
}

func (PimEligibleRoleAssignmentResource) ResourceType() string {
	return "azurerm_pim_eligible_role_assignment"
}

func (PimEligibleRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.PimRoleAssignmentID
}

func (PimEligibleRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Scope for this eligible role assignment, should be a valid resource ID",
			ValidateFunc: validation.Any(
				// Elevated access for a global admin is needed to assign roles in this scope:
				// https://docs.microsoft.com/en-us/azure/role-based-access-control/elevate-access-global-admin#azure-cli
				// It seems only user account is allowed to be elevated access.
				validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),

				billingValidate.EnrollmentID,
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
			Computed:    true,
			ForceNew:    true,
			Description: "The justification for this eligible role assignment",
		},

		"schedule": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: "The schedule details for this eligible role assignment",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date_time": { // defaults to now
						Optional:    true,
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
									Computed: true,
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
									Computed: true,
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
									Computed: true,
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
			Computed:    true,
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

func (PimEligibleRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "Type of principal to which the role will be assigned",
		},
	}
}

func (r PimEligibleRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			schedulesClient := metadata.Client.Authorization.RoleEligibilitySchedulesClient
			requestsClient := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			var config PimEligibleRoleAssignmentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewPimRoleAssignmentID(config.Scope, config.RoleDefinitionId, config.PrincipalId)

			schedule, err := findRoleEligibilitySchedule(ctx, schedulesClient, id)
			if err != nil {
				return err
			}
			if schedule != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			scheduleInfo := &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo{
				Expiration: &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration{
					Type: pointer.To(roleeligibilityschedulerequests.TypeNoExpiration),
				},
			}

			if len(config.ScheduleInfo) > 0 {
				if config.ScheduleInfo[0].StartDateTime != "" {
					scheduleInfo.StartDateTime = pointer.To(config.ScheduleInfo[0].StartDateTime)
				}

				if expiration := config.ScheduleInfo[0].Expiration; len(expiration) > 0 {
					scheduleInfo.Expiration = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration{
						Type: pointer.To(roleeligibilityschedulerequests.TypeNoExpiration),
					}

					switch {
					case expiration[0].DurationDays != 0:
						scheduleInfo.Expiration.Duration = pointer.To(fmt.Sprintf("P%dD", expiration[0].DurationDays))
						scheduleInfo.Expiration.Type = pointer.To(roleeligibilityschedulerequests.TypeAfterDuration)

					case expiration[0].DurationHours != 0:
						scheduleInfo.Expiration.Duration = pointer.To(fmt.Sprintf("PT%dH", expiration[0].DurationHours))
						scheduleInfo.Expiration.Type = pointer.To(roleeligibilityschedulerequests.TypeAfterDuration)

					case expiration[0].EndDateTime != "":
						scheduleInfo.Expiration.EndDateTime = pointer.To(expiration[0].EndDateTime)
						scheduleInfo.Expiration.Type = pointer.To(roleeligibilityschedulerequests.TypeAfterDateTime)
					}
				}
			}

			var ticketInfo *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo

			if len(config.TicketInfo) > 0 {
				ticketInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo{
					TicketNumber: pointer.To(config.TicketInfo[0].TicketNumber),
					TicketSystem: pointer.To(config.TicketInfo[0].TicketSystem),
				}
			}

			scopeId, err := commonids.ParseScopeID(id.Scope)
			if err != nil {
				return err
			}

			payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{
				Properties: &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{
					Justification:    pointer.To(config.Justification),
					PrincipalId:      id.PrincipalId,
					RequestType:      roleeligibilityschedulerequests.RequestTypeAdminAssign,
					RoleDefinitionId: id.RoleDefinitionId,
					Scope:            pointer.To(scopeId.ID()),
					ScheduleInfo:     scheduleInfo,
					TicketInfo:       ticketInfo,
				},
			}

			roleEligibilityScheduleRequestName, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating uuid: %+v", err)
			}

			requestId := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID(id.Scope, roleEligibilityScheduleRequestName)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}

			// TODO: Remove this WaitForState workaround once eventual-consistency retries are added to the upstream SDK
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Retry"},
				Target:  []string{"Created"},
				Refresh: func() (interface{}, string, error) {
					// Retry new requests to smooth over AAD replication issues with the subject principal
					result, err := requestsClient.Create(ctx, requestId, payload)
					if err != nil {
						if result.OData != nil && result.OData.Error != nil && result.OData.Error.Code != nil && *result.OData.Error.Code == "SubjectNotFound" {
							return result, "Retry", nil
						}

						return result, "Error", fmt.Errorf("creating %s: %+v", requestId, err)
					}

					return result, "Created", nil
				},
				MinTimeout: 10 * time.Second,
				Timeout:    time.Until(deadline),
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be created: %+v", requestId, err)
			}

			// Wait for the request to be processed and a schedule to be created, so that subsequent reads will succeed
			stateConf = &pluginsdk.StateChangeConf{
				Pending:    []string{"NotFound"},
				Target:     []string{"Exists"},
				Refresh:    pollForRoleEligibilitySchedule(ctx, schedulesClient, id),
				MinTimeout: 10 * time.Second,
				Timeout:    time.Until(deadline),
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become found: %+v", requestId, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PimEligibleRoleAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			schedulesClient := metadata.Client.Authorization.RoleEligibilitySchedulesClient
			requestsClient := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			// Retrieve existing state as we may not be able to populate everything after the initial request has expired
			var state PimEligibleRoleAssignmentModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Look for a Schedule for this eligible role assignment. Note that the Schedule ID changes each time the assignment is manipulated,
			// whilst still remaining valid for the configured role, scope and principal, so we must search for it.
			schedule, err := findRoleEligibilitySchedule(ctx, schedulesClient, *id)
			if err != nil {
				return err
			}
			if schedule == nil {
				return metadata.MarkAsGone(id)
			}

			// Look for the latest associated assignment request as this contains some fields we want
			request, err := findRoleEligibilityScheduleRequest(ctx, requestsClient, schedule, id)
			if err != nil {
				return err
			}

			state.Scope = id.Scope

			// PIM Role Assignments are represented by a Schedule object and one or more Request objects that comprise the audit history. Requests return
			// more information, but expire after 45 days, so after this time we can only partially populate the resource attributes from the Schedule.
			if request != nil && request.Properties != nil {
				// A request is still present and was found, so populate from the request
				state.Justification = pointer.From(request.Properties.Justification)
				state.PrincipalId = request.Properties.PrincipalId
				state.PrincipalType = string(pointer.From(request.Properties.PrincipalType))
				state.RoleDefinitionId = request.Properties.RoleDefinitionId

				if ticketInfo := request.Properties.TicketInfo; ticketInfo != nil {
					if len(state.TicketInfo) == 0 {
						state.TicketInfo = make([]PimEligibleRoleAssignmentTicketInfo, 1)
					}

					if ticketInfo.TicketNumber != nil {
						state.TicketInfo[0].TicketNumber = *ticketInfo.TicketNumber
					}
					if ticketInfo.TicketSystem != nil {
						state.TicketInfo[0].TicketSystem = *ticketInfo.TicketSystem
					}
				}

				if scheduleInfo := request.Properties.ScheduleInfo; scheduleInfo != nil {
					if len(state.ScheduleInfo) == 0 {
						state.ScheduleInfo = make([]PimEligibleRoleAssignmentScheduleInfo, 1)
					}

					// Only set the StartDateTime if not already present in state, because the value returned by the server advances
					// in short intervals until the request has been fully processed, causing unnecessary persistent diffs
					if state.ScheduleInfo[0].StartDateTime == "" && scheduleInfo.StartDateTime != nil {
						state.ScheduleInfo[0].StartDateTime = *scheduleInfo.StartDateTime
					}

					if expiration := scheduleInfo.Expiration; expiration != nil {
						if len(state.ScheduleInfo[0].Expiration) == 0 {
							state.ScheduleInfo[0].Expiration = make([]PimEligibleRoleAssignmentScheduleInfoExpiration, 1)
						}

						// Only set the EndDateTime if not already present in state, because the value returned by the server advances
						// in short intervals until the request has been fully processed, causing unnecessary persistent diffs
						if state.ScheduleInfo[0].Expiration[0].EndDateTime == "" && expiration.EndDateTime != nil {
							state.ScheduleInfo[0].Expiration[0].EndDateTime = pointer.From(expiration.EndDateTime)
						}

						if expiration.Duration != nil && *expiration.Duration != "" {
							durationRaw := *expiration.Duration

							reHours := regexp.MustCompile(`PT(\d+)H`)
							matches := reHours.FindStringSubmatch(durationRaw)
							if len(matches) == 2 {
								hours, err := strconv.ParseInt(matches[1], 10, 0)
								if err != nil {
									return fmt.Errorf("parsing duration: %+v", err)
								}
								state.ScheduleInfo[0].Expiration[0].DurationHours = hours
							}

							reDays := regexp.MustCompile(`P(\d+)D`)
							matches = reDays.FindStringSubmatch(durationRaw)
							if len(matches) == 2 {
								days, err := strconv.ParseInt(matches[1], 10, 0)
								if err != nil {
									return fmt.Errorf("parsing duration: %+v", err)
								}
								state.ScheduleInfo[0].Expiration[0].DurationDays = days
							}
						}
					}
				}
			} else if props := schedule.Properties; props != nil {
				// The request has likely expired, so populate from the schedule (not all fields will be available)
				state.PrincipalId = pointer.From(props.PrincipalId)
				state.PrincipalType = string(pointer.From(props.PrincipalType))
				state.RoleDefinitionId = pointer.From(props.RoleDefinitionId)

				if props.StartDateTime != nil {
					if len(state.ScheduleInfo) == 0 {
						state.ScheduleInfo = make([]PimEligibleRoleAssignmentScheduleInfo, 1)
					}

					// Only set the StartDateTime if not already present in state, because the value returned by the server advances
					// in short intervals until the request has been fully processed, causing unnecessary persistent diffs
					if state.ScheduleInfo[0].StartDateTime == "" {
						state.ScheduleInfo[0].StartDateTime = *props.StartDateTime
					}
				}

				if props.EndDateTime != nil {
					if len(state.ScheduleInfo) == 0 {
						state.ScheduleInfo = make([]PimEligibleRoleAssignmentScheduleInfo, 1)
					}
					if len(state.ScheduleInfo[0].Expiration) == 0 {
						state.ScheduleInfo[0].Expiration = make([]PimEligibleRoleAssignmentScheduleInfoExpiration, 1)
					}

					// Only set the EndDateTime if not already present in state, because the value returned by the server advances
					// in short intervals until the request has been fully processed, causing unnecessary persistent diffs
					if state.ScheduleInfo[0].Expiration[0].EndDateTime == "" {
						state.ScheduleInfo[0].Expiration[0].EndDateTime = *props.EndDateTime
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (PimEligibleRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			schedulesClient := metadata.Client.Authorization.RoleEligibilitySchedulesClient
			requestsClient := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state PimActiveRoleAssignmentModel
			if err = metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}

			schedule, err := findRoleEligibilitySchedule(ctx, schedulesClient, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if schedule == nil {
				log.Printf("[DEBUG] Eligible Role Assignment request has been canceled")
				return nil
			}
			if schedule.Properties == nil {
				return fmt.Errorf("retrieving %s: response with nil properties received", id)
			}

			switch pointer.From(schedule.Properties.Status) {
			case roleeligibilityschedules.StatusPendingApproval, roleeligibilityschedules.StatusPendingApprovalProvisioning,
				roleeligibilityschedules.StatusPendingEvaluation, roleeligibilityschedules.StatusGranted,
				roleeligibilityschedules.StatusPendingProvisioning, roleeligibilityschedules.StatusPendingAdminDecision:

				// Attempt to find a Request for this Schedule
				request, err := findRoleEligibilityScheduleRequest(ctx, requestsClient, schedule, id)
				if err != nil {
					return err
				}

				// Pending scheduled eligible role assignments should be removed by Cancel operation
				scheduleRequestId, err := roleeligibilityschedulerequests.ParseScopedRoleEligibilityScheduleRequestID(pointer.From(request.Id))
				if err != nil {
					return err
				}
				if _, err = requestsClient.Cancel(ctx, *scheduleRequestId); err != nil {
					return err
				}

			default:
				// Remove eligible role assignment by sending an AdminRemove request
				payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{
					Properties: &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{
						PrincipalId:      id.PrincipalId,
						RoleDefinitionId: id.RoleDefinitionId,
						RequestType:      roleeligibilityschedulerequests.RequestTypeAdminRemove,
						Justification:    pointer.To("Removed by Terraform"),
					},
				}

				// Include the ticket information from state for auditing purposes
				if len(state.TicketInfo) == 1 {
					payload.Properties.TicketInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo{}
					payload.Properties.TicketInfo.TicketNumber = &state.TicketInfo[0].TicketNumber
					payload.Properties.TicketInfo.TicketSystem = &state.TicketInfo[0].TicketSystem
				}

				roleEligibilityScheduleRequestName, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating uuid: %+v", err)
				}

				deleteId := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID(id.Scope, roleEligibilityScheduleRequestName)

				// Wait for removal request to be processed
				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{"Pending"},
					Target:  []string{"Submitted", "GoneAway"},
					Refresh: func() (interface{}, string, error) {
						// Removal request is not accepted within a minimum duration window, so retry it
						result, err := requestsClient.Create(ctx, deleteId, payload)
						if err != nil {
							if result.OData != nil && result.OData.Error != nil {
								if code := result.OData.Error.Code; code != nil {
									// API sometimes returns this error for a short while before relenting
									if *code == "ActiveDurationTooShort" {
										return result, "Pending", nil
									}

									// The principal is gone, so the eligible role assignment must also have gone away
									if *code == "RoleAssignmentDoesNotExist" {
										return result, "GoneAway", nil
									}
								}
							}

							return nil, "Error", fmt.Errorf("sending removal request for %s: %+v", id, err)
						}

						return result, "Submitted", nil
					},
					MinTimeout: 1 * time.Minute,
					Timeout:    time.Until(deadline),
				}

				if _, err = stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for removal request %s to be processed: %+v", id, err)
				}
			}

			// Wait for role eligibility schedule to disappear
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Exists"},
				Target:     []string{"NotFound"},
				Refresh:    pollForRoleEligibilitySchedule(ctx, schedulesClient, *id),
				MinTimeout: 10 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be removed: %+v", id, err)
			}

			return nil
		},
	}
}

func pollForRoleEligibilitySchedule(ctx context.Context, client *roleeligibilityschedules.RoleEligibilitySchedulesClient, id parse.PimRoleAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Polling for %s", id)

		schedule, err := findRoleEligibilitySchedule(ctx, client, id)
		if err != nil {
			return schedule, "Error", err
		}

		if schedule == nil {
			return schedule, "NotFound", nil
		}

		return schedule, "Exists", nil
	}
}

func findRoleEligibilitySchedule(ctx context.Context, client *roleeligibilityschedules.RoleEligibilitySchedulesClient, id parse.PimRoleAssignmentId) (*roleeligibilityschedules.RoleEligibilitySchedule, error) {
	scopeId, err := commonids.ParseScopeID(id.Scope)
	if err != nil {
		return nil, err
	}

	schedulesResult, err := client.ListForScopeComplete(ctx, *scopeId, roleeligibilityschedules.ListForScopeOperationOptions{
		Filter: pointer.To(fmt.Sprintf("(principalId eq '%s')", id.PrincipalId)),
	})
	if err != nil {
		return nil, fmt.Errorf("listing Role Eligiblity Schedules for %s: %+v", scopeId, err)
	}

	for _, schedule := range schedulesResult.Items {
		if props := schedule.Properties; props != nil {
			if props.RoleDefinitionId != nil && strings.EqualFold(*props.RoleDefinitionId, id.RoleDefinitionId) &&
				props.Scope != nil && strings.EqualFold(*props.Scope, scopeId.ID()) &&
				props.PrincipalId != nil && strings.EqualFold(*props.PrincipalId, id.PrincipalId) &&
				props.MemberType != nil && *props.MemberType == roleeligibilityschedules.MemberTypeDirect {
				return &schedule, nil
			}
		}
	}

	return nil, nil
}

func findRoleEligibilityScheduleRequest(ctx context.Context, client *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient, schedule *roleeligibilityschedules.RoleEligibilitySchedule, id *parse.PimRoleAssignmentId) (*roleeligibilityschedulerequests.RoleEligibilityScheduleRequest, error) {
	// Request ID was provided, so retrieve it individually
	if schedule.Properties.RoleEligibilityScheduleRequestId != nil {
		requestId, err := roleeligibilityschedulerequests.ParseScopedRoleEligibilityScheduleRequestID(*schedule.Properties.RoleEligibilityScheduleRequestId)
		if err != nil { //
			return nil, err
		}

		requestResp, err := client.Get(ctx, *requestId)
		if err != nil && !response.WasNotFound(requestResp.HttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", requestId, err)
		}

		if !response.WasNotFound(requestResp.HttpResponse) {
			return requestResp.Model, nil
		}
	}

	// Request ID not provided or was invalid, list by scope and filter by principal for a best-effort search
	if principalId := schedule.Properties.PrincipalId; principalId != nil && id != nil {
		scopeId, err := commonids.ParseScopeID(id.Scope)
		if err != nil {
			return nil, err
		}

		requestsResult, err := client.ListForScopeComplete(ctx, *scopeId, roleeligibilityschedulerequests.ListForScopeOperationOptions{
			Filter: pointer.To(fmt.Sprintf("principalId eq '%s'", *principalId)),
		})
		if err != nil {
			return nil, fmt.Errorf("listing Role Eligibility Requests for principal_id %q: %+v", *principalId, err)
		}
		for _, item := range requestsResult.Items {
			if props := item.Properties; props != nil {
				if props.TargetRoleEligibilityScheduleId != nil && strings.EqualFold(*props.TargetRoleEligibilityScheduleId, id.ID()) &&
					props.RequestType == roleeligibilityschedulerequests.RequestTypeAdminAssign && strings.EqualFold(props.PrincipalId, *principalId) {
					return pointer.To(item), nil
				}
			}
		}
	}

	// No request was found, it probably expired
	return nil, nil
}
