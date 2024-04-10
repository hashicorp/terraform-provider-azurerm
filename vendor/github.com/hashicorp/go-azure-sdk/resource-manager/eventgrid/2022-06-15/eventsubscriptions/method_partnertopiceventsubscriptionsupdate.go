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

type PartnerTopicEventSubscriptionsUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *EventSubscription
}

// PartnerTopicEventSubscriptionsUpdate ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsUpdate(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscriptionUpdateParameters) (result PartnerTopicEventSubscriptionsUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
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

// PartnerTopicEventSubscriptionsUpdateThenPoll performs PartnerTopicEventSubscriptionsUpdate then polls until it's completed
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsUpdateThenPoll(ctx context.Context, id PartnerTopicEventSubscriptionId, input EventSubscriptionUpdateParameters) error {
	result, err := c.PartnerTopicEventSubscriptionsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PartnerTopicEventSubscriptionsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after PartnerTopicEventSubscriptionsUpdate: %+v", err)
	}

	return nil
}
