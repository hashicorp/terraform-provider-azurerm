package machinelearningcomputes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type ComputeDeleteOperationOptions struct {
	UnderlyingResourceAction *UnderlyingResourceAction
}

func DefaultComputeDeleteOperationOptions() ComputeDeleteOperationOptions {
	return ComputeDeleteOperationOptions{}
}

func (o ComputeDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ComputeDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.UnderlyingResourceAction != nil {
		out["underlyingResourceAction"] = *o.UnderlyingResourceAction
	}

	return out
}

// ComputeDelete ...
func (c MachineLearningComputesClient) ComputeDelete(ctx context.Context, id ComputeId, options ComputeDeleteOperationOptions) (result ComputeDeleteOperationResponse, err error) {
	req, err := c.preparerForComputeDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForComputeDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ComputeDeleteThenPoll performs ComputeDelete then polls until it's completed
func (c MachineLearningComputesClient) ComputeDeleteThenPoll(ctx context.Context, id ComputeId, options ComputeDeleteOperationOptions) error {
	result, err := c.ComputeDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ComputeDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ComputeDelete: %+v", err)
	}

	return nil
}

// preparerForComputeDelete prepares the ComputeDelete request.
func (c MachineLearningComputesClient) preparerForComputeDelete(ctx context.Context, id ComputeId, options ComputeDeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForComputeDelete sends the ComputeDelete request. The method will close the
// http.Response Body if it receives an error.
func (c MachineLearningComputesClient) senderForComputeDelete(ctx context.Context, req *http.Request) (future ComputeDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
