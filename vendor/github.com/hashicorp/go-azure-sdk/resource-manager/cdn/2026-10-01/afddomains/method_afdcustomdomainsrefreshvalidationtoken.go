package afddomains

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

type AFDCustomDomainsRefreshValidationTokenOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// AFDCustomDomainsRefreshValidationToken ...
func (c AFDDomainsClient) AFDCustomDomainsRefreshValidationToken(ctx context.Context, id CustomDomainId) (result AFDCustomDomainsRefreshValidationTokenOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/refreshValidationToken", id.ID()),
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

// AFDCustomDomainsRefreshValidationTokenThenPoll performs AFDCustomDomainsRefreshValidationToken then polls until it's completed
func (c AFDDomainsClient) AFDCustomDomainsRefreshValidationTokenThenPoll(ctx context.Context, id CustomDomainId) error {
	return c.AFDCustomDomainsRefreshValidationTokenCallbackThenPoll(ctx, id, nil)
}

// AFDCustomDomainsRefreshValidationTokenCallbackThenPoll performs AFDCustomDomainsRefreshValidationToken, runs the optional callback function, then polls until it's completed
func (c AFDDomainsClient) AFDCustomDomainsRefreshValidationTokenCallbackThenPoll(ctx context.Context, id CustomDomainId, callback func() error) error {
	result, err := c.AFDCustomDomainsRefreshValidationToken(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AFDCustomDomainsRefreshValidationToken: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AFDCustomDomainsRefreshValidationToken: %+v", err)
	}

	return nil
}
