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

type PatchSchedulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *RedisPatchSchedule
}

// PatchSchedulesCreateOrUpdate ...
func (c RedisClient) PatchSchedulesCreateOrUpdate(ctx context.Context, id RediId, input RedisPatchSchedule) (result PatchSchedulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForPatchSchedulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatchSchedulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatchSchedulesCreateOrUpdate prepares the PatchSchedulesCreateOrUpdate request.
func (c RedisClient) preparerForPatchSchedulesCreateOrUpdate(ctx context.Context, id RediId, input RedisPatchSchedule) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/patchSchedules/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatchSchedulesCreateOrUpdate handles the response to the PatchSchedulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForPatchSchedulesCreateOrUpdate(resp *http.Response) (result PatchSchedulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
