package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &aadDiagnosticSettingCreatePoller{}

type aadDiagnosticSettingCreatePoller struct {
	client *diagnosticsettings.DiagnosticSettingsClient
	id     diagnosticsettings.DiagnosticSettingId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewAadDiagnosticSettingCreatePoller(client *diagnosticsettings.DiagnosticSettingsClient, id diagnosticsettings.DiagnosticSettingId) *aadDiagnosticSettingCreatePoller {
	return &aadDiagnosticSettingCreatePoller{
		client: client,
		id:     id,
	}
}

func (p aadDiagnosticSettingCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingInProgress, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollingSuccess, nil
}
