// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	PlanSuffix = "@TIDgmz7xq9ge3py@PUBIDnewrelicinc1635200720692.newrelic_liftr_payg"
)

type NewRelicMonitorModel struct {
	Name                  string                         `tfschema:"name"`
	ResourceGroupName     string                         `tfschema:"resource_group_name"`
	AccountCreationSource monitors.AccountCreationSource `tfschema:"account_creation_source"`
	AccountId             string                         `tfschema:"account_id"`
	IngestionKey          string                         `tfschema:"ingestion_key"`
	Location              string                         `tfschema:"location"`
	OrganizationId        string                         `tfschema:"organization_id"`
	OrgCreationSource     monitors.OrgCreationSource     `tfschema:"org_creation_source"`
	PlanData              []PlanDataModel                `tfschema:"plan"`
	UserId                string                         `tfschema:"user_id"`
	UserInfo              []UserInfoModel                `tfschema:"user"`
}

type PlanDataModel struct {
	EffectiveDate string             `tfschema:"effective_date"`
	BillingCycle  string             `tfschema:"billing_cycle"`
	PlanDetails   string             `tfschema:"plan_id"`
	UsageType     monitors.UsageType `tfschema:"usage_type"`
}

type UserInfoModel struct {
	EmailAddress string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	PhoneNumber  string `tfschema:"phone_number"`
}

type NewRelicMonitorResource struct{}

var _ sdk.Resource = NewRelicMonitorResource{}

func (r NewRelicMonitorResource) ResourceType() string {
	return "azurerm_new_relic_monitor"
}

func (r NewRelicMonitorResource) ModelObject() interface{} {
	return &NewRelicMonitorModel{}
}

func (r NewRelicMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return monitors.ValidateMonitorID
}

