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

type GetFailoverAllTestDetailsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExpressRouteFailoverTestDetails
}

type GetFailoverAllTestDetailsOperationOptions struct {
	FetchLatest *bool
	Type        *string
}

func DefaultGetFailoverAllTestDetailsOperationOptions() GetFailoverAllTestDetailsOperationOptions {
	return GetFailoverAllTestDetailsOperationOptions{}
}

func (o GetFailoverAllTestDetailsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetFailoverAllTestDetailsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetFailoverAllTestDetailsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.FetchLatest != nil {
		out.Append("fetchLatest", fmt.Sprintf("%v", *o.FetchLatest))
	}
	if o.Type != nil {
		out.Append("type", fmt.Sprintf("%v", *o.Type))
	}
	return &out
}

// GetFailoverAllTestDetails ...
func (c VirtualNetworkGatewaysClient) GetFailoverAllTestDetails(ctx context.Context, id VirtualNetworkGatewayId, options GetFailoverAllTestDetailsOperationOptions) (result GetFailoverAllTestDetailsOperationResponse, err error) {
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

// GetFailoverAllTestDetailsThenPoll performs GetFailoverAllTestDetails then polls until it's completed
func (c VirtualNetworkGatewaysClient) GetFailoverAllTestDetailsThenPoll(ctx context.Context, id VirtualNetworkGatewayId, options GetFailoverAllTestDetailsOperationOptions) error {
	result, err := c.GetFailoverAllTestDetails(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing GetFailoverAllTestDetails: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GetFailoverAllTestDetails: %+v", err)
	}

	return nil
}
