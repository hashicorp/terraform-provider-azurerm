package networkmanagers

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

type NetworkManagerCommitsPostOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NetworkManagerCommit
}

// NetworkManagerCommitsPost ...
func (c NetworkManagersClient) NetworkManagerCommitsPost(ctx context.Context, id NetworkManagerId, input NetworkManagerCommit) (result NetworkManagerCommitsPostOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/commit", id.ID()),
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

// NetworkManagerCommitsPostThenPoll performs NetworkManagerCommitsPost then polls until it's completed
func (c NetworkManagersClient) NetworkManagerCommitsPostThenPoll(ctx context.Context, id NetworkManagerId, input NetworkManagerCommit) error {
	result, err := c.NetworkManagerCommitsPost(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing NetworkManagerCommitsPost: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after NetworkManagerCommitsPost: %+v", err)
	}

	return nil
}
