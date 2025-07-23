package virtualnetworkgateways

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

type GetFailoverSingleTestDetailsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteFailoverSingleTestDetails
}

type GetFailoverSingleTestDetailsOperationOptions struct {
	FailoverTestId  *string
	PeeringLocation *string
}

func DefaultGetFailoverSingleTestDetailsOperationOptions() GetFailoverSingleTestDetailsOperationOptions {
	return GetFailoverSingleTestDetailsOperationOptions{}
}

func (o GetFailoverSingleTestDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetFailoverSingleTestDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetFailoverSingleTestDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FailoverTestId != nil {
		out.Append("failoverTestId", fmt.Sprintf("%v", *o.FailoverTestId))
	}
	if o.PeeringLocation != nil {
		out.Append("peeringLocation", fmt.Sprintf("%v", *o.PeeringLocation))
	}
	return &out
}

// GetFailoverSingleTestDetails ...
func (c VirtualNetworkGatewaysClient) GetFailoverSingleTestDetails(ctx context.Context, id VirtualNetworkGatewayId, options GetFailoverSingleTestDetailsOperationOptions) (result GetFailoverSingleTestDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getFailoverSingleTestDetails", id.ID()),
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

// GetFailoverSingleTestDetailsThenPoll performs GetFailoverSingleTestDetails then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetFailoverSingleTestDetailsThenPoll(ctx context.Context, id VirtualNetworkGatewayId, options GetFailoverSingleTestDetailsOperationOptions) error {
	result, err := c.GetFailoverSingleTestDetails(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetFailoverSingleTestDetails: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetFailoverSingleTestDetails: %+v", err)
	}

	return nil
}
