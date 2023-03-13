package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSchedulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// PatchSchedulesDelete ...
func (c RedisClient) PatchSchedulesDelete(ctx context.Context, id RediId) (result PatchSchedulesDeleteOperationResponse, err error) {
	req, err := c.preparerForPatchSchedulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatchSchedulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatchSchedulesDelete prepares the PatchSchedulesDelete request.
func (c RedisClient) preparerForPatchSchedulesDelete(ctx context.Context, id RediId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/patchSchedules/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatchSchedulesDelete handles the response to the PatchSchedulesDelete request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForPatchSchedulesDelete(resp *http.Response) (result PatchSchedulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
