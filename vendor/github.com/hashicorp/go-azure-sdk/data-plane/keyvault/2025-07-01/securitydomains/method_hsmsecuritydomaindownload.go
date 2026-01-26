package securitydomains

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HSMSecurityDomainDownloadOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SecurityDomainObject
}

// HSMSecurityDomainDownload ...
func (c SecuritydomainsClient) HSMSecurityDomainDownload(ctx context.Context, input CertificateInfoObject) (result HSMSecurityDomainDownloadOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPost,
		Path:       "/securitydomain/download",
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

	result.Poller, err = dataplane.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// HSMSecurityDomainDownloadThenPoll performs HSMSecurityDomainDownload then polls until it's completed
func (c SecuritydomainsClient) HSMSecurityDomainDownloadThenPoll(ctx context.Context, input CertificateInfoObject) error {
	result, err := c.HSMSecurityDomainDownload(ctx, input)
	if err != nil {
		return fmt.Errorf("performing HSMSecurityDomainDownload: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after HSMSecurityDomainDownload: %+v", err)
	}

	return nil
}
