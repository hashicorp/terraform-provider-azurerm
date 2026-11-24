package servers

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

type MigrateNetworkModeOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *MigrateNetworkStatus
}

// MigrateNetworkMode ...
func (c ServersClient) MigrateNetworkMode(ctx context.Context, id FlexibleServerId) (result MigrateNetworkModeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/migrateNetwork", id.ID()),
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

// MigrateNetworkModeThenPoll performs MigrateNetworkMode then polls until it's completed
func (c ServersClient) MigrateNetworkModeThenPoll(ctx context.Context, id FlexibleServerId) error {
	return c.MigrateNetworkModeCallbackThenPoll(ctx, id, nil)
}

// MigrateNetworkModeCallbackThenPoll performs MigrateNetworkMode, runs the optional callback function, then polls until it's completed
func (c ServersClient) MigrateNetworkModeCallbackThenPoll(ctx context.Context, id FlexibleServerId, callback func() error) error {
	result, err := c.MigrateNetworkMode(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MigrateNetworkMode: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MigrateNetworkMode: %+v", err)
	}

	return nil
}
