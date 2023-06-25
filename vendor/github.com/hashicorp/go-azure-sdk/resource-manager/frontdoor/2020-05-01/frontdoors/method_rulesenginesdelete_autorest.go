package frontdoors

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

type RulesEnginesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RulesEnginesDelete ...
func (c FrontDoorsClient) RulesEnginesDelete(ctx context.Context, id RulesEngineId) (result RulesEnginesDeleteOperationResponse, err error) {
	req, err := c.preparerForRulesEnginesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRulesEnginesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RulesEnginesDeleteThenPoll performs RulesEnginesDelete then polls until it's completed
func (c FrontDoorsClient) RulesEnginesDeleteThenPoll(ctx context.Context, id RulesEngineId) error {
	result, err := c.RulesEnginesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RulesEnginesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RulesEnginesDelete: %+v", err)
	}

	return nil
}

// preparerForRulesEnginesDelete prepares the RulesEnginesDelete request.
func (c FrontDoorsClient) preparerForRulesEnginesDelete(ctx context.Context, id RulesEngineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRulesEnginesDelete sends the RulesEnginesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForRulesEnginesDelete(ctx context.Context, req *http.Request) (future RulesEnginesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
