// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			identityValue, err := identity.ExpandSystemAssigned(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			// Currently the API does not accept `None` type: https://github.com/Azure/azure-rest-api-specs/issues/29257
			if identityValue.Type != identity.TypeNone {
				properties.Identity = identityValue
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

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
