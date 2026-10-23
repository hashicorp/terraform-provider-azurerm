// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package new_relic_monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenPlanDataModel(input *monitors.PlanData) []PlanDataModel {
	var outputList []PlanDataModel
	if input == nil {
		return []PlanDataModel{}
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
		return []UserInfoModel{}
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

const (
	PlanSuffix = "@TIDgmz7xq9ge3py@PUBIDnewrelicinc1635200720692.newrelic_liftr_payg"
)
