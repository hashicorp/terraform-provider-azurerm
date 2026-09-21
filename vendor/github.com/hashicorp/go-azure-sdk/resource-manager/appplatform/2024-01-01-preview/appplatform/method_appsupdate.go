package appplatform

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

type AppsUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AppResource
}

// AppsUpdate ...
func (c AppPlatformClient) AppsUpdate(ctx context.Context, id AppId, input AppResource) (result AppsUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPatch,
		Path:       id.ID(),
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

// AppsUpdateThenPoll performs AppsUpdate then polls until it's completed
func (c AppPlatformClient) AppsUpdateThenPoll(ctx context.Context, id AppId, input AppResource) error {
	return c.AppsUpdateCallbackThenPoll(ctx, id, input, nil)
}

// AppsUpdateCallbackThenPoll performs AppsUpdate, runs the optional callback function, then polls until it's completed
func (c AppPlatformClient) AppsUpdateCallbackThenPoll(ctx context.Context, id AppId, input AppResource, callback func() error) error {
	result, err := c.AppsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AppsUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AppsUpdate: %+v", err)
	}

	return nil
}
