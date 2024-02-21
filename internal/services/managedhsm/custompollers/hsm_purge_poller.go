// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &hsmDownloadPoller{}

func NewHSMPurgePoller(client *managedhsms.ManagedHsmsClient, id managedhsms.DeletedManagedHSMId) pollers.PollerType {
	return &hsmPurgePoller{
		client:          client,
		purgeId:         id,
		purgeAgainUntil: time.Now().Add(time.Minute),
	}
}

type hsmPurgePoller struct {
	client          *managedhsms.ManagedHsmsClient
	purgeId         managedhsms.DeletedManagedHSMId
	purgeAgainUntil time.Time
}

func (p *hsmPurgePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	deletedResp, err := p.client.GetDeleted(ctx, p.purgeId)

	res := &pollers.PollResult{
		PollInterval: time.Second * 20,
		Status:       pollers.PollingStatusInProgress,
	}
	if response.WasNotFound(deletedResp.HttpResponse) {
		res.Status = pollers.PollingStatusSucceeded
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("retrieving deleted managed HSM %s: %+v", p.purgeId, err)
	}

	if time.Now().After(p.purgeAgainUntil) {
		p.purgeAgainUntil = time.Now().Add(time.Minute)
		purgeResp, _ := p.client.PurgeDeleted(ctx, p.purgeId)
		if response.WasNotFound(purgeResp.HttpResponse) {
			res.Status = pollers.PollingStatusSucceeded
		}
	}

	return res, nil
}
