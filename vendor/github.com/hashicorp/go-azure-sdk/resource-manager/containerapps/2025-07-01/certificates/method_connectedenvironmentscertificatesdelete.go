package certificates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsCertificatesDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ConnectedEnvironmentsCertificatesDelete ...
func (c CertificatesClient) ConnectedEnvironmentsCertificatesDelete(ctx context.Context, id ConnectedEnvironmentCertificateId) (result ConnectedEnvironmentsCertificatesDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// ConnectedEnvironmentsCertificatesDeleteThenPoll performs ConnectedEnvironmentsCertificatesDelete then polls until it's completed
func (c CertificatesClient) ConnectedEnvironmentsCertificatesDeleteThenPoll(ctx context.Context, id ConnectedEnvironmentCertificateId) error {
	result, err := c.ConnectedEnvironmentsCertificatesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConnectedEnvironmentsCertificatesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ConnectedEnvironmentsCertificatesDelete: %+v", err)
	}

	return nil
}
