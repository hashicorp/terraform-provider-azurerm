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
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentschedulerequests"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = PimActiveRoleAssignmentResource{}

type PimActiveRoleAssignmentResource struct{}

func (PimActiveRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
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

		"principal_id": {
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
			Description: "The ticket details.",
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
			Description: "The schedule details of this role assignment.",
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
			Description: "The justification of the role assignment.",
		},
	}
}

func (PimActiveRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The type of principal.",
		},
	}
}

func (PimActiveRoleAssignmentResource) ModelObject() interface{} {
	return &PimActiveRoleAssignmentResourceSchema{}
}

func (PimActiveRoleAssignmentResource) ResourceType() string {
	return "azurerm_pim_active_role_assignment"
}

func (r PimActiveRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clientInstances := metadata.Client.Authorization.RoleAssignmentScheduleInstancesClient
			clientRequest := metadata.Client.Authorization.RoleAssignmentScheduleRequestClient

			scope := metadata.ResourceData.Get("scope").(string)
			roleDefinitionId := metadata.ResourceData.Get("role_definition_id").(string)
			principalId := metadata.ResourceData.Get("principal_id").(string)

			id := parse.NewPimRoleAssignmentID(scope, roleDefinitionId, principalId)

			filter := &roleassignmentscheduleinstances.ListForScopeOperationOptions{
				Filter: pointer.To(fmt.Sprintf("(principalId eq '%s' and roleDefinitionId eq '%s')", id.PrincipalId, id.RoleDefinitionId)),
			}

			items, err := clientInstances.ListForScopeComplete(ctx, id.ScopeID(), *filter)
			if err != nil {
				return fmt.Errorf("listing role assignments on scope %s: %+v", id, err)
			}
			for _, item := range items.Items {
				if *item.Properties.MemberType == roleassignmentscheduleinstances.MemberTypeDirect {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var config PimActiveRoleAssignmentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload roleassignmentschedulerequests.RoleAssignmentScheduleRequest

			r.mapPimActiveRoleAssignmentResourceSchemaToRoleAssignmentScheduleRequest(config, &payload)

			payload.Properties.RequestType = roleassignmentschedulerequests.RequestTypeAdminAssign

			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating uuid: %+v", err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}

			requestId := roleassignmentschedulerequests.NewScopedRoleAssignmentScheduleRequestID(config.Scope, uuid)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Missing"},
				Target:     []string{"Created"},
				Refresh:    createActiveRoleAssignment(ctx, clientRequest, requestId, &payload),
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
				Refresh:    waitForActiveRoleAssignment(ctx, clientInstances, config.Scope, config.PrincipalId, config.RoleDefinitionId, "Found"),
				MinTimeout: 30 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become found: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PimActiveRoleAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clientInstances := metadata.Client.Authorization.RoleAssignmentScheduleInstancesClient
			clientRequest := metadata.Client.Authorization.RoleAssignmentScheduleRequestClient

			schema := PimActiveRoleAssignmentResourceSchema{}

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			filter := &roleassignmentscheduleinstances.ListForScopeOperationOptions{
				Filter: pointer.To(fmt.Sprintf("(principalId eq '%s' and roleDefinitionId eq '%s')", id.PrincipalId, id.RoleDefinitionId)),
			}

			items, err := clientInstances.ListForScopeComplete(ctx, id.ScopeID(), *filter)
			if err != nil {
				return fmt.Errorf("listing role assignments on scope %s: %+v", id, err)
			}
			var instance *roleassignmentscheduleinstances.RoleAssignmentScheduleInstance
			for _, item := range items.Items {
				if *item.Properties.MemberType == roleassignmentscheduleinstances.MemberTypeDirect {
					instance = &item
				}
			}
			if instance == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Scope = id.Scope

			guid, err := parse.RoleAssignmentScheduleIdFromInstance(instance)
			if err != nil {
				return err
			}
			scheduleRequestId := roleassignmentschedulerequests.NewScopedRoleAssignmentScheduleRequestID(id.Scope, *guid)

			resp, err := clientRequest.Get(ctx, scheduleRequestId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Scope = id.Scope

				if err := r.mapRoleAssignmentScheduleRequestToPimActiveRoleAssignmentResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			// The API returns the start date time of when the request has been processed, instead of the date/time of when the request was submitted.
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

func (PimActiveRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clientRequest := metadata.Client.Authorization.RoleAssignmentScheduleRequestClient
			clientInstances := metadata.Client.Authorization.RoleAssignmentScheduleInstancesClient

			id, err := parse.PimRoleAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config PimActiveRoleAssignmentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			payload := roleassignmentschedulerequests.RoleAssignmentScheduleRequest{}
			payload.Properties = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties{}
			payload.Properties.PrincipalId = id.PrincipalId
			payload.Properties.RoleDefinitionId = id.RoleDefinitionId
			payload.Properties.RequestType = roleassignmentschedulerequests.RequestTypeAdminRemove
			payload.Properties.ScheduleInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo{}

			if config.Justification != "" {
				payload.Properties.Justification = &config.Justification
			}
			if len(config.TicketInfo) == 1 {
				payload.Properties.TicketInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo{}
				payload.Properties.TicketInfo.TicketNumber = &config.TicketInfo[0].TicketNumber
				payload.Properties.TicketInfo.TicketSystem = &config.TicketInfo[0].TicketSystem
			}

			uuid, err := uuid.GenerateUUID()
			if err != nil {
				return fmt.Errorf("generating uuid: %+v", err)
			}
			deleteId := roleassignmentschedulerequests.NewScopedRoleAssignmentScheduleRequestID(id.Scope, uuid)

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal error: context has no deadline")
			}
			// wait for resource to deleted
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Exist"},
				Target:     []string{"Deleted"},
				Refresh:    deleteActiveRoleAssignment(ctx, clientRequest, deleteId, &payload),
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
				Refresh:    waitForActiveRoleAssignment(ctx, clientInstances, id.Scope, id.PrincipalId, id.RoleDefinitionId, "Missing"),
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

func (PimActiveRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ValidatePimRoleAssignmentID
}

type PimActiveRoleAssignmentResourceSchema struct {
	RoleDefinitionId string                                              `tfschema:"role_definition_id"`
	Scope            string                                              `tfschema:"scope"`
	PrincipalId      string                                              `tfschema:"principal_id"`
	PrincipalType    string                                              `tfschema:"principal_type"`
	Justification    string                                              `tfschema:"justification"`
	TicketInfo       []PimActiveRoleAssignmentResourceSchemaTicketInfo   `tfschema:"ticket"`
	ScheduleInfo     []PimActiveRoleAssignmentResourceSchemaScheduleInfo `tfschema:"schedule"`
}

type PimActiveRoleAssignmentResourceSchemaTicketInfo struct {
	TicketNumber string `tfschema:"number"`
	TicketSystem string `tfschema:"system"`
}

type PimActiveRoleAssignmentResourceSchemaScheduleInfo struct {
	StartDateTime string                                                        `tfschema:"start_date_time"`
	Expiration    []PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration `tfschema:"expiration"`
}

type PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration struct {
	DurationDays  int    `tfschema:"duration_days"`
	DurationHours int    `tfschema:"duration_hours"`
	EndDateTime   string `tfschema:"end_date_time"`
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaToRoleAssignmentScheduleRequest(input PimActiveRoleAssignmentResourceSchema, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequest) {
	if output.Properties == nil {
		output.Properties = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties{}
	}

	r.mapPimActiveRoleAssignmentResourceSchemaToRoleAssignmentScheduleRequestProperties(input, output.Properties)
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaToRoleAssignmentScheduleRequestProperties(input PimActiveRoleAssignmentResourceSchema, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties) {
	output.Justification = &input.Justification
	output.PrincipalId = input.PrincipalId

	output.RoleDefinitionId = input.RoleDefinitionId
	output.Scope = &input.Scope

	if len(input.TicketInfo) > 0 {
		r.mapPimActiveRoleAssignmentResourceSchemaTicketInfoToRoleAssignmentScheduleRequestProperties(input.TicketInfo[0], output)
	}

	if len(input.ScheduleInfo) > 0 {
		r.mapRoleAssignmentScheduleRequestResourceScheduleInfoSchemaToRoleAssignmentScheduleRequestProperties(input.ScheduleInfo[0], output)
	} else {
		output.ScheduleInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo{}
		output.ScheduleInfo.Expiration = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration{}
		output.ScheduleInfo.Expiration.Type = pointer.To(roleassignmentschedulerequests.TypeNoExpiration)
	}
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaTicketInfoToRoleAssignmentScheduleRequestProperties(input PimActiveRoleAssignmentResourceSchemaTicketInfo, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties) {
	if output.TicketInfo == nil {
		output.TicketInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo{}
	}
	r.mapPimActiveRoleAssignmentResourceSchemaTicketInfoToRoleAssignmentScheduleRequestPropertiesTicketInfo(input, output.TicketInfo)
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaTicketInfoToRoleAssignmentScheduleRequestPropertiesTicketInfo(input PimActiveRoleAssignmentResourceSchemaTicketInfo, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo) {
	output.TicketNumber = pointer.To(input.TicketNumber)
	output.TicketSystem = pointer.To(input.TicketSystem)
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestResourceScheduleInfoSchemaToRoleAssignmentScheduleRequestProperties(input PimActiveRoleAssignmentResourceSchemaScheduleInfo, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties) {
	if output.ScheduleInfo == nil {
		output.ScheduleInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo{}
	}
	r.mapPimActiveRoleAssignmentResourceSchemaScheduleInfoToRoleAssignmentScheduleRequestPropertiesScheduleInfo(input, output.ScheduleInfo)
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaScheduleInfoToRoleAssignmentScheduleRequestPropertiesScheduleInfo(input PimActiveRoleAssignmentResourceSchemaScheduleInfo, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo) {
	output.StartDateTime = pointer.To(input.StartDateTime)

	if output.Expiration == nil {
		output.Expiration = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration{}
	}
	if len(input.Expiration) > 0 {
		r.mapPimActiveRoleAssignmentResourceSchemaScheduleInfoExpirationToRoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration(input.Expiration[0], output.Expiration)
	} else {
		// when no expiration is specified, set the type to No Expiration
		output.Expiration.Type = pointer.To(roleassignmentschedulerequests.TypeNoExpiration)
	}
}

func (r PimActiveRoleAssignmentResource) mapPimActiveRoleAssignmentResourceSchemaScheduleInfoExpirationToRoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration(input PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration, output *roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration) {
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
		output.Type = pointer.To(roleassignmentschedulerequests.TypeAfterDuration)

	case input.EndDateTime != "":
		output.Type = pointer.To(roleassignmentschedulerequests.TypeAfterDateTime)

	default:
		output.Type = pointer.To(roleassignmentschedulerequests.TypeNoExpiration)
	}
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestToPimActiveRoleAssignmentResourceSchema(input roleassignmentschedulerequests.RoleAssignmentScheduleRequest, output *PimActiveRoleAssignmentResourceSchema) error {
	if input.Properties == nil {
		input.Properties = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties{}
	}

	if err := r.mapRoleAssignmentScheduleRequestPropertiesToPimActiveRoleAssignmentResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "RoleAssignmentScheduleRequestProperties", "Properties", err)
	}

	return nil
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestPropertiesToPimActiveRoleAssignmentResourceSchema(input roleassignmentschedulerequests.RoleAssignmentScheduleRequestProperties, output *PimActiveRoleAssignmentResourceSchema) error {
	output.Justification = pointer.From(input.Justification)
	output.PrincipalId = input.PrincipalId
	output.PrincipalType = string(*input.PrincipalType)
	output.RoleDefinitionId = input.RoleDefinitionId

	if input.TicketInfo == nil {
		input.TicketInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo{}
	}

	if input.ScheduleInfo == nil {
		input.ScheduleInfo = &roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo{}
	}

	if input.TicketInfo != nil && (input.TicketInfo.TicketNumber != nil ||
		input.TicketInfo.TicketSystem != nil) {
		tempTicketInfo := &PimActiveRoleAssignmentResourceSchemaTicketInfo{}
		if err := r.mapRoleAssignmentScheduleRequestPropertiesTicketInfoToPimActiveRoleAssignmentResourceSchemaTicketInfo(*input.TicketInfo, tempTicketInfo); err != nil {
			return err
		} else {
			output.TicketInfo = make([]PimActiveRoleAssignmentResourceSchemaTicketInfo, 0)
			output.TicketInfo = append(output.TicketInfo, *tempTicketInfo)
		}
	}

	tempScheduleInfo := &PimActiveRoleAssignmentResourceSchemaScheduleInfo{}
	if err := r.mapRoleAssignmentScheduleRequestPropertiesScheduleInfoToPimActiveRoleAssignmentResourceSchemaScheduleInfo(*input.ScheduleInfo, tempScheduleInfo); err != nil {
		return err
	} else {
		output.ScheduleInfo = make([]PimActiveRoleAssignmentResourceSchemaScheduleInfo, 0)
		output.ScheduleInfo = append(output.ScheduleInfo, *tempScheduleInfo)
	}

	return nil
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestPropertiesTicketInfoToPimActiveRoleAssignmentResourceSchemaTicketInfo(input roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesTicketInfo, output *PimActiveRoleAssignmentResourceSchemaTicketInfo) error {
	output.TicketNumber = pointer.From(input.TicketNumber)
	output.TicketSystem = pointer.From(input.TicketSystem)
	return nil
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestPropertiesScheduleInfoToPimActiveRoleAssignmentResourceSchemaScheduleInfo(input roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfo, output *PimActiveRoleAssignmentResourceSchemaScheduleInfo) error {
	output.StartDateTime = pointer.From(input.StartDateTime)

	tempExpiration := &PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration{}
	if err := r.mapRoleAssignmentScheduleRequestPropertiesExpirationToPimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration(*input.Expiration, tempExpiration); err != nil {
		return err
	} else {
		output.Expiration = make([]PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration, 0)
		output.Expiration = append(output.Expiration, *tempExpiration)
	}
	return nil
}

func (r PimActiveRoleAssignmentResource) mapRoleAssignmentScheduleRequestPropertiesExpirationToPimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration(input roleassignmentschedulerequests.RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration, output *PimActiveRoleAssignmentResourceSchemaScheduleInfoExpiration) error {

	if input.Duration != nil && *input.Duration != "" {
		durationRaw := *input.Duration

		reHours := regexp.MustCompile(`PT(\d+)H`)
		matches := reHours.FindStringSubmatch(durationRaw)
		if len(matches) == 2 {
			hours, err := strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("could not decode hours from RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration: %+v", err)
			}
			output.DurationHours = hours
		}
		reDays := regexp.MustCompile(`P(\d+)D`)
		matches = reDays.FindStringSubmatch(durationRaw)
		if len(matches) == 2 {
			days, err := strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("could not decode days from RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration: %+v", err)
			}
			output.DurationDays = days
		}
	}

	output.EndDateTime = pointer.From(input.EndDateTime)
	return nil
}

func createActiveRoleAssignment(ctx context.Context, client *roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient, id roleassignmentschedulerequests.ScopedRoleAssignmentScheduleRequestId, payload *roleassignmentschedulerequests.RoleAssignmentScheduleRequest) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {

		// Azure can error when the subject doesn't exist yet due to AAD replication
		// Retry deletes while that error exists.
		result, err := client.Create(ctx, id, *payload)
		if err != nil {
			if *result.OData.Error.Code == "SubjectNotFound" {
				return nil, "Exist", nil
			}

			return nil, "Exist", fmt.Errorf("creating %s: %+v", id, err)
		}

		return result, "Created", nil
	}
}

func waitForActiveRoleAssignment(ctx context.Context, client *roleassignmentscheduleinstances.RoleAssignmentScheduleInstancesClient, scope string, principalId string, roleDefinitionId string, target string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if Role Assignment is %s on %q with role %q for %q.", target, scope, roleDefinitionId, principalId)

		instanceId := commonids.NewScopeID(scope)
		filter := &roleassignmentscheduleinstances.ListForScopeOperationOptions{
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
				*item.Properties.MemberType == roleassignmentscheduleinstances.MemberTypeDirect {
				state = "Found"
				result = item
			}
		}

		if target == "Missing" && state == "Missing" {
			result = &roleassignmentscheduleinstances.RoleAssignmentScheduleInstance{}
		}

		return result, state, nil
	}
}

func deleteActiveRoleAssignment(ctx context.Context, client *roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient, id roleassignmentschedulerequests.ScopedRoleAssignmentScheduleRequestId, payload *roleassignmentschedulerequests.RoleAssignmentScheduleRequest) pluginsdk.StateRefreshFunc {
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
