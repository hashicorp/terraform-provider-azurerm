package eventsubscriptions

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

type PartnerTopicEventSubscriptionsDeleteOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// PartnerTopicEventSubscriptionsDelete ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsDelete(ctx context.Context, id PartnerTopicEventSubscriptionId) (result PartnerTopicEventSubscriptionsDeleteOperationResponse, err error) {
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

// PartnerTopicEventSubscriptionsDeleteThenPoll performs PartnerTopicEventSubscriptionsDelete then polls until it's completed
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsDeleteThenPoll(ctx context.Context, id PartnerTopicEventSubscriptionId) error {
	result, err := c.PartnerTopicEventSubscriptionsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PartnerTopicEventSubscriptionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after PartnerTopicEventSubscriptionsDelete: %+v", err)
	}

	return nil
}
