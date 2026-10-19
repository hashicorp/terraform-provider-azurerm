// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pim_eligible_role_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				schedule, err := findRoleEligibilitySchedule(ctx, schedulesClient, id)
				if err != nil {
					return err
				}
				if schedule != nil {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			properties := &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{
				Justification:    pointer.To(config.Justification),
				PrincipalId:      id.PrincipalId,
				RequestType:      roleeligibilityschedulerequests.RequestTypeAdminAssign,
				RoleDefinitionId: id.RoleDefinitionId,
				Scope:            pointer.To(scopeId.ID()),
				ScheduleInfo:     scheduleInfo,
				TicketInfo:       ticketInfo,
			}

			condition := config.Condition
			conditionVersion := config.ConditionVersion

			if condition != "" && conditionVersion != "" {
				properties.Condition = pointer.To(condition)
				properties.ConditionVersion = pointer.To(conditionVersion)
			}

			payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{
				Properties: properties,
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
				Refresh: func() (any, string, error) {
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
