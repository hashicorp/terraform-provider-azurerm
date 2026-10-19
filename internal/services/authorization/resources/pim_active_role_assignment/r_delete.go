// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pim_active_role_assignment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedulerequests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedules"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			schedulesClient := metadata.Client.Authorization.RoleAssignmentSchedulesClient
			requestsClient := metadata.Client.Authorization.RoleAssignmentScheduleRequestClient

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

			schedule, err := findRoleAssignmentSchedule(ctx, schedulesClient, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if schedule == nil {
				log.Printf("[DEBUG] Role Assignment request has been canceled")
				return nil
			}
			if schedule.Properties == nil {
				return fmt.Errorf("retrieving %s: response with nil properties received", id)
			}

			switch pointer.From(schedule.Properties.Status) {
			case roleassignmentschedules.StatusPendingApproval, roleassignmentschedules.StatusPendingApprovalProvisioning,
				roleassignmentschedules.StatusPendingEvaluation, roleassignmentschedules.StatusGranted,
				roleassignmentschedules.StatusPendingProvisioning, roleassignmentschedules.StatusPendingAdminDecision:

				// Attempt to find a Request for this Schedule
				request, err := findRoleAssignmentScheduleRequest(ctx, requestsClient, schedule, id)
				if err != nil {
					return err
				}

				// Pending scheduled role assignments should be removed by Cancel operation
				scheduleRequestId, err := roleassignmentschedulerequests.ParseScopedRoleAssignmentScheduleRequestID(pointer.From(request.Id))
				if err != nil {
					return err
				}
				if _, err = requestsClient.Cancel(ctx, *scheduleRequestId); err != nil {
					return err
				}

				return nil
			default:
				// Remove active role assignment by sending an AdminRemove request
				payload := roleassignmentschedulerequests.RoleAssignmentScheduleRequest{
					Properties: &roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties{
						PrincipalId:      id.PrincipalId,
						RoleDefinitionId: id.RoleDefinitionId,
						RequestType:      roleassignmentschedulerequests.RequestTypeAdminRemove,
						Justification:    pointer.To("Removed by Terraform"),
					},
				}

				// Include the ticket information from state for auditing purposes
				if len(state.TicketInfo) == 1 {
					payload.Properties.TicketInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo{}
					payload.Properties.TicketInfo.TicketNumber = &state.TicketInfo[0].TicketNumber
					payload.Properties.TicketInfo.TicketSystem = &state.TicketInfo[0].TicketSystem
				}

				roleAssignmentScheduleRequestName, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating uuid: %+v", err)
				}

				deleteId := roleassignmentschedulerequests.NewScopedRoleAssignmentScheduleRequestID(id.Scope, roleAssignmentScheduleRequestName)

				// Wait for removal request to be processed
				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{"Pending"},
					Target:  []string{"Submitted", "GoneAway"},
					Refresh: func() (any, string, error) {
						// Removal request is not accepted within a minimum duration window, so retry it
						result, err := requestsClient.Create(ctx, deleteId, payload)
						if err != nil {
							if result.OData != nil && result.OData.Error != nil {
								if code := result.OData.Error.Code; code != nil {
									// API sometimes returns this error for a short while before relenting
									if *code == "ActiveDurationTooShort" {
										return result, "Pending", nil
									}

									// The principal is gone, so the role assignment must also have gone away
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

			// Wait for role assignment schedule to disappear
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Exists"},
				Target:     []string{"NotFound"},
				Refresh:    pollForRoleAssignmentSchedule(ctx, schedulesClient, *id),
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
