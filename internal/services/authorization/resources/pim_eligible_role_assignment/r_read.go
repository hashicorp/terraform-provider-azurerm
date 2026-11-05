// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pim_eligible_role_assignment

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
				state.PrincipalType = pointer.FromEnum(request.Properties.PrincipalType)
				state.RoleDefinitionId = request.Properties.RoleDefinitionId

				state.Condition = pointer.From(request.Properties.Condition)
				state.ConditionVersion = pointer.From(request.Properties.ConditionVersion)

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
				state.PrincipalType = pointer.FromEnum(props.PrincipalType)
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

func pollForRoleEligibilitySchedule(ctx context.Context, client *roleeligibilityschedules.RoleEligibilitySchedulesClient, id parse.PimRoleAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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
		Filter: pointer.To(fmt.Sprintf("(principalId eq '%s') and atScope()", id.PrincipalId)),
	})
	if err != nil {
		return nil, fmt.Errorf("listing Role Eligibility Schedules for %s: %+v", scopeId, err)
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
			Filter: pointer.To(fmt.Sprintf("principalId eq '%s' and atScope()", *principalId)),
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
