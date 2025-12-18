// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)
