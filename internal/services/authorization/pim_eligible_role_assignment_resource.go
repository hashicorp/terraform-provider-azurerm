package authorization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	// nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = PimEligibleRoleAssignmentResource{}

type PimEligibleRoleAssignmentResource struct{}

func (PimEligibleRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The scope.",
		},

		"role_definition_id": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The role definition id.",
		},

		"principal_id": { // not sure how to validate guids or if possible. service principals will give a poor error message back.
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The principal id.",
		},

		"ticket": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			ForceNew:    true,
			Description: "Ticket details relating to the assignment.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"number": {
						Optional:    true,
						Type:        pluginsdk.TypeString,
						Description: "The ticket number.",
					},
					"system": {
						Optional:    true,
						Type:        pluginsdk.TypeString,
						Description: "The ticket system.",
					},
				},
			},
		},

		"schedule": {
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			ForceNew:    true,
			Description: "The schedule details of this eligible role assignment.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date_time": { // defaults to now
						Optional:    true,
						Computed:    true,
						ForceNew:    true,
						Type:        pluginsdk.TypeString,
						Description: "The start date time.",
					},
					"expiration": { // if none specified, it's a permanent assignment
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"duration_days": {
									Optional: true,
									Computed: true,
									ForceNew: true,
									Type:     pluginsdk.TypeInt,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_hours",
										"schedule.0.expiration.0.end_date_time",
									},
									Description: "The duration of the assignment in days.",
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
									Description: "The duration of the assignment in hours.",
								},
								"end_date_time": {
									Optional: true,
									Computed: true,
									ForceNew: true,
									Type:     pluginsdk.TypeString,
									ConflictsWith: []string{
										"schedule.0.expiration.0.duration_days",
										"schedule.0.expiration.0.duration_hours",
									},
									Description: "The end date time of the assignment.",
								},
							},
						},
					},
				},
			},
		},

		"justification": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The justification of the eligible role assignment.",
		},
	}
}

func (PimEligibleRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:        pluginsdk.TypeString,
			Description: "The type of principal.",
			Computed:    true,
		},
	}
}

func (PimEligibleRoleAssignmentResource) ModelObject() interface{} {
	return &PimEligibleRoleAssignmentResourceSchema{}
}

func (PimEligibleRoleAssignmentResource) ResourceType() string {
	return "azurerm_pim_eligible_role_assignment"
}

