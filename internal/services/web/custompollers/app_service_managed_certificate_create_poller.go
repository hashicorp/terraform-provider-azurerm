package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &AppServiceManagedCertificateCreatePoller{}

type AppServiceManagedCertificateCreatePoller struct {
	client *certificates.CertificatesClient
	id     certificates.CertificateId
}

func NewAppServiceManagedCertificateCreatePoller(client *certificates.CertificatesClient, id certificates.CertificateId) AppServiceManagedCertificateCreatePoller {
	return AppServiceManagedCertificateCreatePoller{
		client: client,
		id:     id,
	}
}

func (p AppServiceManagedCertificateCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				Status:       pollers.PollingStatusInProgress,
				PollInterval: 10 * time.Second,
			}, nil
		}
		return nil, fmt.Errorf("retrieving %s: %w", p.id, err)
	}

	return &pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}, nil
}
