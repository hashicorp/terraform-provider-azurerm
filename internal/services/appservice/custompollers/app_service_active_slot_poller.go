// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &appServiceActiveSlotPoller{}

type appServiceActiveSlotPoller struct {
	client *webapps.WebAppsClient
	id     webapps.SlotId
	appId  commonids.AppServiceId
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func NewAppServiceActiveSlotPoller(client *webapps.WebAppsClient, id commonids.AppServiceId, slotId webapps.SlotId) *appServiceActiveSlotPoller {
	return &appServiceActiveSlotPoller{
		client: client,
		id:     slotId,
		appId:  id,
	}
}

func (p appServiceActiveSlotPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.appId)
	if err == nil {
		if resp.Model != nil && resp.Model.Properties != nil {
			swapStatus := resp.Model.Properties.SlotSwapStatus
			if swapStatus == nil || pointer.From(swapStatus.SourceSlotName) != p.id.SlotName {
				return &pollingInProgress, err
			}
			return &pollingSuccess, nil
		}
	}
	return nil, fmt.Errorf("retrieving %s: %+v", p.appId, err)
}
