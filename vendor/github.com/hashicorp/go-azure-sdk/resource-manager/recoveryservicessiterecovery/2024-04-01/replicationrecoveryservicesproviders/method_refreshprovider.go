package replicationrecoveryservicesproviders

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

type RefreshProviderOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RecoveryServicesProvider
}

// RefreshProvider ...
func (c ReplicationRecoveryServicesProvidersClient) RefreshProvider(ctx context.Context, id ReplicationRecoveryServicesProviderId) (result RefreshProviderOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/refreshProvider", id.ID()),
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

// RefreshProviderThenPoll performs RefreshProvider then polls until it's completed
func (c ReplicationRecoveryServicesProvidersClient) RefreshProviderThenPoll(ctx context.Context, id ReplicationRecoveryServicesProviderId) error {
	return c.RefreshProviderCallbackThenPoll(ctx, id, nil)
}

// RefreshProviderCallbackThenPoll performs RefreshProvider, runs the optional callback function, then polls until it's completed
func (c ReplicationRecoveryServicesProvidersClient) RefreshProviderCallbackThenPoll(ctx context.Context, id ReplicationRecoveryServicesProviderId, callback func() error) error {
	result, err := c.RefreshProvider(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RefreshProvider: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after RefreshProvider: %+v", err)
	}

	return nil
}
