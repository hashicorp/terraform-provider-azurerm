package expressroutecircuits

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

type StartCircuitLinkFailoverTestOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

type StartCircuitLinkFailoverTestOperationOptions struct {
	CircuitTestCategory *string
	LinkType            *string
}

func DefaultStartCircuitLinkFailoverTestOperationOptions() StartCircuitLinkFailoverTestOperationOptions {
	return StartCircuitLinkFailoverTestOperationOptions{}
}

func (o StartCircuitLinkFailoverTestOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o StartCircuitLinkFailoverTestOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o StartCircuitLinkFailoverTestOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CircuitTestCategory != nil {
		out.Append("circuitTestCategory", fmt.Sprintf("%v", *o.CircuitTestCategory))
	}
	if o.LinkType != nil {
		out.Append("linkType", fmt.Sprintf("%v", *o.LinkType))
	}
	return &out
}

// StartCircuitLinkFailoverTest ...
func (c ExpressRouteCircuitsClient) StartCircuitLinkFailoverTest(ctx context.Context, id ExpressRouteCircuitId, options StartCircuitLinkFailoverTestOperationOptions) (result StartCircuitLinkFailoverTestOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/startCircuitLinkFailoverTest", id.ID()),
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

// StartCircuitLinkFailoverTestThenPoll performs StartCircuitLinkFailoverTest then polls until it's completed
func (c ExpressRouteCircuitsClient) StartCircuitLinkFailoverTestThenPoll(ctx context.Context, id ExpressRouteCircuitId, options StartCircuitLinkFailoverTestOperationOptions) error {
	return c.StartCircuitLinkFailoverTestCallbackThenPoll(ctx, id, options, nil)
}

// StartCircuitLinkFailoverTestCallbackThenPoll performs StartCircuitLinkFailoverTest, runs the optional callback function, then polls until it's completed
func (c ExpressRouteCircuitsClient) StartCircuitLinkFailoverTestCallbackThenPoll(ctx context.Context, id ExpressRouteCircuitId, options StartCircuitLinkFailoverTestOperationOptions, callback func() error) error {
	result, err := c.StartCircuitLinkFailoverTest(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing StartCircuitLinkFailoverTest: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StartCircuitLinkFailoverTest: %+v", err)
	}

	return nil
}
