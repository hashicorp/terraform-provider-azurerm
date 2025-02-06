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

type ApplicationLiveViewsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ApplicationLiveViewsDelete ...
func (c AppPlatformClient) ApplicationLiveViewsDelete(ctx context.Context, id ApplicationLiveViewId) (result ApplicationLiveViewsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
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

// ApplicationLiveViewsDeleteThenPoll performs ApplicationLiveViewsDelete then polls until it's completed
func (c AppPlatformClient) ApplicationLiveViewsDeleteThenPoll(ctx context.Context, id ApplicationLiveViewId) error {
	result, err := c.ApplicationLiveViewsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ApplicationLiveViewsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ApplicationLiveViewsDelete: %+v", err)
	}

	return nil
}
