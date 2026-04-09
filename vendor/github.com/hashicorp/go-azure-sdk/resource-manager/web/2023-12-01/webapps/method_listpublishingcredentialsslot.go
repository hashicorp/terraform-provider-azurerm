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

type ListPublishingCredentialsSlotOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *User
}

// ListPublishingCredentialsSlot ...
func (c WebAppsClient) ListPublishingCredentialsSlot(ctx context.Context, id SlotId) (result ListPublishingCredentialsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/config/publishingcredentials/list", id.ID()),
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

// ListPublishingCredentialsSlotThenPoll performs ListPublishingCredentialsSlot then polls until it's completed
func (c WebAppsClient) ListPublishingCredentialsSlotThenPoll(ctx context.Context, id SlotId) error {
	result, err := c.ListPublishingCredentialsSlot(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ListPublishingCredentialsSlot: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ListPublishingCredentialsSlot: %+v", err)
	}

	return nil
}
