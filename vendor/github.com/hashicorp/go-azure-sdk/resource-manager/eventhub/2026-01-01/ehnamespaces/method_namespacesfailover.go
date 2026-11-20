package ehnamespaces

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

type NamespacesFailoverOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// NamespacesFailover ...
func (c EHNamespacesClient) NamespacesFailover(ctx context.Context, id NamespaceId, input FailOver) (result NamespacesFailoverOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/failover", id.ID()),
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

// NamespacesFailoverThenPoll performs NamespacesFailover then polls until it's completed
func (c EHNamespacesClient) NamespacesFailoverThenPoll(ctx context.Context, id NamespaceId, input FailOver) error {
	return c.NamespacesFailoverCallbackThenPoll(ctx, id, input, nil)
}

// NamespacesFailoverCallbackThenPoll performs NamespacesFailover, runs the optional callback function, then polls until it's completed
func (c EHNamespacesClient) NamespacesFailoverCallbackThenPoll(ctx context.Context, id NamespaceId, input FailOver, callback func() error) error {
	result, err := c.NamespacesFailover(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing NamespacesFailover: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after NamespacesFailover: %+v", err)
	}

	return nil
}