func (r PimEligibleRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clientInstances := metadata.Client.Authorization.RoleEligibilityScheduleInstancesClient
			clientRequest := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			scope := metadata.ResourceData.Get("scope").(string)
			roleDefinitionId := metadata.ResourceData.Get("role_definition_id").(string)
			principalId := metadata.ResourceData.Get("principal_id").(string)

			id := parse.NewPimRoleAssignmentID(scope, roleDefinitionId, principalId)

			filter := &roleeligibilityscheduleinstances.ListForScopeOperationOptions{
				Filter: pointer.To(fmt.Sprintf("(principalId eq '%s' and roleDefinitionId eq '%s')", id.PrincipalId, id.RoleDefinitionId)),
			}

			items, err := clientInstances.ListForScopeComplete(ctx, id.ScopeID(), *filter)
			if err != nil {
				return fmt.Errorf("listing role assignments on scope %s: %+v", id, err)
			}
			for _, item := range items.Items {
				if *item.Properties.MemberType == roleeligibilityscheduleinstances.MemberTypeDirect {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var config PimEligibleRoleAssignmentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload roleeligibilityschedulerequests.RoleEligibilityScheduleRequest

			r.mapPimEligibleRoleAssignmentResourceSchemaToRoleEligibilityScheduleRequest(config, &payload)

			payload.Properties.RequestType = roleeligibilityschedulerequests.RequestTypeAdminAssign

			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating uuid: %+v", err)
			}

			requestId := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID(config.Scope, uuid)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Missing"},
				Target:     []string{"Created"},
				Refresh:    createEligibilityRoleAssignment(ctx, clientRequest, requestId, &payload),
				MinTimeout: 30 * time.Second,
				Timeout:    time.Until(deadline),
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be created: %+v", id, err)
			}

			// wait for resource to exist
			stateConf = &pluginsdk.StateChangeConf{
				Pending:    []string{"Missing"},
				Target:     []string{"Found"},
				Refresh:    waitForEligibleRoleAssignmentSchedule(ctx, clientInstances, config.Scope, config.PrincipalId, config.RoleDefinitionId, "Found"),
				MinTimeout: 30 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
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
			clientInstances := metadata.Client.Authorization.RoleEligibilityScheduleInstancesClient
			clientRequest := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			schema := PimEligibleRoleAssignmentResourceSchema{}

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			filter := &roleeligibilityscheduleinstances.ListForScopeOperationOptions{
				Filter: pointer.To(fmt.Sprintf("(principalId eq '%s' and roleDefinitionId eq '%s')", id.PrincipalId, id.RoleDefinitionId)),
			}

			items, err := clientInstances.ListForScopeComplete(ctx, id.ScopeID(), *filter)
			if err != nil {
				return fmt.Errorf("listing role assignments on scope %s: %+v", id, err)
			}
			var instance *roleeligibilityscheduleinstances.RoleEligibilityScheduleInstance
			for _, item := range items.Items {
				if *item.Properties.MemberType == roleeligibilityscheduleinstances.MemberTypeDirect {
					instance = &item
				}
			}
			if instance == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Scope = id.Scope

			guid, err := parse.RoleEligibilityScheduleIdFromInstance(instance)
			if err != nil {
				return err
			}
			scheduleRequestId := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID(id.Scope, *guid)

			resp, err := clientRequest.Get(ctx, scheduleRequestId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Scope = id.Scope

				if err := r.mapRoleAssignmentScheduleRequestToPimEligibleRoleAssignmentResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			if val, ok := metadata.ResourceData.GetOk("schedule.0.start_date_time"); ok && len(schema.ScheduleInfo) > 0 {
				schema.ScheduleInfo[0].StartDateTime = val.(string)
			}
			if val, ok := metadata.ResourceData.GetOk("schedule.0.expiration.0.end_date_time"); ok && len(schema.ScheduleInfo) > 0 && len(schema.ScheduleInfo[0].Expiration) > 0 {
				schema.ScheduleInfo[0].Expiration[0].EndDateTime = val.(string)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (PimEligibleRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clientInstances := metadata.Client.Authorization.RoleEligibilityScheduleInstancesClient
			clientRequest := metadata.Client.Authorization.RoleEligibilityScheduleRequestClient

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config PimEligibleRoleAssignmentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{}
			payload.Properties = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{}
			payload.Properties.PrincipalId = id.PrincipalId
			payload.Properties.RoleDefinitionId = id.RoleDefinitionId
			payload.Properties.RequestType = roleeligibilityschedulerequests.RequestTypeAdminRemove
			payload.Properties.ScheduleInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo{}

			if config.Justification != "" {
				payload.Properties.Justification = &config.Justification
			}
			if len(config.TicketInfo) == 1 {
				payload.Properties.TicketInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo{}
				payload.Properties.TicketInfo.TicketNumber = &config.TicketInfo[0].TicketNumber
				payload.Properties.TicketInfo.TicketSystem = &config.TicketInfo[0].TicketSystem
			}

			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating uuid: %+v", err)
			}
			deleteId := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID(id.Scope, uuid)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}

			// wait for resource to deleted
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Exist"},
				Target:     []string{"Deleted"},
				Refresh:    deleteEligibilityRoleAssignmentSchedule(ctx, clientRequest, deleteId, &payload),
				MinTimeout: 1 * time.Minute,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become deleted: %+v", id, err)
			}

			// wait for role assignment to be missing
			stateConf = &pluginsdk.StateChangeConf{
				Pending:    []string{"Found"},
				Target:     []string{"Missing"},
				Refresh:    waitForEligibleRoleAssignmentSchedule(ctx, clientInstances, id.Scope, id.PrincipalId, id.RoleDefinitionId, "Missing"),
				MinTimeout: 30 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become missing: %+v", id, err)
			}

			return nil
		},
	}
}

func (PimEligibleRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidatePimRoleAssignmentID
}

