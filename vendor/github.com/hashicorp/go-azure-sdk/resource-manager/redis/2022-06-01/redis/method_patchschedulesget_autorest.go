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

type PatchSchedulesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *RedisPatchSchedule
}

// PatchSchedulesGet ...
func (c RedisClient) PatchSchedulesGet(ctx context.Context, id RediId) (result PatchSchedulesGetOperationResponse, err error) {
	req, err := c.preparerForPatchSchedulesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatchSchedulesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatchSchedulesGet prepares the PatchSchedulesGet request.
func (c RedisClient) preparerForPatchSchedulesGet(ctx context.Context, id RediId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/patchSchedules/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatchSchedulesGet handles the response to the PatchSchedulesGet request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForPatchSchedulesGet(resp *http.Response) (result PatchSchedulesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
