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

type BackendHealthOnDemandOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApplicationGatewayBackendHealthOnDemand
}

type BackendHealthOnDemandOperationOptions struct {
	Expand *string
}

func DefaultBackendHealthOnDemandOperationOptions() BackendHealthOnDemandOperationOptions {
	return BackendHealthOnDemandOperationOptions{}
}

func (o BackendHealthOnDemandOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o BackendHealthOnDemandOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o BackendHealthOnDemandOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// BackendHealthOnDemand ...
func (c ApplicationGatewaysClient) BackendHealthOnDemand(ctx context.Context, id ApplicationGatewayId, input ApplicationGatewayOnDemandProbe, options BackendHealthOnDemandOperationOptions) (result BackendHealthOnDemandOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/getBackendHealthOnDemand", id.ID()),
		OptionsObject: options,
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

// BackendHealthOnDemandThenPoll performs BackendHealthOnDemand then polls until it's completed
func (c ApplicationGatewaysClient) BackendHealthOnDemandThenPoll(ctx context.Context, id ApplicationGatewayId, input ApplicationGatewayOnDemandProbe, options BackendHealthOnDemandOperationOptions) error {
	result, err := c.BackendHealthOnDemand(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing BackendHealthOnDemand: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after BackendHealthOnDemand: %+v", err)
	}

	return nil
}