type PimEligibleRoleAssignmentResourceSchema struct {
	RoleDefinitionId string                                                `tfschema:"role_definition_id"`
	Scope            string                                                `tfschema:"scope"`
	PrincipalId      string                                                `tfschema:"principal_id"`
	PrincipalType    string                                                `tfschema:"principal_type"`
	Justification    string                                                `tfschema:"justification"`
	TicketInfo       []PimEligibleRoleAssignmentResourceSchemaTicketInfo   `tfschema:"ticket"`
	ScheduleInfo     []PimEligibleRoleAssignmentResourceSchemaScheduleInfo `tfschema:"schedule"`
}

type PimEligibleRoleAssignmentResourceSchemaTicketInfo struct {
	TicketNumber string `tfschema:"number"`
	TicketSystem string `tfschema:"system"`
}

type PimEligibleRoleAssignmentResourceSchemaScheduleInfo struct {
	StartDateTime string                                                          `tfschema:"start_date_time"`
	Expiration    []PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration `tfschema:"expiration"`
}

type PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration struct {
	DurationDays  int    `tfschema:"duration_days"`
	DurationHours int    `tfschema:"duration_hours"`
	EndDateTime   string `tfschema:"end_date_time"`
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaToRoleEligibilityScheduleRequest(input PimEligibleRoleAssignmentResourceSchema, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequest) {
	if output.Properties == nil {
		output.Properties = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{}
	}

	r.mapPimEligibleRoleAssignmentResourceSchemaToRoleEligibilityScheduleRequestProperties(input, output.Properties)
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaToRoleEligibilityScheduleRequestProperties(input PimEligibleRoleAssignmentResourceSchema, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties) {
	output.Justification = &input.Justification
	output.PrincipalId = input.PrincipalId

	output.RoleDefinitionId = input.RoleDefinitionId
	output.Scope = &input.Scope

	if len(input.TicketInfo) > 0 {
		r.mapPimEligibleRoleAssignmentResourceSchemaTicketInfoToRoleEligibilityScheduleRequestProperties(input.TicketInfo[0], output)
	}

	if len(input.ScheduleInfo) > 0 {
		r.mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoToRoleEligibilityScheduleRequestProperties(input.ScheduleInfo[0], output)
	} else {
		output.ScheduleInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo{}
		output.ScheduleInfo.Expiration = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration{}
		output.ScheduleInfo.Expiration.Type = pointer.To(roleeligibilityschedulerequests.TypeNoExpiration)
	}
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaTicketInfoToRoleEligibilityScheduleRequestProperties(input PimEligibleRoleAssignmentResourceSchemaTicketInfo, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties) {
	if output.TicketInfo == nil {
		output.TicketInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo{}
	}
	r.mapPimEligibleRoleAssignmentResourceSchemaTicketInfoToRoleEligibilityScheduleRequestPropertiesTicketInfo(input, output.TicketInfo)
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaTicketInfoToRoleEligibilityScheduleRequestPropertiesTicketInfo(input PimEligibleRoleAssignmentResourceSchemaTicketInfo, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo) {
	output.TicketNumber = pointer.To(input.TicketNumber)
	output.TicketSystem = pointer.To(input.TicketSystem)
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoToRoleEligibilityScheduleRequestProperties(input PimEligibleRoleAssignmentResourceSchemaScheduleInfo, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties) {
	if output.ScheduleInfo == nil {
		output.ScheduleInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo{}
	}

	r.mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoToRoleEligibilityScheduleRequestPropertiesScheduleInfo(input, output.ScheduleInfo)
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoToRoleEligibilityScheduleRequestPropertiesScheduleInfo(input PimEligibleRoleAssignmentResourceSchemaScheduleInfo, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo) {
	output.StartDateTime = pointer.To(input.StartDateTime)

	if output.Expiration == nil {
		output.Expiration = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration{}
	}
	if len(input.Expiration) > 0 {
		r.mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoExpirationToRoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration(input.Expiration[0], output.Expiration)
	} else {
		// when no expiration is specified, set the type to No Expiration
		output.Expiration.Type = pointer.To(roleeligibilityschedulerequests.TypeNoExpiration)
	}
}

func (r PimEligibleRoleAssignmentResource) mapPimEligibleRoleAssignmentResourceSchemaScheduleInfoExpirationToRoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration(input PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration, output *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration) {
	switch {
	case input.DurationDays != 0:
		output.Duration = pointer.To(fmt.Sprintf("P%dD", input.DurationDays))

	case input.DurationHours != 0:
		output.Duration = pointer.To(fmt.Sprintf("PT%dH", input.DurationHours))

	case input.EndDateTime != "":
		output.EndDateTime = pointer.To(input.EndDateTime)

	}

	// value of duration and end date determine expiration type
	switch {
	case output.Duration != nil && *output.Duration != "":
		output.Type = pointer.To(roleeligibilityschedulerequests.TypeAfterDuration)

	case input.EndDateTime != "":
		output.Type = pointer.To(roleeligibilityschedulerequests.TypeAfterDateTime)

	default:
		output.Type = pointer.To(roleeligibilityschedulerequests.TypeNoExpiration)
	}
}

func (r PimEligibleRoleAssignmentResource) mapRoleAssignmentScheduleRequestToPimEligibleRoleAssignmentResourceSchema(input roleeligibilityschedulerequests.RoleEligibilityScheduleRequest, output *PimEligibleRoleAssignmentResourceSchema) error {
	if input.Properties == nil {
		input.Properties = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties{}
	}

	if err := r.mapRoleEligibilityScheduleRequestPropertiesToPimEligibleRoleAssignmentResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "RoleEligibilityScheduleRequestProperties", "Properties", err)
	}

	return nil
}

func (r PimEligibleRoleAssignmentResource) mapRoleEligibilityScheduleRequestPropertiesToPimEligibleRoleAssignmentResourceSchema(input roleeligibilityschedulerequests.RoleEligibilityScheduleRequestProperties, output *PimEligibleRoleAssignmentResourceSchema) error {
	output.Justification = pointer.From(input.Justification)
	output.PrincipalId = input.PrincipalId
	output.PrincipalType = string(*input.PrincipalType)
	output.RoleDefinitionId = input.RoleDefinitionId

	if input.TicketInfo == nil {
		input.TicketInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo{}
	}

	if input.ScheduleInfo == nil {
		input.ScheduleInfo = &roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo{}
	}

	if input.TicketInfo != nil && (input.TicketInfo.TicketNumber != nil ||
		input.TicketInfo.TicketSystem != nil) {
		tempTicketInfo := &PimEligibleRoleAssignmentResourceSchemaTicketInfo{}
		if err := r.mapRoleEligibilityScheduleRequestPropertiesTicketInfoToPimEligibleRoleAssignmentResourceSchemaTicketInfo(*input.TicketInfo, tempTicketInfo); err != nil {
			return err
		} else {
			output.TicketInfo = make([]PimEligibleRoleAssignmentResourceSchemaTicketInfo, 0)
			output.TicketInfo = append(output.TicketInfo, *tempTicketInfo)
		}
	}

	tempScheduleInfo := &PimEligibleRoleAssignmentResourceSchemaScheduleInfo{}
	if err := r.mapRoleEligibilityScheduleRequestPropertiesScheduleInfoToPimEligibleRoleAssignmentResourceSchemaScheduleInfo(*input.ScheduleInfo, tempScheduleInfo); err != nil {
		return err
	} else {
		output.ScheduleInfo = make([]PimEligibleRoleAssignmentResourceSchemaScheduleInfo, 0)
		output.ScheduleInfo = append(output.ScheduleInfo, *tempScheduleInfo)
	}

	return nil
}

func (r PimEligibleRoleAssignmentResource) mapRoleEligibilityScheduleRequestPropertiesTicketInfoToPimEligibleRoleAssignmentResourceSchemaTicketInfo(input roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesTicketInfo, output *PimEligibleRoleAssignmentResourceSchemaTicketInfo) error {
	output.TicketNumber = pointer.From(input.TicketNumber)
	output.TicketSystem = pointer.From(input.TicketSystem)
	return nil
}

func (r PimEligibleRoleAssignmentResource) mapRoleEligibilityScheduleRequestPropertiesScheduleInfoToPimEligibleRoleAssignmentResourceSchemaScheduleInfo(input roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfo, output *PimEligibleRoleAssignmentResourceSchemaScheduleInfo) error {
	output.StartDateTime = pointer.From(input.StartDateTime)

	tempExpiration := &PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration{}
	if err := r.mapRoleEligibilityScheduleRequestPropertiesExpirationToPimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration(*input.Expiration, tempExpiration); err != nil {
		return err
	} else {
		output.Expiration = make([]PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration, 0)
		output.Expiration = append(output.Expiration, *tempExpiration)
	}
	return nil
}

func (r PimEligibleRoleAssignmentResource) mapRoleEligibilityScheduleRequestPropertiesExpirationToPimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration(input roleeligibilityschedulerequests.RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration, output *PimEligibleRoleAssignmentResourceSchemaScheduleInfoExpiration) error {
	if input.Duration != nil && *input.Duration != "" {
		durationRaw := *input.Duration

		reHours := regexp.MustCompile(`PT(\d+)H`)
		matches := reHours.FindStringSubmatch(durationRaw)
		if len(matches) == 2 {
			hours, err := strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("could not decode hours from RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration: %+v", err)
			}
			output.DurationHours = hours
		}
		reDays := regexp.MustCompile(`P(\d+)D`)
		matches = reDays.FindStringSubmatch(durationRaw)
		if len(matches) == 2 {
			days, err := strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("could not decode days from RoleEligibilityScheduleRequestPropertiesScheduleInfoExpiration: %+v", err)
			}
			output.DurationDays = days
		}
	}

	output.EndDateTime = pointer.From(input.EndDateTime)
	return nil
}

func createEligibilityRoleAssignment(ctx context.Context, client *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient, id roleeligibilityschedulerequests.ScopedRoleEligibilityScheduleRequestId, payload *roleeligibilityschedulerequests.RoleEligibilityScheduleRequest) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {

		// Azure can error when the subject doesn't exist yet due to AAD replication
		// Retry deletes while that error exists.
		result, err := client.Create(ctx, id, *payload)
		if err != nil {
			if *result.OData.Error.Code == "SubjectNotFound" {
				return nil, "Missing", nil
			}

			return nil, "Missing", fmt.Errorf("creating %s: %+v", id, err)
		}

		return result, "Created", nil
	}
}

func waitForEligibleRoleAssignmentSchedule(ctx context.Context, client *roleeligibilityscheduleinstances.RoleEligibilityScheduleInstancesClient, scope string, principalId string, roleDefinitionId string, target string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if Role Assignment is %s on %q with role %q for %q.", target, scope, roleDefinitionId, principalId)

		instanceId := commonids.NewScopeID(scope)
		filter := &roleeligibilityscheduleinstances.ListForScopeOperationOptions{
			Filter: pointer.To(fmt.Sprintf("assignedTo('%s')", principalId)),
		}

		items, err := client.ListForScopeComplete(ctx, instanceId, *filter)
		if err != nil {
			return nil, "", fmt.Errorf("listing role assignments on scope %s: %+v", instanceId, err)
		}
		state := "Missing"
		var result interface{}

		for _, item := range items.Items {
			if *item.Properties.RoleDefinitionId == roleDefinitionId &&
				*item.Properties.MemberType == roleeligibilityscheduleinstances.MemberTypeDirect {
				state = "Found"
				result = item
			}
		}

		if target == "Missing" && state == "Missing" {
			result = &roleeligibilityscheduleinstances.RoleEligibilityScheduleInstance{}
		}

		return result, state, nil
	}
}

func deleteEligibilityRoleAssignmentSchedule(ctx context.Context, client *roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient, id roleeligibilityschedulerequests.ScopedRoleEligibilityScheduleRequestId, payload *roleeligibilityschedulerequests.RoleEligibilityScheduleRequest) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {

		// Azure can error when the role hasn't existed for less than 5 minutes.
		// Retry deletes while that error exists.
		result, err := client.Create(ctx, id, *payload)
		if err != nil {
			if *result.OData.Error.Code == "ActiveDurationTooShort" {
				return nil, "Exist", nil
			}

			if *result.OData.Error.Code == "RoleAssignmentDoesNotExist" {
				return nil, "Deleted", nil
			}

			return nil, "Exist", fmt.Errorf("creating %s: %+v", id, err)
		}

		return result, "Deleted", nil
	}
}