func (r NewRelicMonitorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[\w\-]{1,32}$`),
				`name must be between 1 and 32 characters in length and may contain only letters, numbers, hyphens and underscores`,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"effective_date": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ForceNew:         true,
						DiffSuppressFunc: suppress.RFC3339Time,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					// Enum is removed https://github.com/Azure/azure-rest-api-specs/issues/31093
					"billing_cycle": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "MONTHLY",
						ValidateFunc: validation.StringInSlice([]string{
							"MONTHLY",
							"WEEKLY",
							"YEARLY",
						}, false),
					},

					"plan_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "newrelic-pay-as-you-go-free-live",
						ValidateFunc: validation.StringInSlice([]string{
							"newrelic-pay-as-you-go-free-live",
						}, false),
					},

					"usage_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(monitors.UsageTypePAYG),
						ValidateFunc: validation.StringInSlice([]string{
							string(monitors.UsageTypeCOMMITTED),
							string(monitors.UsageTypePAYG),
						}, false),
					},
				},
			},
		},

		"user": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"first_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone_number": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"account_creation_source": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(monitors.AccountCreationSourceLIFTR),
			ValidateFunc: validation.StringInSlice([]string{
				string(monitors.AccountCreationSourceLIFTR),
				string(monitors.AccountCreationSourceNEWRELIC),
			}, false),
		},

		"account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedIdentityOptionalForceNew(),

		"ingestion_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"organization_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"org_creation_source": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(monitors.OrgCreationSourceLIFTR),
			ValidateFunc: validation.StringInSlice([]string{
				string(monitors.OrgCreationSourceLIFTR),
				string(monitors.OrgCreationSourceNEWRELIC),
			}, false),
		},

		"user_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r NewRelicMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NewRelicMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NewRelicMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NewRelic.MonitorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &monitors.NewRelicMonitorResource{
				Location: location.Normalize(model.Location),
				Properties: monitors.MonitorProperties{
					AccountCreationSource:     &model.AccountCreationSource,
					NewRelicAccountProperties: expandNewRelicAccountPropertiesModel(&model),
					OrgCreationSource:         &model.OrgCreationSource,
					PlanData:                  expandPlanDataModel(model.PlanData),
					UserInfo:                  expandUserInfoModel(model.UserInfo),
				},
			}

			identityValue, err := identity.ExpandSystemAssigned(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			// Currently the API does not accept `None` type: https://github.com/Azure/azure-rest-api-specs/issues/29257
			if identityValue.Type != identity.TypeNone {
				properties.Identity = identityValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NewRelicMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitorsClient

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var originalModel NewRelicMonitorModel
			if err = metadata.Decode(&originalModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := NewRelicMonitorModel{
				Name:              id.MonitorName,
				ResourceGroupName: id.ResourceGroupName,
				IngestionKey:      originalModel.IngestionKey,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if err := metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				properties := &model.Properties
				if properties.AccountCreationSource != nil {
					state.AccountCreationSource = *properties.AccountCreationSource
				}

				if properties.NewRelicAccountProperties != nil {
					if properties.NewRelicAccountProperties.AccountInfo != nil {
						state.AccountId = pointer.From(properties.NewRelicAccountProperties.AccountInfo.AccountId)
					}

					if properties.NewRelicAccountProperties.OrganizationInfo != nil {
						state.OrganizationId = pointer.From(properties.NewRelicAccountProperties.OrganizationInfo.OrganizationId)
					}

					state.UserId = pointer.From(properties.NewRelicAccountProperties.UserId)
				}

				if properties.OrgCreationSource != nil {
					state.OrgCreationSource = *properties.OrgCreationSource
				}

				state.PlanData = flattenPlanDataModel(properties.PlanData)

				state.UserInfo = flattenUserInfoModel(properties.UserInfo)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NewRelicMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NewRelic.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NewRelicMonitorModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.UserInfo) == 0 {
				return fmt.Errorf("deleting: `user` not found")
			}

			if err = client.DeleteThenPoll(ctx, *id, monitors.DeleteOperationOptions{UserEmail: &model.UserInfo[0].EmailAddress}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandNewRelicAccountPropertiesModel(model *NewRelicMonitorModel) *monitors.NewRelicAccountProperties {
	output := monitors.NewRelicAccountProperties{}
	containsValue := false

	accountInfo := expandAccountInfoModel(model)
	if accountInfo != nil {
		output.AccountInfo = accountInfo
		containsValue = true
	}

	organizationInfo := expandOrganizationInfoModel(model)
	if organizationInfo != nil {
		output.OrganizationInfo = organizationInfo
		containsValue = true
	}

	if model.UserId != "" {
		output.UserId = &model.UserId
		containsValue = true
	}

	if !containsValue {
		return nil
	}

	return &output
}

func expandAccountInfoModel(model *NewRelicMonitorModel) *monitors.AccountInfo {
	output := monitors.AccountInfo{}
	containsValue := false

	if model.AccountId != "" {
		output.AccountId = &model.AccountId
		containsValue = true
	}

	if model.IngestionKey != "" {
		output.IngestionKey = &model.IngestionKey
		containsValue = true
	}

	if !containsValue {
		return nil
	}

	return &output
}

func expandOrganizationInfoModel(model *NewRelicMonitorModel) *monitors.OrganizationInfo {
	if model.OrganizationId == "" {
		return nil
	}

	return &monitors.OrganizationInfo{
		OrganizationId: &model.OrganizationId,
	}
}

func expandPlanDataModel(inputList []PlanDataModel) *monitors.PlanData {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := monitors.PlanData{
		BillingCycle:  &input.BillingCycle,
		EffectiveDate: &input.EffectiveDate,
		UsageType:     &input.UsageType,
	}
	if input.PlanDetails != "" {
		output.PlanDetails = pointer.To(input.PlanDetails + PlanSuffix)
	}

	return &output
}

func expandUserInfoModel(inputList []UserInfoModel) *monitors.UserInfo {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := monitors.UserInfo{}

	if input.EmailAddress != "" {
		output.EmailAddress = &input.EmailAddress
	}

	if input.FirstName != "" {
		output.FirstName = &input.FirstName
	}

	if input.LastName != "" {
		output.LastName = &input.LastName
	}

	if input.PhoneNumber != "" {
		output.PhoneNumber = &input.PhoneNumber
	}

	return &output
}

func flattenPlanDataModel(input *monitors.PlanData) []PlanDataModel {
	var outputList []PlanDataModel
	if input == nil {
		return outputList
	}
	output := PlanDataModel{}
	if input.BillingCycle != nil {
		output.BillingCycle = *input.BillingCycle
	}

	if input.EffectiveDate != nil {
		output.EffectiveDate = *input.EffectiveDate
	}

	if input.PlanDetails != nil {
		output.PlanDetails = strings.TrimSuffix(*input.PlanDetails, PlanSuffix)
	}

	if input.UsageType != nil {
		output.UsageType = *input.UsageType
	}

	return append(outputList, output)
}

func flattenUserInfoModel(input *monitors.UserInfo) []UserInfoModel {
	var outputList []UserInfoModel
	if input == nil {
		return outputList
	}
	output := UserInfoModel{}

	if input.EmailAddress != nil {
		output.EmailAddress = *input.EmailAddress
	}

	if input.FirstName != nil {
		output.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		output.LastName = *input.LastName
	}

	if input.PhoneNumber != nil {
		output.PhoneNumber = *input.PhoneNumber
	}

	return append(outputList, output)
}
