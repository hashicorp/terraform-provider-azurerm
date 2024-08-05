package share

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

type ProviderShareSubscriptionsRevokeOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ProviderShareSubscription
}

// ProviderShareSubscriptionsRevoke ...
func (c ShareClient) ProviderShareSubscriptionsRevoke(ctx context.Context, id ProviderShareSubscriptionId) (result ProviderShareSubscriptionsRevokeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/revoke", id.ID()),
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

// ProviderShareSubscriptionsRevokeThenPoll performs ProviderShareSubscriptionsRevoke then polls until it's completed
func (c ShareClient) ProviderShareSubscriptionsRevokeThenPoll(ctx context.Context, id ProviderShareSubscriptionId) error {
	result, err := c.ProviderShareSubscriptionsRevoke(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ProviderShareSubscriptionsRevoke: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ProviderShareSubscriptionsRevoke: %+v", err)
	}

	return nil
}
