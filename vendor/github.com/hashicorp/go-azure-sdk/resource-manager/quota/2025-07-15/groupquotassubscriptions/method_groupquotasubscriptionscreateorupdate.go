package groupquotassubscriptions

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

type GroupQuotaSubscriptionsCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GroupQuotaSubscriptionId
}

// GroupQuotaSubscriptionsCreateOrUpdate ...
func (c GroupQuotasSubscriptionsClient) GroupQuotaSubscriptionsCreateOrUpdate(ctx context.Context, id SubscriptionId) (result GroupQuotaSubscriptionsCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
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

// GroupQuotaSubscriptionsCreateOrUpdateThenPoll performs GroupQuotaSubscriptionsCreateOrUpdate then polls until it's completed
func (c GroupQuotasSubscriptionsClient) GroupQuotaSubscriptionsCreateOrUpdateThenPoll(ctx context.Context, id SubscriptionId) error {
	return c.GroupQuotaSubscriptionsCreateOrUpdateCallbackThenPoll(ctx, id, nil)
}

// GroupQuotaSubscriptionsCreateOrUpdateCallbackThenPoll performs GroupQuotaSubscriptionsCreateOrUpdate, runs the optional callback function, then polls until it's completed
func (c GroupQuotasSubscriptionsClient) GroupQuotaSubscriptionsCreateOrUpdateCallbackThenPoll(ctx context.Context, id SubscriptionId, callback func() error) error {
	result, err := c.GroupQuotaSubscriptionsCreateOrUpdate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GroupQuotaSubscriptionsCreateOrUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GroupQuotaSubscriptionsCreateOrUpdate: %+v", err)
	}

	return nil
}
