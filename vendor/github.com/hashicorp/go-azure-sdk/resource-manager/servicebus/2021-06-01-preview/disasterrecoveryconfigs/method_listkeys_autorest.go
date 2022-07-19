package disasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// ListKeys ...
func (c DisasterRecoveryConfigsClient) ListKeys(ctx context.Context, id AuthorizationRuleId) (result ListKeysOperationResponse, err error) {
	req, err := c.preparerForListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "ListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "ListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disasterrecoveryconfigs.DisasterRecoveryConfigsClient", "ListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKeys prepares the ListKeys request.
func (c DisasterRecoveryConfigsClient) preparerForListKeys(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKeys handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (c DisasterRecoveryConfigsClient) responderForListKeys(resp *http.Response) (result ListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
