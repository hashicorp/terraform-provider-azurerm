package cloudendpointresource

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

type CloudEndpointsPreRestoreOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// CloudEndpointsPreRestore ...
func (c CloudEndpointResourceClient) CloudEndpointsPreRestore(ctx context.Context, id CloudEndpointId, input PreRestoreRequest) (result CloudEndpointsPreRestoreOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/prerestore", id.ID()),
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

// CloudEndpointsPreRestoreThenPoll performs CloudEndpointsPreRestore then polls until it's completed
func (c CloudEndpointResourceClient) CloudEndpointsPreRestoreThenPoll(ctx context.Context, id CloudEndpointId, input PreRestoreRequest) error {
	result, err := c.CloudEndpointsPreRestore(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CloudEndpointsPreRestore: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CloudEndpointsPreRestore: %+v", err)
	}

	return nil
}
