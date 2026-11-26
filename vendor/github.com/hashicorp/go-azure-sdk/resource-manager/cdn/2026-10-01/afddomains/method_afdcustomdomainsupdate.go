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

type AFDCustomDomainsUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AFDDomain
}

// AFDCustomDomainsUpdate ...
func (c AFDDomainsClient) AFDCustomDomainsUpdate(ctx context.Context, id CustomDomainId, input AFDDomainUpdateParameters) (result AFDCustomDomainsUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
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

// AFDCustomDomainsUpdateThenPoll performs AFDCustomDomainsUpdate then polls until it's completed
func (c AFDDomainsClient) AFDCustomDomainsUpdateThenPoll(ctx context.Context, id CustomDomainId, input AFDDomainUpdateParameters) error {
	return c.AFDCustomDomainsUpdateCallbackThenPoll(ctx, id, input, nil)
}

// AFDCustomDomainsUpdateCallbackThenPoll performs AFDCustomDomainsUpdate, runs the optional callback function, then polls until it's completed
func (c AFDDomainsClient) AFDCustomDomainsUpdateCallbackThenPoll(ctx context.Context, id CustomDomainId, input AFDDomainUpdateParameters, callback func() error) error {
	result, err := c.AFDCustomDomainsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AFDCustomDomainsUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AFDCustomDomainsUpdate: %+v", err)
	}

	return nil
}
