package servers

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

type DeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// Delete ...
<<<<<<< HEAD
func (c ServersClient) Delete(ctx context.Context, id ServerId) (result DeleteOperationResponse, err error) {
=======
func (c ServersClient) Delete(ctx context.Context, id commonids.SqlServerId) (result DeleteOperationResponse, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
			http.StatusOK,
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

// DeleteThenPoll performs Delete then polls until it's completed
<<<<<<< HEAD
func (c ServersClient) DeleteThenPoll(ctx context.Context, id ServerId) error {
=======
func (c ServersClient) DeleteThenPoll(ctx context.Context, id commonids.SqlServerId) error {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	result, err := c.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}
