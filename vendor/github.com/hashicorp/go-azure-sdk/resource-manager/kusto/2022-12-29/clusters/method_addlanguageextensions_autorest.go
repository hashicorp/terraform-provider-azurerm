package clusters

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

type AddLanguageExtensionsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AddLanguageExtensions ...
func (c ClustersClient) AddLanguageExtensions(ctx context.Context, id ClusterId, input LanguageExtensionsList) (result AddLanguageExtensionsOperationResponse, err error) {
	req, err := c.preparerForAddLanguageExtensions(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "AddLanguageExtensions", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAddLanguageExtensions(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "AddLanguageExtensions", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AddLanguageExtensionsThenPoll performs AddLanguageExtensions then polls until it's completed
func (c ClustersClient) AddLanguageExtensionsThenPoll(ctx context.Context, id ClusterId, input LanguageExtensionsList) error {
	result, err := c.AddLanguageExtensions(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AddLanguageExtensions: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AddLanguageExtensions: %+v", err)
	}

	return nil
}

// preparerForAddLanguageExtensions prepares the AddLanguageExtensions request.
func (c ClustersClient) preparerForAddLanguageExtensions(ctx context.Context, id ClusterId, input LanguageExtensionsList) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/addLanguageExtensions", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAddLanguageExtensions sends the AddLanguageExtensions request. The method will close the
// http.Response Body if it receives an error.
func (c ClustersClient) senderForAddLanguageExtensions(ctx context.Context, req *http.Request) (future AddLanguageExtensionsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
