package databases

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

type ImportOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// Import ...
func (c DatabasesClient) Import(ctx context.Context, id DatabaseId, input ImportClusterParameters) (result ImportOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/import", id.ID()),
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

// ImportThenPoll performs Import then polls until it's completed
func (c DatabasesClient) ImportThenPoll(ctx context.Context, id DatabaseId, input ImportClusterParameters) error {
	return c.ImportCallbackThenPoll(ctx, id, input, nil)
}

// ImportCallbackThenPoll performs Import, runs the optional callback function, then polls until it's completed
func (c DatabasesClient) ImportCallbackThenPoll(ctx context.Context, id DatabaseId, input ImportClusterParameters, callback func() error) error {
	result, err := c.Import(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Import: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Import: %+v", err)
	}

	return nil
}
