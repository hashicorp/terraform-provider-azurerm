package settings

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Setting
}

// Update ...
func (c SettingsClient) Update(ctx context.Context, id SettingId, input Setting) (result UpdateOperationResponse, err error) {
	req, err := c.preparerForUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "settings.SettingsClient", "Update", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "settings.SettingsClient", "Update", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "settings.SettingsClient", "Update", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdate prepares the Update request.
func (c SettingsClient) preparerForUpdate(ctx context.Context, id SettingId, input Setting) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdate handles the response to the Update request. The method always
// closes the http.Response Body.
func (c SettingsClient) responderForUpdate(resp *http.Response) (result UpdateOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	if err != nil {
		return
	}
	model, err := unmarshalSettingImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
