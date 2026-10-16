package openapis

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

type CassandraClustersDeallocateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type CassandraClustersDeallocateOperationOptions struct {
	XMsForceDeallocate *string
}

func DefaultCassandraClustersDeallocateOperationOptions() CassandraClustersDeallocateOperationOptions {
	return CassandraClustersDeallocateOperationOptions{}
}

func (o CassandraClustersDeallocateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsForceDeallocate != nil {
		out.Append("x-ms-force-deallocate", fmt.Sprintf("%v", *o.XMsForceDeallocate))
	}
	return &out
}

func (o CassandraClustersDeallocateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CassandraClustersDeallocateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CassandraClustersDeallocate ...
func (c OpenapisClient) CassandraClustersDeallocate(ctx context.Context, id CassandraClusterId, options CassandraClustersDeallocateOperationOptions) (result CassandraClustersDeallocateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/deallocate", id.ID()),
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

// CassandraClustersDeallocateThenPoll performs CassandraClustersDeallocate then polls until it's completed
func (c OpenapisClient) CassandraClustersDeallocateThenPoll(ctx context.Context, id CassandraClusterId, options CassandraClustersDeallocateOperationOptions) error {
	return c.CassandraClustersDeallocateCallbackThenPoll(ctx, id, options, nil)
}

// CassandraClustersDeallocateCallbackThenPoll performs CassandraClustersDeallocate, runs the optional callback function, then polls until it's completed
func (c OpenapisClient) CassandraClustersDeallocateCallbackThenPoll(ctx context.Context, id CassandraClusterId, options CassandraClustersDeallocateOperationOptions, callback func() error) error {
	result, err := c.CassandraClustersDeallocate(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersDeallocate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CassandraClustersDeallocate: %+v", err)
	}

	return nil
}
