package customdomains

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

type DisableCustomHTTPSOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CustomDomain
}

// DisableCustomHTTPS ...
func (c CustomDomainsClient) DisableCustomHTTPS(ctx context.Context, id EndpointCustomDomainId) (result DisableCustomHTTPSOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/disableCustomHttps", id.ID()),
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

// DisableCustomHTTPSThenPoll performs DisableCustomHTTPS then polls until it's completed
func (c CustomDomainsClient) DisableCustomHTTPSThenPoll(ctx context.Context, id EndpointCustomDomainId) error {
	return c.DisableCustomHTTPSCallbackThenPoll(ctx, id, nil)
}

// DisableCustomHTTPSCallbackThenPoll performs DisableCustomHTTPS, runs the optional callback function, then polls until it's completed
func (c CustomDomainsClient) DisableCustomHTTPSCallbackThenPoll(ctx context.Context, id EndpointCustomDomainId, callback func() error) error {
	result, err := c.DisableCustomHTTPS(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DisableCustomHTTPS: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DisableCustomHTTPS: %+v", err)
	}

	return nil
}
