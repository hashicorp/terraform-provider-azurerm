package connectionanalyzers

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

type NetworkWatchersConnectionAnalyzersQueryOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ConnectionAnalyzerQueryStatusResult
}

// NetworkWatchersConnectionAnalyzersQuery ...
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersQuery(ctx context.Context, id ConnectionAnalyzerId) (result NetworkWatchersConnectionAnalyzersQueryOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/query", id.ID()),
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

// NetworkWatchersConnectionAnalyzersQueryThenPoll performs NetworkWatchersConnectionAnalyzersQuery then polls until it's completed
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersQueryThenPoll(ctx context.Context, id ConnectionAnalyzerId) error {
	return c.NetworkWatchersConnectionAnalyzersQueryCallbackThenPoll(ctx, id, nil)
}

// NetworkWatchersConnectionAnalyzersQueryCallbackThenPoll performs NetworkWatchersConnectionAnalyzersQuery, runs the optional callback function, then polls until it's completed
func (c ConnectionAnalyzersClient) NetworkWatchersConnectionAnalyzersQueryCallbackThenPoll(ctx context.Context, id ConnectionAnalyzerId, callback func() error) error {
	result, err := c.NetworkWatchersConnectionAnalyzersQuery(ctx, id)
	if err != nil {
		return fmt.Errorf("performing NetworkWatchersConnectionAnalyzersQuery: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after NetworkWatchersConnectionAnalyzersQuery: %+v", err)
	}

	return nil
}
