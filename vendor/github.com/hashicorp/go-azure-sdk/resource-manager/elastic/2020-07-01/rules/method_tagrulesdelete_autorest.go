package rules

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

type TagRulesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TagRulesDelete ...
func (c RulesClient) TagRulesDelete(ctx context.Context, id TagRuleId) (result TagRulesDeleteOperationResponse, err error) {
	req, err := c.preparerForTagRulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTagRulesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TagRulesDeleteThenPoll performs TagRulesDelete then polls until it's completed
func (c RulesClient) TagRulesDeleteThenPoll(ctx context.Context, id TagRuleId) error {
	result, err := c.TagRulesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TagRulesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TagRulesDelete: %+v", err)
	}

	return nil
}

// preparerForTagRulesDelete prepares the TagRulesDelete request.
func (c RulesClient) preparerForTagRulesDelete(ctx context.Context, id TagRuleId) (*http.Request, error) {
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

// senderForTagRulesDelete sends the TagRulesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c RulesClient) senderForTagRulesDelete(ctx context.Context, req *http.Request) (future TagRulesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
