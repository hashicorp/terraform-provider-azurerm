// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &recoverKeyPoller{}

func NewRecoverKeyPoller(uri string) pollers.PollerType {
	return &recoverKeyPoller{
		uri: uri,
	}
}

type recoverKeyPoller struct {
	uri string
}

func (p *recoverKeyPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {

	res := &pollers.PollResult{
		PollInterval: time.Second * 20,
		Status:       pollers.PollingStatusInProgress,
	}
	conn, err := http.Get(p.uri)
	if err != nil {
		log.Printf("[DEBUG] Didn't find KeyVault secret at %q", p.uri)
		return res, fmt.Errorf("checking secret at %q: %s", p.uri, err)
	}

	defer conn.Body.Close()
	if response.WasNotFound(conn) {
		res.Status = pollers.PollingStatusSucceeded
		return res, nil
	}

	res.Status = pollers.PollingStatusSucceeded
	return res, nil
}
