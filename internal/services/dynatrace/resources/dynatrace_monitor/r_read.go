// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				identityProps, err := flattenDynatraceIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening identity: %+v", err)
				}
				monitoringStatus := true
				if *props.MonitoringStatus == monitors.MonitoringStatusDisabled {
					monitoringStatus = false
				}

				state := MonitorsResourceModel{
					Name:                          id.MonitorName,
					ResourceGroup:                 id.ResourceGroupName,
					Location:                      model.Location,
					MonitoringStatus:              monitoringStatus,
					MarketplaceSubscriptionStatus: string(*props.MarketplaceSubscriptionStatus),
					Identity:                      identityProps,
					PlanData:                      FlattenDynatracePlanData(props.PlanData),
					UserInfo:                      FlattenDynatraceUserInfo(props.UserInfo),
				}

				if environmentProps := metadata.ResourceData.Get("environment_properties"); environmentProps != nil && len(environmentProps.([]any)) > 0 {
					state.EnvironmentProperties = FlattenDynatraceEnvironmentProperties(props.DynatraceEnvironmentProperties)
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func flattenDynatraceIdentity(input *monitors.IdentityProperties) ([]identity.ModelSystemAssigned, error) {
	if input == nil {
		return nil, errors.New("flattening Dynatrace identity: input is nil")
	}
	var identityProp identity.ModelSystemAssigned

	identityProp.Type = identity.Type(input.Type)

	if input.PrincipalId != nil {
		identityProp.PrincipalId = pointer.From(input.PrincipalId)
	}

	if input.TenantId != nil {
		identityProp.TenantId = pointer.From(input.TenantId)
	}

	return []identity.ModelSystemAssigned{
		identityProp,
	}, nil
}

type PlanData struct {
	BillingCycle  string `tfschema:"billing_cycle"`
	EffectiveDate string `tfschema:"effective_date"`
	PlanDetails   string `tfschema:"plan"`
	UsageType     string `tfschema:"usage_type"`
}

type UserInfo struct {
	Country      string `tfschema:"country"`
	EmailAddress string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	PhoneNumber  string `tfschema:"phone_number"`
}
