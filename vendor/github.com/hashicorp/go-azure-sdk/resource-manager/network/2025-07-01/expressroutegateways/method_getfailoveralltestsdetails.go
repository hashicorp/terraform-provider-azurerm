package expressroutegateways

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

type GetFailoverAllTestsDetailsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteFailoverTestDetails
}

type GetFailoverAllTestsDetailsOperationOptions struct {
	FetchLatest *bool
	Type        *string
}

func DefaultGetFailoverAllTestsDetailsOperationOptions() GetFailoverAllTestsDetailsOperationOptions {
	return GetFailoverAllTestsDetailsOperationOptions{}
}

func (o GetFailoverAllTestsDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetFailoverAllTestsDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetFailoverAllTestsDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FetchLatest != nil {
		out.Append("fetchLatest", fmt.Sprintf("%v", *o.FetchLatest))
	}
	if o.Type != nil {
		out.Append("type", fmt.Sprintf("%v", *o.Type))
	}
	return &out
}

// GetFailoverAllTestsDetails ...
func (c ExpressRouteGatewaysClient) GetFailoverAllTestsDetails(ctx context.Context, id ExpressRouteGatewayId, options GetFailoverAllTestsDetailsOperationOptions) (result GetFailoverAllTestsDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getFailoverAllTestsDetails", id.ID()),
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

// GetFailoverAllTestsDetailsThenPoll performs GetFailoverAllTestsDetails then polls until it's completed
func (c ExpressRouteGatewaysClient) GetFailoverAllTestsDetailsThenPoll(ctx context.Context, id ExpressRouteGatewayId, options GetFailoverAllTestsDetailsOperationOptions) error {
	return c.GetFailoverAllTestsDetailsCallbackThenPoll(ctx, id, options, nil)
}

// GetFailoverAllTestsDetailsCallbackThenPoll performs GetFailoverAllTestsDetails, runs the optional callback function, then polls until it's completed
func (c ExpressRouteGatewaysClient) GetFailoverAllTestsDetailsCallbackThenPoll(ctx context.Context, id ExpressRouteGatewayId, options GetFailoverAllTestsDetailsOperationOptions, callback func() error) error {
	result, err := c.GetFailoverAllTestsDetails(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetFailoverAllTestsDetails: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetFailoverAllTestsDetails: %+v", err)
	}

	return nil
}
