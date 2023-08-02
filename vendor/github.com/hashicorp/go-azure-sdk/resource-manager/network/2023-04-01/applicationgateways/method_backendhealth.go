package applicationgateways

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

type BackendHealthOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type BackendHealthOperationOptions struct {
	Expand *string
}

func DefaultBackendHealthOperationOptions() BackendHealthOperationOptions {
	return BackendHealthOperationOptions{}
}

func (o BackendHealthOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o BackendHealthOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o BackendHealthOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// BackendHealth ...
func (c ApplicationGatewaysClient) BackendHealth(ctx context.Context, id ApplicationGatewayId, options BackendHealthOperationOptions) (result BackendHealthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/backendhealth", id.ID()),
		OptionsObject: options,
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

// BackendHealthThenPoll performs BackendHealth then polls until it's completed
func (c ApplicationGatewaysClient) BackendHealthThenPoll(ctx context.Context, id ApplicationGatewayId, options BackendHealthOperationOptions) error {
	result, err := c.BackendHealth(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing BackendHealth: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackendHealth: %+v", err)
	}

	return nil
}
