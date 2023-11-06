package databases

import (
	"context"
	"fmt"
	"net/http"

<<<<<<< HEAD
=======
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type FailoverOperationOptions struct {
	ReplicaType *ReplicaType
}

func DefaultFailoverOperationOptions() FailoverOperationOptions {
	return FailoverOperationOptions{}
}

func (o FailoverOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o FailoverOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o FailoverOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ReplicaType != nil {
		out.Append("replicaType", fmt.Sprintf("%v", *o.ReplicaType))
	}
	return &out
}

// Failover ...
<<<<<<< HEAD
func (c DatabasesClient) Failover(ctx context.Context, id DatabaseId, options FailoverOperationOptions) (result FailoverOperationResponse, err error) {
=======
func (c DatabasesClient) Failover(ctx context.Context, id commonids.SqlDatabaseId, options FailoverOperationOptions) (result FailoverOperationResponse, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		Path:          fmt.Sprintf("%s/failover", id.ID()),
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

// FailoverThenPoll performs Failover then polls until it's completed
<<<<<<< HEAD
func (c DatabasesClient) FailoverThenPoll(ctx context.Context, id DatabaseId, options FailoverOperationOptions) error {
=======
func (c DatabasesClient) FailoverThenPoll(ctx context.Context, id commonids.SqlDatabaseId, options FailoverOperationOptions) error {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	result, err := c.Failover(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Failover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Failover: %+v", err)
	}

	return nil
}
