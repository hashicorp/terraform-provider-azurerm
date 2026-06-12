package redisenterprise

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

type DatabasesRegenerateKeyOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AccessKeys
}

// DatabasesRegenerateKey ...
func (c RedisEnterpriseClient) DatabasesRegenerateKey(ctx context.Context, id DatabaseId, input RegenerateKeyParameters) (result DatabasesRegenerateKeyOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/regenerateKey", id.ID()),
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

// DatabasesRegenerateKeyThenPoll performs DatabasesRegenerateKey then polls until it's completed
func (c RedisEnterpriseClient) DatabasesRegenerateKeyThenPoll(ctx context.Context, id DatabaseId, input RegenerateKeyParameters) error {
	return c.DatabasesRegenerateKeyCallbackThenPoll(ctx, id, input, nil)
}

// DatabasesRegenerateKeyCallbackThenPoll performs DatabasesRegenerateKey, runs the optional callback function, then polls until it's completed
func (c RedisEnterpriseClient) DatabasesRegenerateKeyCallbackThenPoll(ctx context.Context, id DatabaseId, input RegenerateKeyParameters, callback func() error) error {
	result, err := c.DatabasesRegenerateKey(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesRegenerateKey: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DatabasesRegenerateKey: %+v", err)
	}

	return nil
}
