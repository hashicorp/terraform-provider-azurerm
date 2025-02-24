package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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
func (c DatabasesClient) Failover(ctx context.Context, id commonids.SqlDatabaseId, options FailoverOperationOptions) (result FailoverOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/failover", id.ID()),
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
func (c DatabasesClient) FailoverThenPoll(ctx context.Context, id commonids.SqlDatabaseId, options FailoverOperationOptions) error {
	result, err := c.Failover(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Failover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Failover: %+v", err)
	}

	return nil
}
