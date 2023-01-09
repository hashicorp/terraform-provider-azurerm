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

type RemoveLanguageExtensionsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RemoveLanguageExtensions ...
func (c ClustersClient) RemoveLanguageExtensions(ctx context.Context, id ClusterId, input LanguageExtensionsList) (result RemoveLanguageExtensionsOperationResponse, err error) {
	req, err := c.preparerForRemoveLanguageExtensions(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "RemoveLanguageExtensions", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRemoveLanguageExtensions(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "RemoveLanguageExtensions", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RemoveLanguageExtensionsThenPoll performs RemoveLanguageExtensions then polls until it's completed
func (c ClustersClient) RemoveLanguageExtensionsThenPoll(ctx context.Context, id ClusterId, input LanguageExtensionsList) error {
	result, err := c.RemoveLanguageExtensions(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RemoveLanguageExtensions: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RemoveLanguageExtensions: %+v", err)
	}

	return nil
}

// preparerForRemoveLanguageExtensions prepares the RemoveLanguageExtensions request.
func (c ClustersClient) preparerForRemoveLanguageExtensions(ctx context.Context, id ClusterId, input LanguageExtensionsList) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/removeLanguageExtensions", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRemoveLanguageExtensions sends the RemoveLanguageExtensions request. The method will close the
// http.Response Body if it receives an error.
func (c ClustersClient) senderForRemoveLanguageExtensions(ctx context.Context, req *http.Request) (future RemoveLanguageExtensionsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
