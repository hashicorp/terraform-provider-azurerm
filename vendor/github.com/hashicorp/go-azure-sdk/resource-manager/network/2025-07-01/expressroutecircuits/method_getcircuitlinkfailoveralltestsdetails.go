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

type GetCircuitLinkFailoverAllTestsDetailsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteLinkFailoverAllTestsDetails
}

type GetCircuitLinkFailoverAllTestsDetailsOperationOptions struct {
	FailoverTestType *string
	FetchLatest      *bool
}

func DefaultGetCircuitLinkFailoverAllTestsDetailsOperationOptions() GetCircuitLinkFailoverAllTestsDetailsOperationOptions {
	return GetCircuitLinkFailoverAllTestsDetailsOperationOptions{}
}

func (o GetCircuitLinkFailoverAllTestsDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCircuitLinkFailoverAllTestsDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetCircuitLinkFailoverAllTestsDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FailoverTestType != nil {
		out.Append("failoverTestType", fmt.Sprintf("%v", *o.FailoverTestType))
	}
	if o.FetchLatest != nil {
		out.Append("fetchLatest", fmt.Sprintf("%v", *o.FetchLatest))
	}
	return &out
}

// GetCircuitLinkFailoverAllTestsDetails ...
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverAllTestsDetails(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverAllTestsDetailsOperationOptions) (result GetCircuitLinkFailoverAllTestsDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getCircuitLinkFailoverAllTestsDetails", id.ID()),
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

// GetCircuitLinkFailoverAllTestsDetailsThenPoll performs GetCircuitLinkFailoverAllTestsDetails then polls until it's completed
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverAllTestsDetailsThenPoll(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverAllTestsDetailsOperationOptions) error {
	return c.GetCircuitLinkFailoverAllTestsDetailsCallbackThenPoll(ctx, id, options, nil)
}

// GetCircuitLinkFailoverAllTestsDetailsCallbackThenPoll performs GetCircuitLinkFailoverAllTestsDetails, runs the optional callback function, then polls until it's completed
func (c ExpressRouteCircuitsClient) GetCircuitLinkFailoverAllTestsDetailsCallbackThenPoll(ctx context.Context, id ExpressRouteCircuitId, options GetCircuitLinkFailoverAllTestsDetailsOperationOptions, callback func() error) error {
	result, err := c.GetCircuitLinkFailoverAllTestsDetails(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetCircuitLinkFailoverAllTestsDetails: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetCircuitLinkFailoverAllTestsDetails: %+v", err)
	}

	return nil
}
