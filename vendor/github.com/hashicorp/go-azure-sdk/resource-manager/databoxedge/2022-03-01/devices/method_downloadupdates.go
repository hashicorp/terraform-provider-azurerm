package devices

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

type DownloadUpdatesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// DownloadUpdates ...
func (c DevicesClient) DownloadUpdates(ctx context.Context, id DataBoxEdgeDeviceId) (result DownloadUpdatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/downloadUpdates", id.ID()),
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

// DownloadUpdatesThenPoll performs DownloadUpdates then polls until it's completed
func (c DevicesClient) DownloadUpdatesThenPoll(ctx context.Context, id DataBoxEdgeDeviceId) error {
	return c.DownloadUpdatesCallbackThenPoll(ctx, id, nil)
}

// DownloadUpdatesCallbackThenPoll performs DownloadUpdates, runs the optional callback function, then polls until it's completed
func (c DevicesClient) DownloadUpdatesCallbackThenPoll(ctx context.Context, id DataBoxEdgeDeviceId, callback func() error) error {
	result, err := c.DownloadUpdates(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DownloadUpdates: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DownloadUpdates: %+v", err)
	}

	return nil
}
