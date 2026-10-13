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

type EnableCustomHTTPSOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CustomDomain
}

// EnableCustomHTTPS ...
func (c CustomDomainsClient) EnableCustomHTTPS(ctx context.Context, id EndpointCustomDomainId, input CustomDomainHTTPSParameters) (result EnableCustomHTTPSOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/enableCustomHttps", id.ID()),
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// EnableCustomHTTPSThenPoll performs EnableCustomHTTPS then polls until it's completed
func (c CustomDomainsClient) EnableCustomHTTPSThenPoll(ctx context.Context, id EndpointCustomDomainId, input CustomDomainHTTPSParameters) error {
	return c.EnableCustomHTTPSCallbackThenPoll(ctx, id, input, nil)
}

// EnableCustomHTTPSCallbackThenPoll performs EnableCustomHTTPS, runs the optional callback function, then polls until it's completed
func (c CustomDomainsClient) EnableCustomHTTPSCallbackThenPoll(ctx context.Context, id EndpointCustomDomainId, input CustomDomainHTTPSParameters, callback func() error) error {
	result, err := c.EnableCustomHTTPS(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing EnableCustomHTTPS: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after EnableCustomHTTPS: %+v", err)
	}

	return nil
}
