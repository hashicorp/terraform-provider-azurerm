// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.MonitorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MonitorsResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroup, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			monitoringStatus := monitors.MonitoringStatusEnabled
			if !model.MonitoringStatus {
				monitoringStatus = monitors.MonitoringStatusDisabled
			}
			monitorsProps := monitors.MonitorProperties{
				MarketplaceSubscriptionStatus:  pointer.ToEnum[monitors.MarketplaceSubscriptionStatus](model.MarketplaceSubscriptionStatus),
				MonitoringStatus:               pointer.To(monitoringStatus),
				PlanData:                       ExpandDynatracePlanData(model.PlanData),
				UserInfo:                       ExpandDynatraceUserInfo(model.UserInfo),
				DynatraceEnvironmentProperties: ExpandDynatraceEnvironmentProperties(model.EnvironmentProperties),
			}

			dynatraceIdentity, err := expandDynatraceIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			monitor := monitors.MonitorResource{
				Identity:   dynatraceIdentity,
				Location:   model.Location,
				Name:       pointer.To(model.Name),
				Properties: monitorsProps,
				Tags:       pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, monitor, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandDynatraceIdentity(input []identity.ModelSystemAssigned) (*monitors.IdentityProperties, error) {
	config, err := identity.ExpandSystemAssignedFromModel(input)
	if err != nil {
		return nil, err
	}

	if config.Type == identity.TypeNone {
		return pointer.To(monitors.IdentityProperties{}), nil
	}

	dynatraceIdentity := monitors.IdentityProperties{
		Type: monitors.ManagedIdentityType(config.Type),
	}

	return pointer.To(dynatraceIdentity), nil
}
