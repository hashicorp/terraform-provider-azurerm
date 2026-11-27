package edgeactionversions

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

type SwapDefaultOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// SwapDefault ...
func (c EdgeActionVersionsClient) SwapDefault(ctx context.Context, id VersionId) (result SwapDefaultOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/swapDefault", id.ID()),
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

// SwapDefaultThenPoll performs SwapDefault then polls until it's completed
func (c EdgeActionVersionsClient) SwapDefaultThenPoll(ctx context.Context, id VersionId) error {
	return c.SwapDefaultCallbackThenPoll(ctx, id, nil)
}

// SwapDefaultCallbackThenPoll performs SwapDefault, runs the optional callback function, then polls until it's completed
func (c EdgeActionVersionsClient) SwapDefaultCallbackThenPoll(ctx context.Context, id VersionId, callback func() error) error {
	result, err := c.SwapDefault(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SwapDefault: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SwapDefault: %+v", err)
	}

	return nil
}
