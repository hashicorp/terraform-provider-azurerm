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

type GetCircuitLinkFailoverSingleTestDetailsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteLinkFailoverSingleTestDetails
}

type GetCircuitLinkFailoverSingleTestDetailsOperationOptions struct {
	CircuitTestCategory *string
	FailoverTestId      *string
	LinkType            *string
}

func DefaultGetCircuitLinkFailoverSingleTestDetailsOperationOptions() GetCircuitLinkFailoverSingleTestDetailsOperationOptions {
	return GetCircuitLinkFailoverSingleTestDetailsOperationOptions{}
}

func (o GetCircuitLinkFailoverSingleTestDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCircuitLinkFailoverSingleTestDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetCircuitLinkFailoverSingleTestDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CircuitTestCategory != nil {
		out.Append("circuitTestCategory", fmt.Sprintf("%v", *o.CircuitTestCategory))
	}
	if o.FailoverTestId != nil {
		out.Append("failoverTestId", fmt.Sprintf("%v", *o.FailoverTestId))
	}
	if o.LinkType != nil {
		out.Append("linkType", fmt.Sprintf("%v", *o.LinkType))
	}
	return &out
}

// GetCircuitLinkFailoverSingleTestDetails ...
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverSingleTestDetails(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverSingleTestDetailsOperationOptions) (result GetCircuitLinkFailoverSingleTestDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getCircuitLinkFailoverSingleTestDetails", id.ID()),
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

// GetCircuitLinkFailoverSingleTestDetailsThenPoll performs GetCircuitLinkFailoverSingleTestDetails then polls until it's completed
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverSingleTestDetailsThenPoll(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverSingleTestDetailsOperationOptions) error {
	return c.GetCircuitLinkFailoverSingleTestDetailsCallbackThenPoll(ctx, id, options, nil)
}

// GetCircuitLinkFailoverSingleTestDetailsCallbackThenPoll performs GetCircuitLinkFailoverSingleTestDetails, runs the optional callback function, then polls until it's completed
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverSingleTestDetailsCallbackThenPoll(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverSingleTestDetailsOperationOptions, callback func() error) error {
	result, err := c.GetCircuitLinkFailoverSingleTestDetails(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetCircuitLinkFailoverSingleTestDetails: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetCircuitLinkFailoverSingleTestDetails: %+v", err)
	}

	return nil
}
