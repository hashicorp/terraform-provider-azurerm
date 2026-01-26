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

type HSMSecurityDomainUploadOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SecurityDomainOperationStatus
}

// HSMSecurityDomainUpload ...
func (c SecuritydomainsClient) HSMSecurityDomainUpload(ctx context.Context, input SecurityDomainObject) (result HSMSecurityDomainUploadOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       "/securitydomain/upload",
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

// HSMSecurityDomainUploadThenPoll performs HSMSecurityDomainUpload then polls until it's completed
func (c SecuritydomainsClient) HSMSecurityDomainUploadThenPoll(ctx context.Context, input SecurityDomainObject) error {
	result, err := c.HSMSecurityDomainUpload(ctx, input)
	if err != nil {
		return fmt.Errorf("performing HSMSecurityDomainUpload: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after HSMSecurityDomainUpload: %+v", err)
	}

	return nil
}
