package webapps

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

type DeletePrivateEndpointConnectionSlotOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *interface{}
}

// DeletePrivateEndpointConnectionSlot ...
func (c WebAppsClient) DeletePrivateEndpointConnectionSlot(ctx context.Context, id SlotPrivateEndpointConnectionId) (result DeletePrivateEndpointConnectionSlotOperationResponse, err error) {
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

// DeletePrivateEndpointConnectionSlotThenPoll performs DeletePrivateEndpointConnectionSlot then polls until it's completed
func (c WebAppsClient) DeletePrivateEndpointConnectionSlotThenPoll(ctx context.Context, id SlotPrivateEndpointConnectionId) error {
	result, err := c.DeletePrivateEndpointConnectionSlot(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeletePrivateEndpointConnectionSlot: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeletePrivateEndpointConnectionSlot: %+v", err)
	}

	return nil
}
